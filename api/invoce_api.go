package api

import (
	"AREX-Market-Ledger/models"
	pb "AREX-Market-Ledger/proto"
	"AREX-Market-Ledger/repository"
	"context"
	"fmt"
	"time"
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

	sellOrder, err := serviceImpl.db.SellInvoice(ctx, in)
	if err != nil {
		return &pb.CreateInvoiceResp{
			Error: nil,
		}, err
	}

	details := fmt.Sprintf("Issuer number %d has a €%.2f invoice number %s that should be financed",
		sellOrder.Invoice.Issuer.Customer.Id,
		sellOrder.Invoice.FaceValue,
		sellOrder.Invoice.Number)

	t := models.Transaction{
		Date:            time.Now(),
		TransactionType: models.OWN_INVOICE,
		Details:         details,
		Debit:           sellOrder.Invoice.FaceValue,
		Credit:          models.ZeroCreditDebit,
		Customer:        sellOrder.Invoice.Issuer.Customer,
		SellOrder:       sellOrder,
	}

	t, err = serviceImpl.db.SaveTransaction(ctx, t)
	if err != nil {
		return &pb.CreateInvoiceResp{
			Error: nil,
		}, err
	}

	return &pb.CreateInvoiceResp{
		SellOrder: &pb.SellOrder{
			Id: int32(sellOrder.Id),
			Invoice: &pb.Invoice{
				Id:          int32(sellOrder.Invoice.Id),
				Number:      sellOrder.Invoice.Number,
				Description: sellOrder.Invoice.Description,
				FaceValue:   sellOrder.Invoice.FaceValue,
				IssuerId:    0,
			},
			SellerWants: sellOrder.SellerWants,
		},

		Transaction: &pb.Transaction{
			Id:              int32(t.Id),
			Date:            t.Date.String(),
			TransactionType: pb.Transaction_OWN_INVOICE,
			Details:         t.Details,
			Credit:          models.ZeroCreditDebit,
			Debit:           t.Debit,
			CustomerId:      int32(t.Customer.Id),
			SellOrderId:     int32(t.SellOrder.Id),
		},
		Error: nil,
	}, err
}
