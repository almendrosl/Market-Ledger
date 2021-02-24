package models

type SellOrder struct {
	Id          int
	Invoice     Invoice
	SellerWants float32
	Bids        []Bid
	Ledger      []Transaction
}
