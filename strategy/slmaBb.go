package strategy

import (
	"fmt"
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/model/enum/order"
	"github.com/berserkkv/trader/model/enum/state"
	"github.com/berserkkv/trader/ta"
)

type SlmaBb struct {
	bbState state.State
}

func (s *SlmaBb) Name() string {
	return "SlmaBb"
}

func (s *SlmaBb) Run(candles []model.Candle) (order.Command, string) {
	if len(candles) < 20 {
		return order.WAIT, "Not enough candles"
	}

	bb20 := ta.BollingerPercentB(candles, 20)

	closePrices := make([]float64, len(candles))
	for i, c := range candles {
		closePrices[i] = c.Close
	}
	slma := ta.SLMA(closePrices, 20)

	lastClosePrice := candles[len(candles)-1].Close

	info := fmt.Sprintf("price=%.2f, bb20=%.2f, slma=%.2f", lastClosePrice, bb20, slma)

	if bb20 > 0.9 {
		s.bbState = state.Above90
		info += fmt.Sprintf(" bbState=%s", s.bbState)
		return order.WAIT, "Wait " + info
	}

	if bb20 < 0.1 {
		s.bbState = state.Under10
		info += fmt.Sprintf(" bbState=%s", s.bbState)
		return order.WAIT, "Wait " + info
	}

	if s.bbState == state.Above90 && slma > lastClosePrice {
		s.bbState = state.Neutral
		info += fmt.Sprintf(" bbState=%s", s.bbState)
		return order.SHORT, "SHORT " + info
	}

	if s.bbState == state.Under10 && slma < lastClosePrice {
		s.bbState = state.Neutral
		info += fmt.Sprintf(" bbState=%s", s.bbState)
		return order.LONG, "LONG " + info
	}

	return order.WAIT, "Wait " + info

}
