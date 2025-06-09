package strategy

import (
	"fmt"
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/model/enum/order"
	"github.com/berserkkv/trader/ta"
)

type HA_SLMA struct {
	crossed int
}

func (HA_SLMA) Name() string {
	return "HA_SLMA"
}
func (s HA_SLMA) Start(candles []model.Candle) (order.Command, string) {
	closePrice := make([]float64, len(candles))
	for i, c := range candles {
		closePrice[i] = c.Close
	}
	slma20 := ta.SLMA(closePrice, 20)

	ha := ta.CalculateHeikinAshi(candles)

	changed, color := ta.DetectHeikinAshiColorChange(ha)

	info := fmt.Sprintf("price=%.2f, slma=%.2f, HAColor=%s, HAChanged=%t, crossed=%d",
		candles[len(candles)-1].Close, slma20[len(slma20)-1], color, changed, s.crossed)
	z := len(candles) - 1

	if s.crossed != 0 {
		cr := s.crossed
		s.crossed = 0
		if cr == 1 && ta.CandleColor(candles[z]) == 1 {
			return order.LONG, "LONG " + info
		} else if cr == -1 && ta.CandleColor(candles[z]) == -1 {
			return order.SHORT, "SHORT " + info
		}
	} else if changed {
		if color == "green" && candles[z].Close > slma20[len(slma20)-1] && ta.CandleColor(candles[z-1]) == 1 {
			s.crossed = 1
		} else if color == "red" && candles[z].Close < slma20[len(slma20)-1] && ta.CandleColor(candles[z-1]) == -1 {
			s.crossed = -1
		}
	}
	return order.WAIT, "WAIT " + info
}
