package strategy

import (
	"fmt"
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/model/enum/order"
	"github.com/berserkkv/trader/model/enum/state"
	"github.com/berserkkv/trader/ta"
)

type SlmaEmaBb struct {
	slmaEma state.State
	bb      state.State
}

func (*SlmaEmaBb) Name() string {
	return "SlmaEmaBb"
}

func (s *SlmaEmaBb) Start(candles []model.Candle) (order.Command, string) {
	closePrices := make([]float64, len(candles))
	for i, c := range candles {
		closePrices[i] = c.Close
	}
	slmaSlice := ta.SLMA(closePrices, 20)
	bbSlice := ta.BollingerPercentB(candles, 20)
	emaSlice := ta.EMA(candles, 200)

	slma := slmaSlice[len(slmaSlice)-1]
	bb := bbSlice[len(bbSlice)-1]
	ema := emaSlice[len(emaSlice)-1]
	price := candles[len(candles)-1].Close

	if s.slmaEma == "" {
		if slma < ema {
			s.slmaEma = state.CrossDown
		} else {
			s.slmaEma = state.CrossUp
		}
	}

	if s.slmaEma == state.CrossUp && slma <= ema {
		s.slmaEma = state.CrossDown
		s.bb = state.Neutral
		if bb < 10 {
			s.bb = state.Under10
		}
		return order.WAIT, "SlmaEma crossed down"
	}

	if s.slmaEma == state.CrossDown && slma >= ema {
		s.slmaEma = state.CrossUp
		s.bb = state.Neutral
		if bb > 90 {
			s.bb = state.Above90
		}

		return order.WAIT, "SlmaEma crossed up"
	}

	if s.slmaEma == state.CrossUp && bb < 10 {
		s.bb = state.Under10
		return order.WAIT, "BB under 10"
	}

	if s.slmaEma == state.CrossDown && bb > 90 {
		s.bb = state.Above90
		return order.WAIT, "BB above 90"
	}

	if s.slmaEma == state.CrossUp && s.bb == state.Under10 && bb > 10 && price > ema {
		s.bb = state.Neutral
		s.slmaEma = state.Neutral
		return order.LONG, "Long"
	}

	if s.slmaEma == state.CrossDown && s.bb == state.Above90 && bb < 90 && price < ema {
		s.bb = state.Neutral
		s.slmaEma = state.Neutral
		return order.SHORT, "Short"
	}

	return order.WAIT, s.String()
}

func (s *SlmaEmaBb) String() string {
	return fmt.Sprintf("SlmaEma: %s, BB: %s", s.slmaEma, s.bb)
}
