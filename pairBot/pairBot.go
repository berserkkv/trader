package pairBot

import (
	"errors"
	"fmt"
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/model/enum/order"
	"github.com/berserkkv/trader/model/enum/symbol"
	"github.com/berserkkv/trader/model/enum/timeframe"
	"github.com/berserkkv/trader/service/calculator"
	"github.com/berserkkv/trader/service/connector"
	"log/slog"
	"time"
)

type PairBot struct {
	Id                        int64               `gorm:"primaryKey" json:"id"`
	Name                      string              `gorm:"not null;unique" json:"name"`
	Symbol1                   symbol.Symbol       `gorm:"not null;check:name <> ''" json:"symbol1"`
	Symbol2                   symbol.Symbol       `gorm:"not null;check:name <> ''" json:"symbol2"`
	IsNotActive               bool                `gorm:"default:false" json:"isNotActive"`
	Timeframe                 timeframe.Timeframe `gorm:"not null" json:"timeframe"`
	Connector                 connector.Connector `gorm:"-" json:"-"`
	CurrentCapital            float64             `gorm:"not null" json:"currentCapital"`
	LastScanned               time.Time           `gorm:"not null" json:"lastScanned"`
	TotalWins                 int64               `json:"totalWins"`
	TotalLosses               int64               `json:"totalLosses"`
	Leverage                  float64             `json:"leverage"`
	TakeProfit                float64             `json:"takeProfit"`
	StopLoss                  float64             `json:"stopLoss"`
	ZScore                    float64             `json:"zScore"`
	InPos                     bool                `gorm:"default:false" json:"inPos"`
	OrderType1                order.Command       `json:"orderType1"`
	OrderType2                order.Command       `json:"orderType2"`
	OrderCreatedTime          time.Time           `json:"orderCreatedTime1"`
	OrderScannedTime          time.Time           `json:"orderScannedTime2"`
	OrderQuantity1            float64             `json:"orderQuantity1"`
	OrderQuantity2            float64             `json:"orderQuantity2"`
	OrderCapital1             float64             `json:"orderCapital1"`
	OrderCapital2             float64             `json:"orderCapital2"`
	OrderCapitalWithLeverage1 float64             `json:"orderCapitalWithLeverage1"`
	OrderCapitalWithLeverage2 float64             `json:"orderCapitalWithLeverage2"`
	OrderEntryPrice1          float64             `json:"orderEntryPrice1"`
	OrderEntryPrice2          float64             `json:"orderEntryPrice2"`
	OrderStopLoss1            float64             `json:"orderStopLoss1"`
	OrderStopLoss2            float64             `json:"orderStopLoss2"`
	OrderTakeProfit1          float64             `json:"orderTakeProfit1"`
	OrderTakeProfit2          float64             `json:"orderTakeProfit2"`
	OrderFee                  float64             `json:"orderFee"`
	Pnl1                      float64             `json:"pnl1"`
	Pnl2                      float64             `json:"pnl2"`
	Roe1                      float64             `json:"roe1"`
	Roe2                      float64             `json:"roe2"`
}

func NewPairBot(smb1, smb2 symbol.Symbol, timeframe timeframe.Timeframe, capital, leverage, takeProfit, stopLoss float64) *PairBot {
	name := string(smb1) + "_" + string(smb2) + "_" + string(timeframe)
	return &PairBot{
		Name:           name,
		Symbol1:        smb1,
		Symbol2:        smb2,
		Timeframe:      timeframe,
		Connector:      connector.BinanceConnector{},
		CurrentCapital: capital,
		Leverage:       leverage,
		TakeProfit:     takeProfit,
		StopLoss:       stopLoss,
	}

}

func (b *PairBot) OpenPosition(cmd1 order.Command) error {
	if err := b.CanOpenPosition(); err != nil {
		return err
	}
	price1 := b.Connector.GetPrice(b.Symbol1)
	price2 := b.Connector.GetPrice(b.Symbol2)
	var cmd2 order.Command

	if cmd1 == order.LONG {
		cmd2 = order.SHORT
	} else {
		cmd2 = order.LONG
	}

	b.OrderType1 = cmd1
	b.OrderType2 = cmd2

	capital := b.CurrentCapital

	b.CurrentCapital -= capital

	fee := calculator.CalculateMakerFee(capital)

	capital -= fee

	b.OrderCapitalWithLeverage1 = b.Leverage * capital / 2
	b.OrderCapitalWithLeverage2 = b.Leverage * capital / 2

	now := time.Now()
	b.OrderQuantity1 = calculator.CalculateBuyQuantity(price1, b.OrderCapitalWithLeverage1)
	b.OrderQuantity2 = calculator.CalculateBuyQuantity(price2, b.OrderCapitalWithLeverage2)

	b.OrderEntryPrice1 = price1
	b.OrderEntryPrice2 = price2
	b.OrderCapital1 = capital / 2
	b.OrderCapital2 = capital / 2
	b.InPos = true
	b.OrderCreatedTime = now
	b.OrderScannedTime = now
	b.OrderFee = fee

	return nil
}

func (b *PairBot) ClosePosition(curPrice1, curPrice2 float64) (model.PairOrder, error) {

	var pnl1 float64
	var pnlPercent1 float64

	var pnl2 float64
	var pnlPercent2 float64

	fee1 := calculator.CalculateMakerFee(b.OrderCapital1)
	b.OrderCapitalWithLeverage1 -= fee1
	b.OrderCapital1 -= fee1

	fee2 := calculator.CalculateMakerFee(b.OrderCapital2)
	b.OrderCapitalWithLeverage2 -= fee2
	b.OrderCapital2 -= fee2

	pnl1 = calculator.CalculatePNL(curPrice1, b.OrderCapitalWithLeverage1, b.OrderQuantity1, b.OrderType1)
	pnlPercent1 = calculator.CalculateRoe(b.OrderEntryPrice1, curPrice1, b.Leverage, b.OrderType1)

	pnl2 = calculator.CalculatePNL(curPrice2, b.OrderCapitalWithLeverage2, b.OrderQuantity2, b.OrderType2)
	pnlPercent2 = calculator.CalculateRoe(b.OrderEntryPrice2, curPrice2, b.Leverage, b.OrderType2)

	if pnl1+pnl2 > 0 {
		b.TotalWins++
	} else {
		b.TotalLosses++
	}

	b.OrderFee += fee1 + fee2
	totalCapital := b.OrderCapital1 + b.OrderCapital2 + pnl1 + pnl2
	b.CurrentCapital += totalCapital

	closedOrder := model.PairOrder{
		Symbol1:            b.Symbol1,
		Symbol2:            b.Symbol2,
		Type1:              b.OrderType1,
		Type2:              b.OrderType2,
		BotID:              b.Id,
		EntryPrice1:        b.OrderEntryPrice1,
		EntryPrice2:        b.OrderEntryPrice2,
		ExitPrice1:         curPrice1,
		ExitPrice2:         curPrice2,
		Quantity1:          b.OrderQuantity1,
		Quantity2:          b.OrderQuantity2,
		ProfitLoss1:        pnl1,
		ProfitLoss2:        pnl2,
		ProfitLossPercent1: pnlPercent1,
		ProfitLossPercent2: pnlPercent2,
		CreatedTime:        b.OrderCreatedTime,
		ClosedTime:         time.Now(),
		Fee:                b.OrderFee,
		Leverage:           b.Leverage,
	}

	b.InPos = false
	b.OrderEntryPrice1 = 0
	b.OrderStopLoss1 = 0
	b.OrderTakeProfit1 = 0
	b.OrderType1 = ""
	b.OrderCapital1 = 0
	b.OrderCapitalWithLeverage1 = 0
	b.OrderCreatedTime = time.Time{}
	b.OrderQuantity1 = 0
	b.OrderFee = 0
	b.OrderScannedTime = time.Time{}
	b.Pnl1 = 0
	b.Roe1 = 0
	b.OrderEntryPrice2 = 0
	b.OrderStopLoss2 = 0
	b.OrderTakeProfit2 = 0
	b.OrderType2 = ""
	b.OrderCapital2 = 0
	b.OrderCapitalWithLeverage2 = 0
	b.OrderQuantity2 = 0
	b.Pnl2 = 0
	b.Roe2 = 0

	return closedOrder, nil
}

func (b *PairBot) UpdatePnlAndRoe(curPrice1, curPrice2 float64) {
	b.Roe1 = calculator.CalculateRoe(b.OrderEntryPrice1, curPrice1, b.Leverage, b.OrderType1)
	b.Pnl1 = calculator.CalculatePNL(curPrice1, b.OrderCapitalWithLeverage1, b.OrderQuantity1, b.OrderType1)

	b.Roe2 = calculator.CalculateRoe(b.OrderEntryPrice2, curPrice2, b.Leverage, b.OrderType2)
	b.Pnl2 = calculator.CalculatePNL(curPrice2, b.OrderCapitalWithLeverage2, b.OrderQuantity2, b.OrderType2)

}

func (b *PairBot) ShouldOpenPosition() order.Command {
	b.UpdateZScore()

	if b.ZScore >= 2 {
		return order.SHORT
	} else if b.ZScore <= -2 {
		return order.LONG
	}
	return order.WAIT
}

func (b *PairBot) ShouldClosePosition() bool {
	// if orderType is short then zScore was more than 2
	if b.OrderType1 == order.SHORT {
		if b.ZScore <= 1 {
			return true
		}
	} else {
		if b.ZScore >= -1 {
			return true
		}
	}
	return false
}

func (b *PairBot) UpdateZScore() {
	price1, price2 := b.GetKlines()
	if price1 == nil || price2 == nil {
		slog.Error("Get klines from API", "price1", price1, "price2", price2)
		return
	}
	b.ZScore = calculator.CalculatePairTradingSpread(price1, price2)
}

func (b *PairBot) CanOpenPosition() error {
	if b.IsNotActive {
		slog.Debug("bot can't open position, bot not active", "name", b.Name)
		return errors.New("bot can't open position, bot not active")
	}

	if b.InPos {
		slog.Debug("bot is already in open position", "name", b.Name)
		return errors.New("bot is already in open position")
	}

	if b.CurrentCapital <= 85 {
		slog.Debug("bot can't open position, capital not enough", "name", b.Name)
		return errors.New("bot can't open position, capital not enough")
	}

	return nil
}

func (b *PairBot) String() string {
	return fmt.Sprintf("{Name: %s, InPos: %t, ZScore: %.2f, Capital: %.2f}", b.Name, b.InPos, b.ZScore, b.CurrentCapital)
}

func (b *PairBot) GetKlines() ([]float64, []float64) {
	candles1 := b.Connector.GetCandles(b.Symbol1, b.Timeframe, 200)
	candles2 := b.Connector.GetCandles(b.Symbol2, b.Timeframe, 200)

	if candles1 == nil || candles2 == nil {
		return nil, nil
	}

	closedPrice1 := make([]float64, len(candles1))
	closedPrice2 := make([]float64, len(candles2))

	for i := 0; i < len(candles1); i++ {
		closedPrice1[i] = candles1[i].Close
		closedPrice2[i] = candles2[i].Close
	}
	return closedPrice1, closedPrice2
}
