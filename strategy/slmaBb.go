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

func (s *SlmaBb) Start(candles []model.Candle) (order.Command, string) {

	bb20 := ta.BollingerPercentB(candles, 20)

	closePrices := make([]float64, len(candles))
	for i, c := range candles {
		closePrices[i] = c.Close
	}
	slma := ta.SLMA(closePrices, 20)

	lastBB20 := bb20[len(bb20)-1]
	lastSLMA20 := slma[len(slma)-1]
	lastClosePrice := candles[len(candles)-1].Close

	info := fmt.Sprintf("price=%.2f, bb20=%.2f, slma=%.2f", lastClosePrice, lastBB20, lastSLMA20)

	if lastBB20 > 90 {
		s.bbState = state.Above90
		info += fmt.Sprintf(" bbState=%s", s.bbState)
		return order.WAIT, "Wait " + info
	}

	if lastBB20 < 10 {
		s.bbState = state.Under10
		info += fmt.Sprintf(" bbState=%s", s.bbState)
		return order.WAIT, "Wait " + info
	}

	if s.bbState == state.Above90 && lastSLMA20 > lastClosePrice {
		s.bbState = state.Neutral
		info += fmt.Sprintf(" bbState=%s", s.bbState)
		return order.SHORT, "SHORT " + info
	}

	if s.bbState == state.Under10 && lastSLMA20 < lastClosePrice {
		s.bbState = state.Neutral
		info += fmt.Sprintf(" bbState=%s", s.bbState)
		return order.LONG, "LONG " + info
	}

	return order.WAIT, "Wait " + info

}
