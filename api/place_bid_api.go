package api

import (
	"AREX-Market-Ledger/models"
	pb "AREX-Market-Ledger/proto"
	"AREX-Market-Ledger/service"
	"context"
	"errors"
	"fmt"
	"time"
)

func (serviceImpl *MarketLedgerServiceImpl) PlaceBid(ctx context.Context, in *pb.PlaceBidReq) (*pb.PlaceBidResp, error) {
	var bidResp *pb.PlaceBidResp

	so, err := serviceImpl.db.OneSellOrder(ctx, in.Bid.SellOrderId)
	if err != nil {
		return bidResp, err
	}

	if so.SellOrderState != models.UNLOCKED{
		return nil, errors.New("the sell order is not unlocked")
	}

	balance, err := serviceImpl.db.GetCashFromUser(ctx, in.Bid.InvestorId)
	if err != nil {
		return bidResp, err
	}

	if balance < in.Bid.Amount {
		return nil, errors.New("not enough balance")
	}

	bidID, err := serviceImpl.db.SaveBid(ctx, in)
	if err != nil {
		return nil, err
	}

	details := fmt.Sprintf("investor %d places a bid of size €%.2f and amount €%.2f;",
		in.Bid.InvestorId,
		in.Bid.Size,
		in.Bid.Amount)

	_, err = serviceImpl.db.SaveTransaction(ctx, models.Transaction{
		Date:            time.Now(),
		TransactionType: models.CASH,
		Details:         details,
		Credit:          in.Bid.Amount,
		Debit:           models.ZeroCreditDebit,
		Customer:        models.Customer{Id: in.Bid.InvestorId},
		SellOrder:       models.SellOrder{Id: in.Bid.SellOrderId},
	})

	_, err = serviceImpl.db.SaveTransaction(ctx, models.Transaction{
		Date:            time.Now(),
		TransactionType: models.RESERVED,
		Details:         details,
		Debit:           in.Bid.Amount,
		Credit:          models.ZeroCreditDebit,
		Customer:        models.Customer{Id: in.Bid.InvestorId},
		SellOrder:       models.SellOrder{Id: in.Bid.SellOrderId},
	})

	bidResp = &pb.PlaceBidResp{
		Bid: &pb.Bid{
			Id:          bidID,
			Size:        in.Bid.Size,
			Amount:      in.Bid.Amount,
			InvestorId:  in.Bid.InvestorId,
			SellOrderId: in.Bid.SellOrderId,
		},
	}

	go service.MatchingAlgorithm(ctx, models.Bid{
		Id:     bidID,
		Size:   in.Bid.Size,
		Amount: in.Bid.Amount,
		Investor: models.Investor{
			Customer: models.Customer{
				Id: in.Bid.InvestorId,
			},
		},
		SellOrder: models.SellOrder{
			Id: in.Bid.SellOrderId,
		},
	})
	return bidResp, nil
}
