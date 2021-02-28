package repository

import (
	"AREX-Market-Ledger/models"
	"context"
	"log"
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
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}
