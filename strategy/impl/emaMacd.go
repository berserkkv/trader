package strategyImpl

import (
	"fmt"
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/model/enum/order"
	"github.com/berserkkv/trader/service/tools"
	"github.com/berserkkv/trader/ta"
)

type EmaMacd struct{}

func (s *EmaMacd) Name() string {
	return "EmaMacd"
}

func (s *EmaMacd) Run(candles []model.Candle) (order.Command, string) {
	macdLine, signalLine, histogram := ta.MACDSlice(tools.GetClosePrices(candles))
	ema200 := ta.EMA(tools.GetClosePrices(candles), 200)

	n := len(macdLine)
	macd := macdLine[n-1]
	macdPrev := macdLine[n-2]
	signal := signalLine[n-1]
	hist := histogram[n-1]
	price := candles[len(candles)-1].Close

	info := fmt.Sprintf("macd=%.2f, siganl=%.2f, hist=%.2f, price=%.2f, ema=%.2f", macd, signal, hist, price, ema200)

	if macdPrev < 0 && macd > 0 && hist > 0 && price > ema200 {
		return order.LONG, "Long " + info
	} else if macdPrev > 0 && macd < 0 && hist < 0 && price < ema200 {
		return order.SHORT, "Short " + info
	}
	return order.WAIT, "Wait " + info

}
