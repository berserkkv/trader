package bot

import (
	"errors"
	"fmt"
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/model/enum/order"
	"github.com/berserkkv/trader/model/enum/symbol"
	"github.com/berserkkv/trader/model/enum/timeframe"
	"github.com/berserkkv/trader/service/calculator"
	"github.com/berserkkv/trader/service/connector"
	"github.com/berserkkv/trader/strategy"
	"log/slog"
	"time"
)

type Bot struct {
	Id                       int64               `gorm:"primaryKey" json:"id"`
	Name                     string              `gorm:"not null;unique" json:"name"`
	Symbol                   symbol.Symbol       `gorm:"not null;check:name <> ''" json:"symbol"`
	IsNotActive              bool                `gorm:"default:false" json:"isNotActive"`
	TimeFrame                timeframe.Timeframe `gorm:"not null" json:"timeFrame"`
	StrategyName             string              `gorm:"not null" json:"strategyName"`
	Strategy                 strategy.Strategy   `gorm:"-" json:"-"` // Skip interface for DB and JSON
	InitialCapital           float64             `gorm:"not null" json:"initialCapital"`
	CurrentCapital           float64             `gorm:"not null" json:"currentCapital"`
	LastScanned              time.Time           `gorm:"not null" json:"lastScanned"`
	TotalWins                int64               `json:"totalWins"`
	TotalLosses              int64               `json:"totalLosses"`
	TotalTrades              int64               `json:"totalTrades"`
	CurrentWinsStreak        int64               `json:"currentWinsStreak"`
	CurrentLossStreak        int64               `json:"currentLossStreak"`
	MaxWinsStreak            int64               `json:"maxWinsStreak"`
	MaxLossStreak            int64               `json:"maxLossStreak"`
	Leverage                 float64             `json:"leverage"`
	TakeProfit               float64             `json:"takeProfit"`
	StopLoss                 float64             `json:"stopLoss"`
	InPos                    bool                `gorm:"default:false" json:"inPos"`
	OrderType                order.Command       `json:"orderType"`
	OrderCreatedTime         time.Time           `json:"orderCreatedTime"`
	OrderScannedTime         time.Time           `json:"orderScannedTime"`
	OrderQuantity            float64             `json:"orderQuantity"`
	OrderCapital             float64             `json:"orderCapital"`
	OrderCapitalWithLeverage float64             `json:"orderCapitalWithLeverage"`
	OrderEntryPrice          float64             `json:"orderEntryPrice"`
	OrderStopLoss            float64             `json:"orderStopLoss"`
	OrderTakeProfit          float64             `json:"orderTakeProfit"`
	OrderFee                 float64             `json:"orderFee"`
	Pnl                      float64             `json:"pnl"`
	Roe                      float64             `json:"roe"`
}

func NewBot(timeframe timeframe.Timeframe, st strategy.Strategy, smb symbol.Symbol, capital, leverage, takeProfit, stopLoss float64) *Bot {
	name := st.Name() + "_" + string(timeframe) + "_" + string(smb)
	return &Bot{
		Name:           name,
		Symbol:         smb,
		TimeFrame:      timeframe,
		StrategyName:   st.Name(),
		Strategy:       st,
		InitialCapital: capital,
		CurrentCapital: capital,
		Leverage:       leverage,
		TakeProfit:     takeProfit,
		StopLoss:       stopLoss,
	}

}

func (b *Bot) OpenPosition(command order.Command) error {
	if err := b.CanOpenPosition(); err != nil {
		return err
	}
	price := connector.GetPrice(b.Symbol)

	stopLossPercent := 0.5
	takeProfitPercent := 1.5

	if command == order.LONG {
		b.OrderStopLoss = calculator.CalculateStopLossWithPercent(price, stopLossPercent, false)
		b.OrderTakeProfit = calculator.CalculateTakeProfitWithPercent(price, takeProfitPercent, false)
	} else {
		b.OrderStopLoss = calculator.CalculateStopLossWithPercent(price, stopLossPercent, true)
		b.OrderTakeProfit = calculator.CalculateTakeProfitWithPercent(price, takeProfitPercent, true)
	}
	fee := calculator.CalculateMakerFee(b.CurrentCapital)

	b.CurrentCapital -= fee

	b.OrderCapitalWithLeverage = b.Leverage * b.CurrentCapital

	now := time.Now()
	b.OrderQuantity = calculator.CalculateBuyQuantity(price, b.OrderCapitalWithLeverage)
	b.OrderEntryPrice = price
	b.OrderCapital = b.CurrentCapital
	// b.CurrentCapital = 0
	b.InPos = true
	b.OrderType = command
	b.OrderCreatedTime = now
	b.OrderScannedTime = now
	b.OrderFee = fee

	slog.Info("Position opened",
		"name", b.Name,
		"OrderType", b.OrderType,
		"entryPrice", b.OrderEntryPrice,
		"stopLoss", b.OrderStopLoss,
		"takeProfit", b.OrderTakeProfit,
		"asset", b.OrderQuantity,
	)

	return nil
}

func (b *Bot) ClosePosition(curPrice float64) (model.Order, error) {
	var pnl float64
	var pnlPercent float64

	fee := calculator.CalculateMakerFee(b.OrderCapital)
	b.OrderCapitalWithLeverage -= fee
	b.OrderCapital -= fee

	if b.OrderType == order.LONG {
		pnl = calculator.CalculatePNLForLong(curPrice, b.OrderCapitalWithLeverage, b.OrderQuantity)
		pnlPercent = calculator.CalculateRoeForLong(b.OrderEntryPrice, curPrice, b.Leverage)
	} else if b.OrderType == order.SHORT {
		pnl = calculator.CalculatePNLForShort(curPrice, b.OrderCapitalWithLeverage, b.OrderQuantity)
		pnlPercent = calculator.CalculateRoeForShort(b.OrderEntryPrice, curPrice, b.Leverage)
	} else {
		return model.Order{}, fmt.Errorf("invalid order type")
	}

	b.calculateStatistics(pnl)

	b.OrderFee += fee

	b.CurrentCapital = b.OrderCapital + pnl

	closedOrder := model.Order{
		Symbol:            b.Symbol,
		Type:              b.OrderType,
		BotID:             b.Id,
		EntryPrice:        b.OrderEntryPrice,
		ExitPrice:         curPrice,
		Quantity:          b.OrderQuantity,
		ProfitLoss:        pnl,
		ProfitLossPercent: pnlPercent,
		CreatedTime:       b.OrderCreatedTime,
		ClosedTime:        time.Now(),
		Fee:               b.OrderFee,
		Leverage:          b.Leverage,
	}

	b.InPos = false
	b.OrderEntryPrice = 0
	b.OrderStopLoss = 0
	b.OrderTakeProfit = 0
	b.OrderType = ""
	b.OrderCapital = 0
	b.OrderCapitalWithLeverage = 0
	b.OrderCreatedTime = time.Time{}
	b.OrderQuantity = 0
	b.OrderFee = 0
	b.OrderScannedTime = time.Time{}
	b.Pnl = 0
	b.Roe = 0

	return closedOrder, nil
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

func (b *Bot) ShouldClosePosition(curPrice float64) bool {
	if b.OrderType == order.LONG {
		if curPrice >= b.OrderTakeProfit || curPrice <= b.OrderStopLoss {
			return true
		}
	} else {
		if curPrice <= b.OrderTakeProfit || curPrice >= b.OrderStopLoss {
			return true
		}
	}
	return false
}

func (b *Bot) calculateStatistics(pnl float64) {
	if pnl > 0 {
		b.TotalWins++
		if b.CurrentLossStreak > 0 {
			b.MaxLossStreak = max(b.MaxLossStreak, b.CurrentLossStreak)
			b.CurrentLossStreak = 0
		}
		b.CurrentWinsStreak++
		b.MaxWinsStreak = max(b.MaxWinsStreak, b.CurrentWinsStreak)

	} else {
		b.TotalLosses++
		if b.CurrentWinsStreak > 0 {
			b.MaxWinsStreak = max(b.MaxWinsStreak, b.CurrentWinsStreak)
			b.CurrentWinsStreak = 0
		}
		b.CurrentLossStreak++
		b.MaxLossStreak = max(b.MaxLossStreak, b.CurrentLossStreak)
	}
	b.TotalTrades++
}

func (b *Bot) calculateRoe(curPrice float64) float64 {
	if b.OrderType == order.LONG {
		return calculator.CalculateRoeForLong(b.OrderEntryPrice, curPrice, b.Leverage)
	} else if b.OrderType == order.SHORT {
		return calculator.CalculateRoeForShort(b.OrderEntryPrice, curPrice, b.Leverage)
	}
	return 0
}

func (b *Bot) calculatePnl(curPrice float64) float64 {
	if b.OrderType == order.LONG {
		return calculator.CalculatePNLForLong(curPrice, b.OrderCapitalWithLeverage, b.OrderQuantity)
	} else if b.OrderType == order.SHORT {
		return calculator.CalculatePNLForShort(curPrice, b.OrderCapitalWithLeverage, b.OrderQuantity)
	}
	return 0
}

func (b *Bot) UpdatePnlAndRoe(curPrice float64) {
	b.Roe = b.calculateRoe(curPrice)
	b.Pnl = b.calculatePnl(curPrice)
}

func (b *Bot) ShiftStopLoss() {
	realROE := b.Roe / b.Leverage

	if realROE >= 0.2 {
		// Convert to decimal before calculations
		pnlDecimal := realROE / 100.0

		// Shift stop loss to half the profit
		shift := pnlDecimal / 2.0

		newStopLoss := b.OrderEntryPrice * (1.0 + shift)
		if b.OrderType == order.LONG {
			if newStopLoss > b.StopLoss {
				b.OrderStopLoss = newStopLoss
			}
		} else if b.OrderType == order.SHORT {
			if newStopLoss < b.StopLoss {
				b.OrderStopLoss = newStopLoss
			}
		}
	}
}

func (b *Bot) String() string {
	return fmt.Sprintf("{Name: %s, InPos: %t, Capital: %.2f}", b.Name, b.InPos, b.CurrentCapital)
}
