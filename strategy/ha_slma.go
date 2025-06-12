package strategy

import (
	"fmt"
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/model/enum/order"
	"github.com/berserkkv/trader/ta"
)

type HA_SLMA struct{}

func (HA_SLMA) Name() string {
	return "HA_SLMA"
}
func (s HA_SLMA) Start(candles []model.Candle) (order.Command, string) {
	closePrice := make([]float64, len(candles))
	for i, c := range candles {
		closePrice[i] = c.Close
	}
	slma20 := ta.SLMA(closePrice, 20)

	ha := ta.CalculateHeikinAshi(candles[:len(candles)-1])

	changed, color := ta.DetectHeikinAshiColorChange(ha)

	info := fmt.Sprintf("price=%.2f, slma=%.2f, HAColor=%s, HAChanged=%t",
		candles[len(candles)-2].Close, slma20[len(slma20)-1], color, changed)
	z := len(candles) - 1

	if changed {
		if color == "green" && candles[z-1].Close > slma20[len(slma20)-1] && ta.CandleColor(candles[z-2]) == 1 && ta.CandleColor(candles[z]) == 1 {
			return order.LONG, "LONG " + info
		} else if color == "red" && candles[z-1].Close < slma20[len(slma20)-1] && ta.CandleColor(candles[z-2]) == -1 && ta.CandleColor(candles[z]) == -1 {
			return order.SHORT, "SHORT " + info
		}
	}

	return order.WAIT, "WAIT " + info
}
