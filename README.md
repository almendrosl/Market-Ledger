# AREX Market Ledger

To compile protobuf execute in main directory

`protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/marketLedger.proto`

### Run Main.go

To run the main file type the following command

`go run main.go`

### Run Docker Compose

`docker-compose up --build`



_**BloomRPC to replace Postman**_
https://github.com/uw-labs/bloomrpc
