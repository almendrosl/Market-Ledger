package api

import (
	"AREX-Market-Ledger/repository"
	"context"
	"github.com/joho/godotenv"
	"log"
	"net"
	"os"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	pb "AREX-Market-Ledger/proto"
)

func TestMain(m *testing.M) {
	var err error
	err = godotenv.Load(os.ExpandEnv("./../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}

	os.Exit(m.Run())

}

func dialer() func(context.Context, string) (net.Conn, error) {

	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()

	db := repository.InitDB()

	pb.RegisterMarketLedgerServiceServer(server, NewMarketLedgerServiceImpl(db))

	go func() {
		defer db.Conn.Close()
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func TestCreateInvoice(t *testing.T) {
	tests := []struct {
		name string
		req  *pb.CreateInvoiceReq
		res  *pb.CreateInvoiceResp
	}{
		{
			"success response",
			&pb.CreateInvoiceReq{
				SellOrder: &pb.SellOrder{
					Invoice: &pb.Invoice{
						Number:      "234-324234-12",
						Description: "invoice to sell test",
						FaceValue:   3434.43,
						IssuerId:    2,
					},
					SellerWants: 444.78,
				},
			},
			&pb.CreateInvoiceResp{
				SellOrder: nil,
				Transaction: &pb.Transaction{
					Debit: 3434.43,
				},
				Error: nil,
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

			response, err := client.CreateInvoice(ctx, tt.req)

			if response != nil {
				if response.Transaction.Debit != tt.res.Transaction.Debit {
					t.Error("response: expected", tt.res.Transaction.Debit, "received", response.Transaction.Debit)
				}
			}

			if err != nil {
				t.Fatalf("CreateInvoice failed: %v", err)
			}
		})
	}
}
