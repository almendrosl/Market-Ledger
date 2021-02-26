package service

import (
	"AREX-Market-Ledger/models"
	"AREX-Market-Ledger/repository"
	"context"
	"fmt"
	"time"
)

func MatchingAlgorithm(ctx context.Context, bid models.Bid) {
	db := repository.InitDB()
	defer db.Conn.Close()
	ctx = context.Background()

	sellOrder, err := db.OneSellOrder(ctx, bid.SellOrder.Id)
	if err != nil {
		return
	}

	soDiscount := discount(sellOrder.SellerWants, sellOrder.Invoice.FaceValue)
	bidDiscount := discount(bid.Amount, bid.Size)

	if bidDiscount > soDiscount {
		err := db.DeleteBid(ctx, bid)
		if err != nil {
			return
		}
		err = releaseBalance(ctx, db, bid, bid.Amount)
		return
	}

	var totalBids float32
	for _, bidTemp := range sellOrder.Bids {
		totalBids += bidTemp.Size
	}

	sellOrder.SellOrderState = models.LOCKED

	if totalBids < sellOrder.Invoice.FaceValue {
		return
	} else if totalBids == sellOrder.Invoice.FaceValue {
		err := db.UpdateSellOrderState(ctx, sellOrder)
		if err != nil {
			return
		}
		return
	}

	//totalBids > faceValue

	//trim the bid for the rest
	newSize := bid.Size - (totalBids - sellOrder.Invoice.FaceValue)
	newAmount := newSize - (newSize * bidDiscount)
	amountToRelease := bid.Amount - newAmount

	bid.Amount = newAmount
	bid.Size = newSize

	//update the bid
	err = db.UpdateBid(ctx, bid)
	if err != nil {
		return
	}

	err = releaseBalance(ctx, db, bid, amountToRelease)
	if err != nil {
		return
	}
	err = db.UpdateSellOrderState(ctx, sellOrder)
	if err != nil {
		return
	}
}

func discount(minor float32, mayor float32) float32 {
	return 1 - (minor / mayor)
}

func releaseBalance(ctx context.Context, db repository.Database, bid models.Bid, amount float32) error {

	details := fmt.Sprintf("An ammount of €%.2f are released to investor %d;",
		amount,
		bid.Investor.Id,
	)

	_, err := db.SaveTransaction(ctx, models.Transaction{
		Date:            time.Now(),
		TransactionType: models.CASH,
		Details:         details,
		Debit:           amount,
		Credit:          models.ZeroCreditDebit,
		Customer:        models.Customer{Id: bid.Investor.Id},
		SellOrder:       models.SellOrder{Id: bid.SellOrder.Id},
	})
	if err != nil {
		return err
	}

	_, err = db.SaveTransaction(ctx, models.Transaction{
		Date:            time.Now(),
		TransactionType: models.RESERVED,
		Details:         details,
		Credit:          amount,
		Debit:           models.ZeroCreditDebit,
		Customer:        models.Customer{Id: bid.Investor.Id},
		SellOrder:       models.SellOrder{Id: bid.SellOrder.Id},
	})
	if err != nil {
		return err
	}

	return nil
}
