package main

import (
	"AREX-Market-Ledger/api"
	service "AREX-Market-Ledger/proto"
	"AREX-Market-Ledger/repository"
	"fmt"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("sad .env file found")
	}
}

func main() {
	netListener := getNetListener(8080)
	gRPCServer := grpc.NewServer()

	database, err := repository.Initialize(os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"),
		os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))
	if err != nil {
		log.Fatalf("Could not set up database: %v", err)
	}
	defer database.Conn.Close()

	marketLedgerServiceImpl := api.NewMarketLedgerServiceImpl(database)
	service.RegisterMarketLedgerServiceServer(
		gRPCServer,
		marketLedgerServiceImpl,
	)



	// start the server
	if err := gRPCServer.Serve(netListener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

func getNetListener(port uint) net.Listener {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		panic(fmt.Sprintf("failed to listen: %v", err))
	}

	return lis
}
