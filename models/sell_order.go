package models

type SellOrderState string

const (
	UNLOCKED  SellOrderState = "Unlocked"
	LOCKED                   = "Locked"
	COMMITTED                = "Committed"
	REVERSED                 = "Reversed"
)

type SellOrder struct {
	Id             int32
	Invoice        Invoice
	SellerWants    float32
	Bids           []Bid
	Ledger         []Transaction
	SellOrderState SellOrderState
}
