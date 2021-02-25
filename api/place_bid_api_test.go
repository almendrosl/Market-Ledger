package api

import (
	pb "AREX-Market-Ledger/proto"
	"context"
	"google.golang.org/grpc"
	"log"
	"testing"
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
				Bid:  &pb.Bid{
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
					(response.Bid.Size != tt.res.Bid.Size){
					t.Error("response: expected", tt.res.String(), "received", response.String())
				}
			}

			if err != nil {
				t.Fatalf("CreateInvoice failed: %v", err)
			}
		})
	}
}