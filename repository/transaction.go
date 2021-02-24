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
