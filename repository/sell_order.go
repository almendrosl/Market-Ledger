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

func (db Database) OneSellOrder(ctx context.Context, soID int32) (models.SellOrder, error) {
	var so models.SellOrder

	q := `
    SELECT t.id, t.seller_wants, i.id, i.number, i.description, i.face_value , i.issuer_id
    	FROM public.sell_order t
		join invoice i on i.id = t.invoice_id
		WHERE t.id = $1;
    `

	row := db.Conn.QueryRowContext(ctx, q, soID)

	err := row.Scan(&so.Id, &so.SellerWants, &so.Invoice.Id, &so.Invoice.Number,
		&so.Invoice.Description, &so.Invoice.FaceValue, &so.Invoice.Issuer.Id)
	if err != nil {
		return so, err
	}

	qb := `SELECT b.id, b.size, b.amount from bid b
		WHERE b.sell_order_id = $1
    `

	rows, err := db.Conn.QueryContext(ctx, qb, so.Id)
	if err != nil {
		return so, err
	}

	defer rows.Close()

	var bids []models.Bid

	for rows.Next() {
		var bid models.Bid
		rows.Scan(&bid.Id, &bid.Size, &bid.Amount)
		bids = append(bids, bid)
	}

	so.Bids = bids

	return so, nil
}


func (db Database) UpdateSellOrderState(ctx context.Context, sellOrder models.SellOrder) error {
	q := `UPDATE public.sell_order SET sell_order_state =$1 WHERE id =$2`

	stmt, err := db.Conn.PrepareContext(ctx, q)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, sellOrder.SellOrderState, sellOrder.Id)
	if err != nil {
		return err
	}

	return nil
}