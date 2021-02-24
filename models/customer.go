package models

type Customer struct {
	Id     int
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
