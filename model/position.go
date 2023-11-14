package model

import (
	"gorm.io/gorm"
	"time"
)

type Position struct {
	gorm.Model
	Date      time.Time `json:"date"`
	Name      string    `json:"name"`
	Quote     int       `json:"quote"`
	Cost      float32   `json:"cost"`
	StockCode string    `json:"stock_code"`

	Stock `json:"stock" gorm:"references:Code;-"`
}

func (position *Position) getMarketValue() float32 {
	return float32(position.Quote) * position.ClosePrice
}

func (position *Position) getTotalCost() float32 {
	return float32(position.Quote) * position.Cost
}

func (position *Position) getProfit() float32 {
	return (position.ClosePrice - position.Cost) * float32(position.Quote)
}

func (position *Position) Read(db *gorm.DB, id int) error {
	result := db.Take(position, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (position *Position) Create(db *gorm.DB) error {
	result := db.Create(position)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (position *Position) Update(db *gorm.DB) error {
	result := db.Updates(position)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (position *Position) Delete(db *gorm.DB) error {
	result := db.Delete(position)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
