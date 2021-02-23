package main

import (
	service "AREX-Market-Ledger/proto"
	"AREX-Market-Ledger/repository"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	netListener := getNetListener(7000)
	gRPCServer := grpc.NewServer()

	marketLedgerServiceImpl := repository.NewMarketLedgerServiceImpl()
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
