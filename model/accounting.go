package model

import (
	"gorm.io/gorm"
	"time"
)

type AccountingTransactionType int8

const (
	Expense AccountingTransactionType = 1
	Revenue AccountingTransactionType = 2
)

type AccountingTransaction struct {
	gorm.Model
	Date            time.Time                 `json:"date"`
	Note            string                    `json:"note"`
	Amount          int                       `json:"amount"`
	Currency        string                    `json:"currency"`
	TransactionType AccountingTransactionType `json:"transaction_type"`
	Tags            string                    `json:"tags"`
	Category        string                    `json:"category"`
	UserID          uint                      `json:"user_id"`
}
