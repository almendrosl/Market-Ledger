package api

import (
	"AREX-Market-Ledger/models"
	pb "AREX-Market-Ledger/proto"
	"context"
)

func (serviceImpl *MarketLedgerServiceImpl) SellOrders(ctx context.Context, in *pb.Empty) (*pb.SellOrdersResp, error) {
	var sellOrdersResp []*pb.SellOrder
	sor := &pb.SellOrdersResp{SellOrders: sellOrdersResp}

	sellOrdersDB, err := serviceImpl.db.SellOrders(ctx)
	if err != nil {
		return &pb.SellOrdersResp{
		}, err
	}

	for _, so := range sellOrdersDB {
		sellOrder := &pb.SellOrder{
			Id: so.Id,
			Invoice: &pb.Invoice{
				Id:          so.Invoice.Id,
				Number:      so.Invoice.Number,
				Description: so.Invoice.Description,
				FaceValue:   so.Invoice.FaceValue,
				IssuerId:    so.Invoice.Issuer.Id,
			},
			SellerWants:    so.SellerWants,
			SellOrderState: mapSellOrderState(so.SellOrderState),
		}
		sor.SellOrders = append(sor.SellOrders, sellOrder)
	}

	return sor, nil

}

func (serviceImpl *MarketLedgerServiceImpl) OneSellOrder(ctx context.Context, in *pb.OneSellOrderReq) (*pb.OneSellOrderResp, error) {
	var sellOrdersResp *pb.OneSellOrderResp

	so, err := serviceImpl.db.OneSellOrder(ctx, in.SellOrderId)
	if err != nil {
		return &pb.OneSellOrderResp{
		}, err
	}

	var bids []*pb.Bid
	for _, bid := range so.Bids {
		pbBid := &pb.Bid{
			Id:          bid.Id,
			Size:        bid.Size,
			Amount:      bid.Amount,
			InvestorId:  bid.Investor.Id,
			SellOrderId: bid.SellOrder.Id,
		}
		bids = append(bids, pbBid)
	}

	var ledger []*pb.Transaction
	for _, t := range so.Ledger {
		pbT := &pb.Transaction{
			Id:              t.Id,
			Date:            t.Date.String(),
			TransactionType: mapTransactionType(t.TransactionType),
			Details:         t.Details,
			Debit:           t.Debit,
			Credit:          t.Credit,
			CustomerId:      t.Customer.Id,
			SellOrderId:     t.SellOrder.Id,
		}
		ledger = append(ledger, pbT)
	}

	sellOrdersResp = &pb.OneSellOrderResp{
		SellOrder: &pb.SellOrder{
			Id: so.Id,
			Invoice: &pb.Invoice{
				Id:          so.Invoice.Id,
				Number:      so.Invoice.Number,
				Description: so.Invoice.Description,
				FaceValue:   so.Invoice.FaceValue,
				IssuerId:    so.Invoice.Issuer.Id,
			},
			Ledger:         ledger,
			SellerWants:    so.SellerWants,
			SellOrderState: mapSellOrderState(so.SellOrderState),
			Bids:           bids,
		},
	}

	return sellOrdersResp, nil
}

func (serviceImpl *MarketLedgerServiceImpl) FinishSellOrder(ctx context.Context, in *pb.FinishSellOrderReq) (*pb.FinishSellOrderResp, error) {
	var resp *pb.FinishSellOrderResp

	so := models.SellOrder{
		Id:             in.SellOrderId,
		SellOrderState: mapFinishSellOrderState(in.FinishSellOrderType),
	}

	err := serviceImpl.db.UpdateSellOrderState(ctx, so)
	if err != nil {
		return &pb.FinishSellOrderResp{}, err
	}

	if so.SellOrderState == models.REVERSED {
		err = serviceImpl.db.RevertTransaction(ctx, so)
		if err != nil {
			return &pb.FinishSellOrderResp{}, err
		}
	} else {
		err = serviceImpl.db.CommitTransaction(ctx, so)
		if err != nil {
			return &pb.FinishSellOrderResp{}, err
		}
	}

	resp = &pb.FinishSellOrderResp{SellOrder: &pb.SellOrder{Id: so.Id}}

	return resp, nil
}

func mapFinishSellOrderState(orderType pb.FinishSellOrderReq_FinishSellOrderType) models.SellOrderState {
	switch orderType {
	case pb.FinishSellOrderReq_COMMIT:
		return models.COMMITTED
	case pb.FinishSellOrderReq_REJECT:
		return models.REVERSED
	}
	return models.REVERSED
}

func mapSellOrderState(soState models.SellOrderState) pb.SellOrder_SellOrderState {
	switch soState {
	case models.LOCKED:
		return pb.SellOrder_LOCKED
	case models.REVERSED:
		return pb.SellOrder_REVERSED
	case models.COMMITTED:
		return pb.SellOrder_COMMITTED
	case models.UNLOCKED:
		return pb.SellOrder_UNLOCKED
	}
	return pb.SellOrder_LOCKED
}

func mapTransactionType(tType models.TransactionType) pb.Transaction_TransactionType {
	switch tType {
	case models.CASH:
		return pb.Transaction_CASH
	case models.CAPITAL:
		return pb.Transaction_CAPITAL
	case models.OWN_INVOICE:
		return pb.Transaction_OWN_INVOICE
	case models.EXPECTED_PAID:
		return pb.Transaction_EXPECTED_PAID
	case models.EXPECTED_PROFIT:
		return pb.Transaction_EXPECTED_PROFIT
	case models.RESERVED:
		return pb.Transaction_RESERVED
	case models.OWES_INVOICE:
		return pb.Transaction_OWES_INVOICE
	case models.LOSS:
		return pb.Transaction_LOSS
	}
	return pb.Transaction_OWN_INVOICE
}
