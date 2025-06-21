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
	bbUpperBorder := 0.9
	bbLowerBorder := 0.1

	info := fmt.Sprintf("SLMA: %.2f, EMA: %.2f, BB: %.2f, Price: %.2f",
		slma, ema, bb, price)

	if slma < ema {
		s.slmaEma = state.CrossDown
	} else {
		s.slmaEma = state.CrossUp
	}

	if bb < bbLowerBorder {
		s.bb = state.Under10
	} else if bb > bbUpperBorder {
		s.bb = state.Above90
	}

	if s.slmaEma == state.CrossUp && s.bb == state.Under10 && bb > bbLowerBorder && price > ema {
		s.bb = state.Neutral
		s.slmaEma = state.Neutral
		return order.LONG, "Long " + info + s.String()
	} else if s.slmaEma == state.CrossDown && s.bb == state.Above90 && bb < bbUpperBorder && price < ema {
		s.bb = state.Neutral
		s.slmaEma = state.Neutral
		return order.SHORT, "Short " + info + s.String()
	}

	return order.WAIT, "Wait " + info + s.String()
}

func (s *SlmaEmaBb) String() string {
	return fmt.Sprintf("SlmaEma: %s, BB: %s", s.slmaEma, s.bb)
}
