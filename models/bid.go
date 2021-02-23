package models

type Bid struct {
	Id        int
	Size      float32
	Amount    float32
	Investor  Investor
	SellOrder SellOrder
}
