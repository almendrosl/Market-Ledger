package api

import (
	"AREX-Market-Ledger/models"
	pb "AREX-Market-Ledger/proto"
	"AREX-Market-Ledger/repository"
	"context"
	"google.golang.org/grpc"
	"log"
	"testing"
	"time"
)

func TestPlaceBidOrders(t *testing.T) {
	tests := []struct {
		name string
		req  *pb.PlaceBidReq
		res  *pb.PlaceBidResp
	}{
		{
			"success response",
			&pb.PlaceBidReq{
				Bid: &pb.Bid{
					Size:        500,
					Amount:      270,
					InvestorId:  1,
					SellOrderId: 1,
				},
			},
			&pb.PlaceBidResp{
				Bid: &pb.Bid{
					Size:        500,
					Amount:      270,
					InvestorId:  1,
					SellOrderId: 1,
				},
			},
		},
	}

	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewMarketLedgerServiceClient(conn)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			response, err := client.PlaceBid(ctx, tt.req)

			if response != nil {
				if (response.Bid.Amount != tt.res.Bid.Amount) &&
					(response.Bid.Size != tt.res.Bid.Size) {
					t.Error("response: expected", tt.res.String(), "received", response.String())
				}
			}

			if err != nil {
				t.Fatalf("CreateInvoice failed: %v", err)
			}
		})
	}
}

func TestMatchingAlgorithm(t *testing.T) {
	tests := []struct {
		name string
		req  *pb.PlaceBidReq
		res  *models.SellOrder
	}{
		{
			"bidDiscount > soDiscount",
			&pb.PlaceBidReq{
				Bid: &pb.Bid{
					Size:        500,
					Amount:      10,
					InvestorId:  1,
					SellOrderId: 1,
				},
			},
			&models.SellOrder{
				Id:             1,
				Invoice:        models.Invoice{},
				SellerWants:    0,
				Bids:           nil,
				Ledger:         nil,
				SellOrderState: "",
			},
		},
		{
			"totalBids < FaceValue",
			&pb.PlaceBidReq{
				Bid: &pb.Bid{
					Size:        500,
					Amount:      490,
					InvestorId:  1,
					SellOrderId: 1,
				},
			},
			&models.SellOrder{
				Id:          1,
				Invoice:     models.Invoice{},
				SellerWants: 0,
				Bids: []models.Bid{{
					Size:   500,
					Amount: 490,
					Investor: models.Investor{
						Customer: models.Customer{
							Id: 1,
						},
					},
					SellOrder: models.SellOrder{},
				}},
				Ledger:         nil,
				SellOrderState: models.UNLOCKED,
			},
		},
		{
			"totalBids == FaceValue",
			&pb.PlaceBidReq{
				Bid: &pb.Bid{
					Size:        4500,
					Amount:      4490,
					InvestorId:  3,
					SellOrderId: 1,
				},
			},
			&models.SellOrder{
				Id:          1,
				Invoice:     models.Invoice{},
				SellerWants: 0,
				Bids: []models.Bid{
					{
						Size:   500,
						Amount: 490,
						Investor: models.Investor{
							Customer: models.Customer{
								Id: 1,
							},
						},
						SellOrder: models.SellOrder{},
					},
					{
						Size:   4500,
						Amount: 4490,
						Investor: models.Investor{
							Customer: models.Customer{
								Id: 3,
							},
						},
						SellOrder: models.SellOrder{},
					},
				},
				Ledger:         nil,
				SellOrderState: models.LOCKED,
			},
		},
		{
			"totalBids > FaceValue",
			&pb.PlaceBidReq{
				Bid: &pb.Bid{
					Size:        1100,
					Amount:      1090,
					InvestorId:  3,
					SellOrderId: 2,
				},
			},
			&models.SellOrder{
				Id:          2,
				Invoice:     models.Invoice{},
				SellerWants: 0,
				Bids: []models.Bid{
					{
						Size:   1000,
						Amount: 990.9091,
						Investor: models.Investor{
							Customer: models.Customer{
								Id: 3,
							},
						},
						SellOrder: models.SellOrder{},
					},
				},
				Ledger:         nil,
				SellOrderState: models.LOCKED,
			},
		},
	}

	ctx := context.Background()
	db := repository.DbInitTest()
	defer db.Conn.Close()

	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewMarketLedgerServiceClient(conn)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			response, err := client.PlaceBid(ctx, tt.req)

			if response != nil {
				time.Sleep(2 * time.Second)
				so, err := db.OneSellOrder(ctx, tt.res.Id)
				if err != nil {
					t.Fatalf("DB failed: %v", err)
				}

				if tt.res.Bids == nil {
					if so.Bids != nil {
						t.Error("response: expected no bid ", response.Bid.Id, " received ", response.String())
					}
				} else {
					if len(tt.res.Bids) == 1 {
						if (tt.res.SellOrderState != so.SellOrderState) ||
							(tt.res.Bids[0].Amount != so.Bids[0].Amount) ||
							(tt.res.Bids[0].Size != so.Bids[0].Size) {
							t.Error("response: expected ", tt.res, " received ", so)
						}
					} else {
						if (tt.res.SellOrderState != so.SellOrderState) ||
							(tt.res.Bids[0].Amount != so.Bids[0].Amount) ||
							(tt.res.Bids[0].Size != so.Bids[0].Size) {
							t.Error("response: expected ", tt.res, " received ", so)
						}
					}
				}
			}

			if err != nil {
				t.Fatalf("CreateInvoice failed: %v", err)
			}
		})
	}
}
