package model

import "time"

type User struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
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

type Position struct {
	ID          int       `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	IsOpen      bool      `json:"isOpen,omitempty"`
	StartedTime time.Time `json:"startedTime,omitempty"`
	ClosedTime  time.Time `json:"closedTime,omitempty"`
	CreatedTime time.Time `json:"createdTime,omitempty"`
	UpdatedTime time.Time `json:"updatedTime,omitempty"`
	Symbol      string    `json:"symbol" binding:"required"`
	Side        string    `json:"side,omitempty"`
	EntryPrice  float64   `json:"entryPrice,omitempty"`
	ExitPrice   float64   `json:"exitPrice,omitempty"`
	Quantity    float64   `json:"quantity,omitempty"`
	ProfitLoss  float64   `json:"profitLoss,omitempty"`
	StopLoss    float64   `json:"stopLoss,omitempty"`
	TakeProfit  float64   `json:"takeProfit,omitempty"`
}

type Order struct {
	ID        int       `json:"id,omitempty"`
	Symbol    string    `json:"symbol,omitempty"`
	Side      string    `json:"side,omitempty"`
	Type      string    `json:"type,omitempty"`
	Price     float64   `json:"price,omitempty"`
	Quantity  float64   `json:"quantity,omitempty"`
	Status    string    `json:"status,omitempty"` // "OPEN", "FILLED", "CANCELLED"
	CreatedAt time.Time `json:"createdAt,omitempty"`
	FilledAt  time.Time `json:"filledAt,omitempty"`
}
