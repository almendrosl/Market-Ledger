package repository

import (
	"AREX-Market-Ledger/models"
	pb "AREX-Market-Ledger/proto"
	"context"
	"log"
)

func (db Database) SellInvoice(ctx context.Context, in *pb.CreateInvoiceReq) (models.SellOrder, error) {
	var sellOrder models.SellOrder

	tx, err := db.Conn.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	u := "INSERT INTO public.invoice (number, description, face_value, issuer_id) VALUES ($1, $2, $3, $4) RETURNING id"

	row := tx.QueryRowContext(ctx, u, in.SellOrder.Invoice.Number,
		in.SellOrder.Invoice.Description, in.SellOrder.Invoice.FaceValue, in.SellOrder.Invoice.IssuerId)

	var iID int
	err = row.Scan(&iID)
	if err != nil {
		tx.Rollback()
		return sellOrder, err
	}

	// Run a query to get a count of all cats



	rowInvoice := tx.QueryRow("SELECT t.id, t.number, t.description, t.face_value FROM public.invoice t WHERE t.id = $1", iID)


	var invoice models.Invoice
	err = rowInvoice.Scan(&invoice.Id, &invoice.Number, &invoice.Description, &invoice.FaceValue)
	if err != nil {
		tx.Rollback()
		return sellOrder, err
	}

	invoice.Issuer = models.Issuer{
		Customer: models.Customer{
			Id: in.SellOrder.Invoice.IssuerId,
		},
	}

	soRow := tx.QueryRowContext(ctx, "INSERT INTO public.sell_order (invoice_id, seller_wants, sell_order_state) VALUES ($1, $2, $3) RETURNING id",
		iID, in.SellOrder.SellerWants, models.UNLOCKED)
	if err != nil {
		tx.Rollback()
		return sellOrder, err
	}

	var soID int32
	err = soRow.Scan(&soID)
	if err != nil {
		tx.Rollback()
		return sellOrder, err
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	sellOrder = models.SellOrder{
		Id:          soID,
		Invoice:     invoice,
		SellerWants: in.SellOrder.SellerWants,
		Bids:        nil,
		Ledger:      nil,
	}

	return sellOrder, nil
}
