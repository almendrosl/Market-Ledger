package repository

import (
	"AREX-Market-Ledger/models"
	"context"
)

func (db Database) SellOrders(ctx context.Context) ([]models.SellOrder, error) {
	var sellOrders []models.SellOrder

	q := `
    SELECT t.id, t.seller_wants, i.id, i.number, i.description, i.face_value , i.issuer_id
    	FROM public.sell_order t
		join invoice i on i.id = t.invoice_id;
    `

	rows, err := db.Conn.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var so models.SellOrder
		rows.Scan(&so.Id, &so.SellerWants, &so.Invoice.Id, &so.Invoice.Number,
			&so.Invoice.Description, &so.Invoice.FaceValue, &so.Invoice.Issuer.Id)
		sellOrders = append(sellOrders, so)
	}

	return sellOrders, nil
}
