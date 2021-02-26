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
			SellerWants: so.SellerWants,
			SellOrderState: mapSellOrderState(so.SellOrderState),
		}
		sor.SellOrders = append(sor.SellOrders, sellOrder)
	}

	return sor, nil

}

func (serviceImpl *MarketLedgerServiceImpl) OneSellOrder(ctx context.Context, in *pb.OneSellOrderReq) (*pb.OneSellOrderResp, error){
	var sellOrdersResp *pb.OneSellOrderResp

	so, err := serviceImpl.db.OneSellOrder(ctx, in.SellOrderId)
	if err != nil {
		return &pb.OneSellOrderResp{
		}, err
	}

	var bids []*pb.Bid
	for _, bid := range so.Bids{
		pbBid := &pb.Bid{
			Id:          bid.Id,
			Size:        bid.Size,
			Amount:      bid.Amount,
			InvestorId:  bid.Investor.Id,
			SellOrderId: bid.SellOrder.Id,
		}
		bids = append(bids, pbBid)
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
			SellerWants: so.SellerWants,
			SellOrderState: mapSellOrderState(so.SellOrderState),
			Bids: bids,
		},
	}

	return sellOrdersResp, nil
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