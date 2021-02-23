package models

type Customer struct {
	Id      int
	Name    string
	Balance float32
}

type Issuer struct {
	Customer
	Invoices []Invoice
}

type Investor struct {
	Customer
	Bids []Bid
}
