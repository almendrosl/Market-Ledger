package models

import "time"

type TransactionType string

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

type TransactionDCType string

const (
	DEBIT  TransactionDCType = "debit"
	CREDIT TransactionDCType = "credit"
)

type Transaction struct {
	Id                int32
	Date              time.Time
	TransactionType   TransactionType
	Details           string
	TransactionDCType TransactionDCType
	Value             float32
	Customer          Customer
	SellOrder         SellOrder
}
