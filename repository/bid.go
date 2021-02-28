package repository

import (
	"AREX-Market-Ledger/models"
	pb "AREX-Market-Ledger/proto"
	"context"
	"log"
)

func (db Database) SaveBid(ctx context.Context, in *pb.PlaceBidReq) (int32, error) {

	var bID int32

	tx, err := db.Conn.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	u := `INSERT INTO public.bid (size, amount, investor_id, sell_order_id) 
			VALUES ($1, $2, $3, $4) RETURNING id`

	row := tx.QueryRowContext(ctx, u, in.Bid.Size, in.Bid.Amount, in.Bid.InvestorId, in.Bid.SellOrderId)

	err = row.Scan(&bID)
	if err != nil {
		tx.Rollback()
		return bID, err
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	return bID, nil
}

func (db Database) DeleteBid(ctx context.Context, bid models.Bid) error {

	q := `DELETE FROM public.bid WHERE id =$1;`

	stmt, err := db.Conn.PrepareContext(ctx, q)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, bid.Id)
	if err != nil {
		return err
	}

	return nil
}


func (db Database) UpdateBid(ctx context.Context, bid models.Bid) error {
	q := `UPDATE public.bid SET size =$1, amount =$2 WHERE id =$3;`


	stmt, err := db.Conn.PrepareContext(ctx, q)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, bid.Size, bid.Amount, bid.Id)
	if err != nil {
		return err
	}

	return nil
}

func (db Database) GetBidsBySellOrder(ctx context.Context, so models.SellOrder) ([]models.Bid, error) {
	qb := `SELECT b.id, b.size, b.amount from bid b
		WHERE b.sell_order_id = $1
    `

	rows, err := db.Conn.QueryContext(ctx, qb, so.Id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var bids []models.Bid

	for rows.Next() {
		var bid models.Bid
		rows.Scan(&bid.Id, &bid.Size, &bid.Amount)
		bids = append(bids, bid)
	}
	return bids, nil
}