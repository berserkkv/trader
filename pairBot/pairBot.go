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
	CurrentCapital1           float64             `gorm:"not null" json:"currentCapital1"`
	CurrentCapital2           float64             `gorm:"not null" json:"currentCapital2"`
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
	OrderFee1                 float64             `json:"orderFee1"`
	OrderFee2                 float64             `json:"orderFee2"`
	Pnl1                      float64             `json:"pnl1"`
	Pnl2                      float64             `json:"pnl2"`
	Roe1                      float64             `json:"roe1"`
	Roe2                      float64             `json:"roe2"`
}

func NewPairBot(smb1, smb2 symbol.Symbol, timeframe timeframe.Timeframe, capital, leverage, takeProfit, stopLoss float64) *PairBot {
	name := string(smb1) + "_" + string(smb2) + "_" + string(timeframe)
	return &PairBot{
		Name:            name,
		Symbol1:         smb1,
		Symbol2:         smb2,
		Timeframe:       timeframe,
		Connector:       connector.BinanceConnector{},
		CurrentCapital1: capital,
		CurrentCapital2: capital,
		Leverage:        leverage,
		TakeProfit:      takeProfit,
		StopLoss:        stopLoss,
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

	capital1 := b.CurrentCapital1
	capital2 := b.CurrentCapital2

	b.CurrentCapital1 -= capital1
	b.CurrentCapital2 -= capital2

	fee1 := calculator.CalculateMakerFee(capital1)
	fee2 := calculator.CalculateMakerFee(capital2)

	capital1 -= fee1
	capital2 -= fee2

	b.OrderCapitalWithLeverage1 = b.Leverage * capital1
	b.OrderCapitalWithLeverage2 = b.Leverage * capital2

	now := time.Now()
	b.OrderQuantity1 = calculator.CalculateBuyQuantity(price1, b.OrderCapitalWithLeverage1)
	b.OrderQuantity2 = calculator.CalculateBuyQuantity(price2, b.OrderCapitalWithLeverage2)

	b.OrderEntryPrice1 = price1
	b.OrderEntryPrice2 = price2
	b.OrderCapital1 = capital1
	b.OrderCapital2 = capital2
	b.InPos = true
	b.OrderCreatedTime = now
	b.OrderScannedTime = now
	b.OrderFee1 = fee1
	b.OrderFee2 = fee2

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

	b.OrderFee1 += fee1
	b.OrderFee2 += fee2

	b.CurrentCapital1 += b.OrderCapital1 + pnl1
	b.CurrentCapital2 += b.OrderCapital2 + pnl2

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
		Fee1:               b.OrderFee1,
		Fee2:               b.OrderFee2,
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
	b.OrderFee1 = 0
	b.OrderScannedTime = time.Time{}
	b.Pnl1 = 0
	b.Roe1 = 0

	return closedOrder, nil
}

func (b *PairBot) UpdatePnlAndRoe(curPrice1, curPrice2 float64) {
	b.Roe1 = calculator.CalculateRoe(b.OrderEntryPrice1, curPrice1, b.Leverage, b.OrderType1)
	b.Pnl1 = calculator.CalculatePNL(curPrice1, b.OrderCapitalWithLeverage1, b.OrderQuantity1, b.OrderType1)

	b.Roe2 = calculator.CalculateRoe(b.OrderEntryPrice2, curPrice2, b.Leverage, b.OrderType2)
	b.Pnl2 = calculator.CalculatePNL(curPrice2, b.OrderCapitalWithLeverage2, b.OrderQuantity2, b.OrderType2)

}

func (b *PairBot) ShouldOpenPosition() order.Command {
	price1, price2 := b.GetKlines()
	zScore := calculator.CalculatePairTradingSpread(price1, price2)
	b.ZScore = zScore

	if zScore >= 2 {
		return order.SHORT
	} else if zScore <= -2 {
		return order.LONG
	}
	return order.WAIT
}

func (b *PairBot) ShouldClosePosition() bool {
	// if orderType is short then zScore was more than 2
	if b.OrderType1 == order.SHORT {
		if b.ZScore <= -2 {
			return true
		}
	} else {
		if b.ZScore >= 2 {
			return true
		}
	}
	return false
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

	if b.CurrentCapital1+b.CurrentCapital2 <= 90 {
		slog.Debug("bot can't open position, capital not enough", "name", b.Name)
		return errors.New("bot can't open position, capital not enough")
	}

	return nil
}

func (b *PairBot) String() string {
	return fmt.Sprintf("{Name: %s, InPos: %t, ZScore: %.2f, Capital1: %.2f, Capital2: %.2f}", b.Name, b.InPos, b.ZScore, b.CurrentCapital1, b.CurrentCapital2)
}

func (b *PairBot) GetKlines() ([]float64, []float64) {
	candles1 := b.Connector.GetCandles(b.Symbol1, b.Timeframe, 200)
	candles2 := b.Connector.GetCandles(b.Symbol2, b.Timeframe, 200)

	closedPrice1 := make([]float64, len(candles1))
	closedPrice2 := make([]float64, len(candles2))

	for i := 0; i < len(candles1); i++ {
		closedPrice1[i] = candles1[i].Close
		closedPrice2[i] = candles2[i].Close
	}
	return closedPrice1, closedPrice2
}
