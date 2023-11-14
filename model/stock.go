package model

import (
	"gorm.io/gorm"
	"time"
)

type Stock struct {
	TxnDate    time.Time `json:"date" gorm:"primaryKey"`
	Code       string    `json:"code" gorm:"primaryKey"`
	OpenPrice  float32   `json:"open_price"`
	HighPrice  float32   `json:"high_price"`
	LowPrice   float32   `json:"low_price"`
	ClosePrice float32   `json:"close_price"`
	Volume     int64     `json:"volume"`
}

func (stock *Stock) Read(db *gorm.DB, id int) error {
	result := db.Take(stock, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (stock *Stock) Create(db *gorm.DB) error {
	result := db.Create(stock)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (stock *Stock) Update(db *gorm.DB) error {
	result := db.Updates(stock)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (stock *Stock) Delete(db *gorm.DB) error {
	result := db.Delete(stock)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
