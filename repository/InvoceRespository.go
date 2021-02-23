package repository

import (
	pb "AREX-Market-Ledger/proto"
	"context"
	"log"
	"strconv"
)

type MarketLedgerServiceImpl struct {
	pb.UnimplementedMarketLedgerServiceServer
}

func NewMarketLedgerServiceImpl() *MarketLedgerServiceImpl {
	return &MarketLedgerServiceImpl{}
}

//Add function implementation of gRPC Service.
func (serviceImpl *MarketLedgerServiceImpl) CreateInvoice(ctx context.Context, in *pb.CreateInvoiceReq) (*pb.CreateInvoiceResp, error) {

	log.Println("Received request for adding repository with id " + strconv.Itoa(int(in.Invoice.Id)))

	return &pb.CreateInvoiceResp{
		InvoiceId: in.Invoice.Id,
		Error:     nil,
	}, nil
}
