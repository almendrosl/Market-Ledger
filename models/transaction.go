package models

import "time"

type TransactionType string

const ZeroCreditDebit = 0.0

const (
	CASH            TransactionType = "Cash"
	OWN_INVOICE                     = "OwnInvoice"
	RESERVED                        = "Reserved"
	EXPECTED_PAID                   = "ExpectedPaid"
	EXPECTED_PROFIT                 = "ExpectedProfit"
	CAPITAL                         = "Capital"
	OWES_INVOICE                    = "OwesInvoice"
	LOSS                            = "Loss"
)

type Transaction struct {
	Id              int32
	Date            time.Time
	TransactionType TransactionType
	Details         string
	Debit           float32
	Credit          float32
	Customer        Customer
	SellOrder       SellOrder
}
