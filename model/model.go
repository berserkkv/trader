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
}

type Position struct {
	ID          int       // Unique identifier
	IsOpen      bool      // Whether position is currently open
	StartedTime time.Time // position started time
	ClosedTime  time.Time // position closed time
	CreatedTime time.Time // init time
	UpdatedTime time.Time
	Symbol      string  // e.g., "BTC/USDT" or "AAPL"
	Side        string  // "BUY" or "SELL"
	EntryPrice  float64 // Price at which position was entered
	ExitPrice   float64 // Price at which position was exited (0 if open)
	Quantity    float64 // Quantity traded
	ProfitLoss  float64 // Realized or unrealized PnL
	StopLoss    float64 // Optional stop-loss level
	TakeProfit  float64 // Optional take-profit level
}

type Order struct {
	ID        int
	Symbol    string
	Side      string
	Type      string
	Price     float64
	Quantity  float64
	Status    string    // "OPEN", "FILLED", "CANCELLED"
	CreatedAt time.Time // Order creation time
	FilledAt  time.Time // Order filled time
}
