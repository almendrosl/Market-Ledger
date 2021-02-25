package api

import (
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
		}
		sor.SellOrders = append(sor.SellOrders, sellOrder)
	}

	return sor, nil

}
