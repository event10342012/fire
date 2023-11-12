package finance

import (
	"gorm.io/gorm"
	"time"
)

type Stock struct {
	gorm.Model
	TxnDate    time.Time `json:"date"`
	Code       string    `json:"code"`
	OpenPrice  float32   `json:"open_price"`
	HighPrice  float32   `json:"high_price"`
	LowPrice   float32   `json:"low_price"`
	ClosePrice float32   `json:"close_price"`
}
