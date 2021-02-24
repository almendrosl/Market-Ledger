package models

type Customer struct {
	Id     int32
	Name   string
	Ledger []Transaction
}

type Issuer struct {
	Customer
	Invoices []Invoice
}

type Investor struct {
	Customer
	Bids []Bid
}
