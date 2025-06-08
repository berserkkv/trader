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

type PairOrder struct {
	ID                 int64         `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	Symbol1            symbol.Symbol `gorm:"not null" json:"symbol1"`
	Symbol2            symbol.Symbol `gorm:"not null" json:"symbol2"`
	Type1              order.Command `gorm:"not null" json:"type1"`
	Type2              order.Command `gorm:"not null" json:"type2"`
	BotID              int64         `gorm:"not null" json:"botId"`
	EntryPrice1        float64       `gorm:"not null" json:"entryPrice1"`
	EntryPrice2        float64       `gorm:"not null" json:"entryPrice2"`
	ExitPrice1         float64       `gorm:"not null" json:"exitPrice1"`
	ExitPrice2         float64       `gorm:"not null" json:"exitPrice2"`
	Quantity1          float64       `gorm:"not null" json:"quantity1"`
	Quantity2          float64       `gorm:"not null" json:"quantity2"`
	ProfitLoss1        float64       `gorm:"not null" json:"profitLoss1"`
	ProfitLoss2        float64       `gorm:"not null" json:"profitLoss2"`
	ProfitLossPercent1 float64       `gorm:"not null" json:"profitLossPercent1"`
	ProfitLossPercent2 float64       `gorm:"not null" json:"profitLossPercent2"`
	CreatedTime        time.Time     `gorm:"not null" json:"createdTime"`
	ClosedTime         time.Time     `gorm:"not null" json:"closedTime"`
	Fee1               float64       `gorm:"not null" json:"fee1"`
	Fee2               float64       `gorm:"not null" json:"fee2"`
	Leverage           float64       `gorm:"not null" json:"leverage"`
}

type Statistics struct {
	Pnl  float64   `json:"pnl"`
	Time time.Time `json:"time"`
}
