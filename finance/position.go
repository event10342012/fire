package finance

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

	Stock `json:"stock" gorm:"references:Code"`
}

func (p *Position) getMarketValue() float32 {
	return float32(p.Quote) * p.ClosePrice
}

func (p *Position) getTotalCost() float32 {
	return float32(p.Quote) * p.Cost
}

func (p *Position) getProfit() float32 {
	return (p.ClosePrice - p.Cost) * float32(p.Quote)
}
