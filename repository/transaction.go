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

	u := `INSERT INTO public.transaction (date, transaction_type, details, transaction_d_c_type, value, customer_id, sell_order) 
			VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	row := tx.QueryRowContext(ctx, u, t.Date, t.TransactionType, t.Details, t.TransactionDCType, t.Value, t.Customer.Id, t.SellOrder.Id)

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

	q := `SELECT sum(t.value), t.transaction_d_c_type  FROM public.transaction t
			WHERE transaction_type = 'Cash' and customer_id = $1
			group by t.transaction_d_c_type`

	rows, err := db.Conn.QueryContext(ctx, q, uID)
	if err != nil {
		return 0, err
	}

	defer rows.Close()

	var debit float32
	var credit float32

	for rows.Next() {
		var value float32
		var ttype string
		rows.Scan(&value, &ttype)

		switch ttype {
		case "debit":
			debit = value
		case "credit":
			debit = value
		default:
			return 0, err
		}
	}

	return debit - credit, nil
}