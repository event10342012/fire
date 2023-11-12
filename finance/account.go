package finance

import (
	"gorm.io/gorm"
	"time"
)

type TransactionType int

const (
	Expense TransactionType = 1
	Revenue TransactionType = 2
)

type Transaction struct {
	gorm.Model
	Date            time.Time       `json:"date"`
	Name            string          `json:"name"`
	Amount          int             `json:"amount"`
	TransactionType TransactionType `json:"transaction_type"`
	Tags            string          `json:"tags"`
	Category        string          `json:"category"`
	UserID          uint            `json:"user_id"`
}

func (transaction *Transaction) Read(db *gorm.DB, id int) error {
	result := db.Take(transaction, id)

	if result.Error != nil {
		return nil
	}
	return nil
}

func (transaction *Transaction) Create(db *gorm.DB) error {
	result := db.Create(transaction)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (transaction *Transaction) Update(db *gorm.DB) error {
	result := db.Model(transaction).Updates(transaction)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (transaction *Transaction) Delete(db *gorm.DB) error {
	result := db.Delete(transaction)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
