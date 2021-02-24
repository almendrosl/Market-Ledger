package api

import (
	pb "AREX-Market-Ledger/proto"
	"AREX-Market-Ledger/repository"
	"context"
	"log"
	"strconv"
)

type MarketLedgerServiceImpl struct {
	pb.UnimplementedMarketLedgerServiceServer
	db repository.Database
}

func NewMarketLedgerServiceImpl(repo repository.Database) *MarketLedgerServiceImpl {
	return &MarketLedgerServiceImpl{db: repo}
}

//Add function implementation of gRPC Service.
func (serviceImpl *MarketLedgerServiceImpl) CreateInvoice(ctx context.Context, in *pb.CreateInvoiceReq) (*pb.CreateInvoiceResp, error) {

	log.Println("Received request for adding repository with id " + strconv.Itoa(int(in.SellOrder.Invoice.Id)))
	sellOrder, err := serviceImpl.db.SellInvoice(ctx, in)

	if err != nil {
		return &pb.CreateInvoiceResp{
			Error:     nil,
		}, err
	}

	return &pb.CreateInvoiceResp{
		SellOrder: &pb.SellOrder{
			Id:          int32(sellOrder.Id),
			Invoice:     &pb.Invoice{
				Id:          int32(sellOrder.Invoice.Id),
				Number:      sellOrder.Invoice.Number,
				Description: sellOrder.Invoice.Description,
				FaceValue:   sellOrder.Invoice.FaceValue,
				IssuerId:    0,
			},
			SellerWants: sellOrder.SellerWants,
		},
		Error:     nil,
	}, err
}
