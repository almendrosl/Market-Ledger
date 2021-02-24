package models

type SellOrder struct {
	Id          int32
	Invoice     Invoice
	SellerWants float32
	Bids        []Bid
	Ledger      []Transaction
}
