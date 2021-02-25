package api

import (
	pb "AREX-Market-Ledger/proto"
	"context"
	"google.golang.org/grpc"
	"log"
	"testing"
)

func TestSellOrders(t *testing.T) {
	tests := []struct {
		name string
		req  *pb.Empty
		res  *pb.SellOrdersResp
	}{
		{
			"success response",
			&pb.Empty{},
			&pb.SellOrdersResp{
				SellOrders: []*pb.SellOrder{},
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

			response, err := client.SellOrders(ctx, tt.req)

			if response != nil {
				if response.SellOrders == nil {
					t.Error("response: expected", tt.res.SellOrders, "received", response.SellOrders)
				}
			}

			if err != nil {
				t.Fatalf("CreateInvoice failed: %v", err)
			}
		})
	}
}
