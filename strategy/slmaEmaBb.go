package strategy

import (
	"fmt"
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/model/enum/order"
	"github.com/berserkkv/trader/model/enum/state"
	"github.com/berserkkv/trader/service/tools"
	"github.com/berserkkv/trader/ta"
)

type SlmaEmaBb struct {
	bb state.State
}

func (*SlmaEmaBb) Name() string {
	return "SlmaEmaBb"
}

func (s *SlmaEmaBb) Run(candles []model.Candle) (order.Command, string) {
	slma := ta.SLMA(tools.GetClosePrices(candles), 20)
	bb := ta.BollingerPercentB(candles, 20)
	ema := ta.EMA(tools.GetClosePrices(candles), 200)

	price := candles[len(candles)-1].Close
	bbUpperBorder := 0.9
	bbLowerBorder := 0.1

	info := fmt.Sprintf("SLMA: %.2f, EMA: %.2f, BB: %.2f, Price: %.2f",
		slma, ema, bb, price)

	if bb < bbLowerBorder {
		s.bb = state.Under10
	} else if bb > bbUpperBorder {
		s.bb = state.Above90
	}

	if slma > ema && s.bb == state.Under10 && bb > bbLowerBorder && price > ema {
		s.bb = state.Neutral
		return order.LONG, "Long " + info + s.String()
	} else if slma < ema && s.bb == state.Above90 && bb < bbUpperBorder && price < ema {
		s.bb = state.Neutral
		return order.SHORT, "Short " + info + s.String()
	}

	return order.WAIT, "Wait " + info + s.String()
}

func (s *SlmaEmaBb) String() string {
	return fmt.Sprintf(" BB: %s ", s.bb)
}
