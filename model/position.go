package model

import (
	"gorm.io/gorm"
	"time"
)

type AssetPosition struct {
	gorm.Model
	Date      time.Time `json:"date"`
	Name      string    `json:"name"`
	AssetType string    `json:"asset_type"`
	Quantity  int       `json:"quantity"`
	Cost      float32   `json:"cost"`
	AssetID   int       `json:"asset_id"`

	Stock `json:"stock" gorm:"references:Code;-"`
}

func (position *AssetPosition) getMarketValue() float32 {
	return float32(position.Quantity) * position.ClosePrice
}

func (position *AssetPosition) getTotalCost() float32 {
	return float32(position.Quantity) * position.Cost
}

func (position *AssetPosition) getProfit() float32 {
	return (position.ClosePrice - position.Cost) * float32(position.Quantity)
}
