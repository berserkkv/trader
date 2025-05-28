package model

import (
	"github.com/berserkkv/trader/model/enum/order"
	"github.com/berserkkv/trader/model/enum/symbol"
	"github.com/berserkkv/trader/model/enum/timeframe"
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
	Leverage          float64       `gorm:"not null" json:"leverage"`
}

type CreateBotRequest struct {
	Symbol     string  `json:"symbol" binding:"required"`
	Strategy   string  `json:"strategy" binding:"required"`
	Timeframe  string  `json:"timeframe" binding:"required"`
	Capital    float64 `json:"capital" binding:"required"`
	Leverage   float64 `json:"leverage" binding:"required"`
	TakeProfit float64 `json:"takeProfit" binding:"required"`
	StopLoss   float64 `json:"stopLoss" binding:"required"`
}

type Statistics struct {
	ID           uint                `gorm:"primaryKey"`
	StrategyName string              `gorm:"not null" json:"strategy"`
	Symbol       symbol.Symbol       `gorm:"not null" json:"symbol"`
	Timeframe    timeframe.Timeframe `gorm:"not null" json:"timeframe"`
	Balance      float64             `gorm:"not null" json:"balance"`
	Time         time.Time           `gorm:"not null" json:"time"`
}
