package repository

import (
	"AREX-Market-Ledger/models"
	"context"
	"fmt"
	"log"
	"time"
)

func (db Database) SaveTransaction(ctx context.Context, t models.Transaction) (models.Transaction, error) {

	tx, err := db.Conn.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	u := `INSERT INTO public.transaction (date, transaction_type, details, debit, credit, customer_id, sell_order_id) 
			VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	row := tx.QueryRowContext(ctx, u, t.Date, t.TransactionType, t.Details, t.Debit, t.Credit, t.Customer.Id, t.SellOrder.Id)

	err = row.Scan(&t.Id)
	if err != nil {
		tx.Rollback()
		return t, err
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	return t, nil
}

func (db Database) GetCashFromUser(ctx context.Context, uID int32) (float32, error) {
	var balance float32

	q := `SELECT (sum(t.debit) - sum(credit)) as balance  FROM public.transaction t
			WHERE transaction_type = 'Cash' and customer_id = $1
`
	rows, err := db.Conn.QueryContext(ctx, q, uID)
	if err != nil {
		return 0, err
	}

	defer rows.Close()

	for rows.Next() {
		rows.Scan(&balance)
	}

	return balance, nil
}

func (db Database) GetTransactionsBySellOrder(ctx context.Context, so models.SellOrder) ([]models.Transaction, error) {
	var transactions []models.Transaction

	qb := `SELECT id, date, transaction_type, details, debit, credit, customer_id FROM public.transaction
				WHERE transaction.sell_order_id = $1
    `

	rows, err := db.Conn.QueryContext(ctx, qb, so.Id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var transaction models.Transaction
		rows.Scan(&transaction.Id, &transaction.Date, &transaction.TransactionType,
			&transaction.Details, &transaction.Debit, &transaction.Credit,
			&transaction.Customer.Id)
		transaction.SellOrder.Id = so.Id
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func (db Database) RevertTransaction(ctx context.Context, so models.SellOrder) error {
	qb := `UPDATE public.transaction
		SET    credit = CASE WHEN credit=0 THEN debit ELSE credit END
     			, debit = CASE WHEN debit=0 THEN credit ELSE debit END
		where sell_order_id = $1;
    `

	stmt, err := db.Conn.PrepareContext(ctx, qb)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, so.Id)
	if err != nil {
		return err
	}

	so.SellOrderState = models.REVERSED
	err = db.UpdateSellOrderState(ctx, so)
	if err != nil {
		return nil
	}

	return nil
}

func (db Database) CommitTransaction(ctx context.Context, so models.SellOrder) error {
	oso, err := db.OneSellOrder(ctx, so.Id)
	if err != nil {
		return err
	}

	var issuerLoss float32

	for _, bid := range oso.Bids{
		details := fmt.Sprintf("%s owes €%.2f for invoice %s to the investor %s",
			oso.Invoice.Issuer.Name,
			bid.Size,
			oso.Invoice.Number,
			bid.Investor.Name)
		t1 := models.Transaction{
			Date:            time.Now(),
			TransactionType: models.OWES_INVOICE,
			Details:         details,
			Debit:           models.ZeroCreditDebit,
			Credit:          bid.Size,
			Customer:        models.Customer{Id: oso.Invoice.Issuer.Id},
			SellOrder:       models.SellOrder{Id: oso.Id},
		}
		_, err = db.SaveTransaction(ctx,t1)
		if err != nil {
			return err
		}

		details1 := fmt.Sprintf("The amount €%.2f is not longer reserver by investor %s",
			bid.Amount,
			bid.Investor.Name)
		_, err = db.SaveTransaction(ctx, models.Transaction{
			Date:            time.Now(),
			TransactionType: models.RESERVED,
			Details:         details1,
			Debit:           models.ZeroCreditDebit,
			Credit:          bid.Amount,
			Customer:        models.Customer{Id: bid.Investor.Id},
			SellOrder:       models.SellOrder{Id: oso.Id},
		})
		if err != nil {
			return err
		}

		details3 := fmt.Sprintf("The amount €%.2f is paid for the purcharse to party %s for the invoice %s",
			bid.Amount,
			oso.Invoice.Issuer.Name,
			oso.Invoice.Number)
		_, err = db.SaveTransaction(ctx, models.Transaction{
			Date:            time.Now(),
			TransactionType: models.CASH,
			Details:         details3,
			Debit:           bid.Amount,
			Credit:          models.ZeroCreditDebit,
			Customer:        models.Customer{Id: oso.Invoice.Issuer.Id},
			SellOrder:       models.SellOrder{Id: oso.Id},
		})
		if err != nil {
			return err
		}

		details4 := fmt.Sprintf("The party %s expect €%.2f to be paid by %s",
			bid.Investor.Name,
			bid.Size,
			oso.Invoice.Issuer.Name,
		)
		_, err = db.SaveTransaction(ctx, models.Transaction{
			Date:            time.Now(),
			TransactionType: models.EXPECTED_PAID,
			Details:         details4,
			Debit:           bid.Size,
			Credit:          models.ZeroCreditDebit,
			Customer:        models.Customer{Id: bid.Investor.Id},
			SellOrder:       models.SellOrder{Id: oso.Id},
		})
		if err != nil {
			return err
		}

		expectedProfit := bid.Size - bid.Amount
		details5 := fmt.Sprintf("The party %s expected profit €%.2f",
			bid.Investor.Name,
			expectedProfit,
		)
		_, err = db.SaveTransaction(ctx, models.Transaction{
			Date:            time.Now(),
			TransactionType: models.EXPECTED_PROFIT,
			Details:         details5,
			Debit:           models.ZeroCreditDebit,
			Credit:          expectedProfit,
			Customer:        models.Customer{Id: bid.Investor.Id},
			SellOrder:       models.SellOrder{Id: oso.Id},
		})

		if err != nil {
			return err
		}

		issuerLoss += expectedProfit
	}

	details := fmt.Sprintf("The party %s owes €%.2f for invoice %s",
		oso.Invoice.Issuer.Name,
		oso.Invoice.FaceValue,
		oso.Invoice.Number,
		)
	_, err = db.SaveTransaction(ctx, models.Transaction{
		Date:            time.Now(),
		TransactionType: models.OWN_INVOICE,
		Details:         details,
		Debit:           models.ZeroCreditDebit,
		Credit:          oso.Invoice.FaceValue,
		Customer:        models.Customer{Id: oso.Invoice.Issuer.Id},
		SellOrder:       models.SellOrder{Id: oso.Id},
	})
	if err != nil {
		return err
	}


	details1 := fmt.Sprintf(`Loss of €%.2f that is given by the difference between the invoice size, and the total amount received from all the bids`,
		issuerLoss,
	)
	_, err = db.SaveTransaction(ctx, models.Transaction{
		Date:            time.Now(),
		TransactionType: models.LOSS,
		Details:         details1,
		Debit:           issuerLoss,
		Credit:          models.ZeroCreditDebit,
		Customer:        models.Customer{Id: oso.Invoice.Issuer.Id},
		SellOrder:       models.SellOrder{Id: oso.Id},
	})
	if err != nil {
		return err
	}

	oso.SellOrderState = models.COMMITTED
	err = db.UpdateSellOrderState(ctx, oso)
	if err != nil {
		return nil
	}

	return nil
}