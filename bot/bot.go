package bot

import (
	"errors"
	"github.com/berserkkv/trader/model/enum/order"
	"github.com/berserkkv/trader/model/enum/symbol"
	"github.com/berserkkv/trader/model/enum/timeframe"
	"github.com/berserkkv/trader/service/connector"
	"github.com/berserkkv/trader/strategy"
	"log/slog"
)

type Bot struct {
	Id          int64               `gorm:"primaryKey"`
	Name        string              `gorm:"not null;unique"`
	Symbol      symbol.Symbol       `gorm:"not null;check:name <> ''"`
	IsNotActive bool                `gorm:"default:false"`
	TimeFrame   timeframe.Timeframe `gorm:"not null"`

	StrategyName string            `gorm:"not null"`
	Strategy     strategy.Strategy `gorm:"-"` // ✅ Skip interface

	InitialCapital float64 `gorm:"not null"`
	CurrentCapital float64 `gorm:"not null"`

	Win         int64
	Loss        int64
	TotalTrades int64

	CurrentWinsStreak int64
	CurrentLossStreak int64
	MaxWinsStreak     int64
	MaxLossStreak     int64

	InPos        bool `gorm:"default:false"` // position is not open
	PositionType order.Command
	Asset        float64
	EntryPrice   float64
	StopLoss     float64
	TakeProfit   float64
}

func NewBot(timeframe timeframe.Timeframe, st strategy.Strategy, smb symbol.Symbol) *Bot {
	name := st.Name() + "_" + string(timeframe) + "_" + string(smb)
	return &Bot{
		Name:           name,
		Symbol:         smb,
		TimeFrame:      timeframe,
		StrategyName:   st.Name(),
		Strategy:       st,
		InitialCapital: 1000,
	}

}

func (b *Bot) OpenPosition(command order.Command) error {
	if err := b.CanOpenPosition(); err != nil {
		return err
	}
	price := connector.GetPrice(b.Symbol)

	b.Asset = b.CurrentCapital / price
	b.EntryPrice = price

	stopLossPercent := 0.002    // 0.2%
	takeProfitPercent := 0.0025 // 0.25%

	if command == order.LONG {
		b.StopLoss = price * (1 - stopLossPercent)
		b.TakeProfit = price * (1 + takeProfitPercent)
	} else {
		b.StopLoss = price * (1 + stopLossPercent)
		b.TakeProfit = price * (1 - takeProfitPercent)
	}

	b.InPos = true
	b.PositionType = command

	slog.Info("Position opened",
		"name", b.Name,
		"PositionType", b.PositionType,
		"entryPrice", b.EntryPrice,
		"stopLoss", b.StopLoss,
		"takeProfit", b.TakeProfit,
		"asset", b.Asset,
	)

	return nil
}

func (b *Bot) CanOpenPosition() error {
	if b.IsNotActive {
		slog.Debug("bot can't open position, bot not active", "name", b.Name)
		return errors.New("bot can't open position, bot not active")
	}

	if b.InPos {
		slog.Debug("bot is already in open position", "name", b.Name)
		return errors.New("bot is already in open position")
	}

	if b.CurrentCapital <= 10 {
		slog.Debug("bot can't open position, capital not enough", "name", b.Name)
		return errors.New("bot can't open position, capital not enough")
	}

	return nil
}
