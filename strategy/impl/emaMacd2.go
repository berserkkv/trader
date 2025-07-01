package strategyImpl

import (
	"fmt"
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/model/enum/order"
	"github.com/berserkkv/trader/service/tools"
	"github.com/berserkkv/trader/ta"
)

type EmaMacd2 struct{}

func (*EmaMacd2) Name() string {
	return "EmaMacd2"
}

func (s *EmaMacd2) Run(candles []model.Candle) (order.Command, string) {
	if len(candles) == 0 {
		return order.WAIT, "candles is empty"
	}
	closePrices := tools.GetClosePrices(candles)
	ema200 := ta.EMA(closePrices, 200)
	macdLine, signalLine, _ := ta.MACDSlice(closePrices)
	macdPrev := macdLine[len(macdLine)-2]
	signalPrev := signalLine[len(signalLine)-2]
	macd := macdLine[len(macdLine)-1]
	signal := signalLine[len(signalLine)-1]
	price := candles[len(candles)-1].Close

	info := fmt.Sprintf("macd=%.2f, signal=%.2f, macdPrev=%.2f, signalPrev=%.2f, ema200=%.2f, price=%.2f", macd, signal, macdPrev, signalPrev, ema200, price)

	if price > ema200 && macdPrev < signalPrev && macd > signal {
		return order.LONG, "Long " + info
	}
	if price < ema200 && macdPrev > signalPrev && macd < signal {
		return order.SHORT, "Short " + info
	}
	return order.WAIT, info
}
