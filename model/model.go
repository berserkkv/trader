package model

import (
	"github.com/berserkkv/trader/model/enum/order"
	"github.com/berserkkv/trader/model/enum/symbol"
	"time"
)

type Candle struct {
	Time   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64
}

type HeikinAshi struct {
	Open, High, Low, Close float64
	Color                  string
}

type Order struct {
	ID                int64         `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	Symbol            symbol.Symbol `gorm:"not null" json:"symbol"`
	Type              order.Command `gorm:"not null" json:"type"`
	BotID             int64         `gorm:"not null" json:"botId"`
	EntryPrice        float64       `gorm:"not null" json:"entryPrice"`
	ExitPrice         float64       `gorm:"not null" json:"exitPrice"`
	Quantity          float64       `gorm:"not null" json:"quantity"`
	ProfitLoss        float64       `gorm:"not null" json:"profitLoss"`
	ProfitLossPercent float64       `gorm:"not null" json:"profitLossPercent"`
	CreatedTime       time.Time     `gorm:"not null" json:"createdTime"`
	ClosedTime        time.Time     `gorm:"not null" json:"closedTime"`
	Fee               float64       `gorm:"not null" json:"fee"`
}

type CreateBotRequest struct {
	Symbol    string  `json:"symbol" binding:"required"`
	Strategy  string  `json:"strategy" binding:"required"`
	Timeframe string  `json:"timeframe" binding:"required"`
	Capital   float64 `json:"capital" binding:"required"`
}
