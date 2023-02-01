# Market Ledger

To compile protobuf execute in main directory

`protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/marketLedger.proto`

### Script

###### Execute script script.sql in the postgres db the first time
### Make postgres_test Data Base
`
create database postgres_test;`

To run the main file type the following command

### Run Main.go

To run the main file type the following command

`go run main.go`

### Run Docker Compose

`docker-compose up --build`



_**BloomRPC to replace Postman**_
https://github.com/uw-labs/bloomrpc

#Examples:
To create an invoice:
![Alt text](images/CreateInvoice.PNG?raw=true "Title")

View All Sell Orders:
![Alt text](images/AllSellOrders.PNG?raw=true "Title")
PLace a bid:
![Alt text](images/PlaceBid.PNG?raw=true "Title")

View Detail of a Sell Order and the ledger:
![Alt text](images/ViewOneSellOrder.PNG?raw=true "Title")

Finish Sell Order:
![Alt text](images/FinishSellOrder.PNG?raw=true "Title")
