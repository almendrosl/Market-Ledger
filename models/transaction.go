package models

import "time"

type TransactionType string

const (
	Cash           TransactionType = "Cash"
	OwnInvoice                     = "OwnInvoice"
	Reserved                       = "Reserved"
	ExpectedPaid                   = "ExpectedPaid"
	ExpectedProfit                 = "ExpectedProfit"
	Capital                        = "Capital"
	OwesInvoice                    = "OwesInvoice"
	Loss                           = "Loss"
)


type Transaction struct {
	Id              int
	Date            time.Time
	TransactionType TransactionType
	Details         string
	Debt            float32
	Credit          float32
	Customer        Customer
	SellOrder       SellOrder
}
