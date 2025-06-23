package strategyImpl

import (
	"fmt"
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/model/enum/order"
	"github.com/berserkkv/trader/model/enum/state"
	"github.com/berserkkv/trader/service/tools"
	"github.com/berserkkv/trader/ta"
)

type SlmaEmaBb2 struct {
	bbPrevState state.State
}

func (s *SlmaEmaBb2) Name() string {
	return "SlmaEmaBb2"
}

func (s *SlmaEmaBb2) Run(candles []model.Candle) (order.Command, string) {
	slma20 := ta.SLMA(tools.GetClosePrices(candles), 20)
	ema200 := ta.EMA(tools.GetClosePrices(candles), 200)
	price := candles[len(candles)-1].Close
	bb20 := ta.BollingerPercentB(candles, 20)
	bbUpperBorder := 0.9
	bbLowerBorder := 0.1
	bbMiddle := 0.5

	info := fmt.Sprintf("EMA: %.2f, SLMA: %.2f, Price: %.2f, BB: %.2f", ema200, slma20, price, bb20)

	if bb20 < bbLowerBorder {
		s.bbPrevState = state.Under10
	} else if bb20 > bbUpperBorder {
		s.bbPrevState = state.Above90
	}

	if slma20 > ema200 && price > slma20 && s.bbPrevState == state.Under10 && bb20 > bbMiddle {
		s.bbPrevState = state.Neutral
		return order.LONG, "Long " + info
	} else if slma20 < ema200 && price < slma20 && s.bbPrevState == state.Above90 && bb20 < bbMiddle {
		s.bbPrevState = state.Neutral
		return order.SHORT, "Short " + info
	}
	return order.WAIT, info + string(s.bbPrevState)
}
