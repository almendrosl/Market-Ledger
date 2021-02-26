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
		Date:              time.Now(),
		TransactionType:   models.CASH,
		Details:           details,
		TransactionDCType: models.CREDIT,
		Value:             in.Bid.Amount,
		Customer:          models.Customer{Id: in.Bid.InvestorId},
		SellOrder:         models.SellOrder{Id: in.Bid.SellOrderId},
	})

	_, err = serviceImpl.db.SaveTransaction(ctx, models.Transaction{
		Date:              time.Now(),
		TransactionType:   models.RESERVED,
		Details:           details,
		TransactionDCType: models.DEBIT,
		Value:             in.Bid.Amount,
		Customer:          models.Customer{Id: in.Bid.InvestorId},
		SellOrder:         models.SellOrder{Id: in.Bid.SellOrderId},
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
