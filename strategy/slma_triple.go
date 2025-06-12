package strategy

import (
	"fmt"
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/model/enum/order"
	"github.com/berserkkv/trader/ta"
)

type SLMATriple struct{}

func (SLMATriple) Name() string {
	return "SLMATriple"
}
func (s SLMATriple) Start(candles []model.Candle) (order.Command, string) {
	cur := ta.CandleColor(candles[len(candles)-1])
	prev := ta.CandleColor(candles[len(candles)-2])
	prev2 := ta.CandleColor(candles[len(candles)-3])

	if cur == 1 && prev == 1 && prev2 == 1 {
		return order.LONG, "LONG all 3 green"
	} else if cur == -1 && prev == -1 && prev2 == -1 {
		return order.SHORT, "SHORT all 3 red"
	}

	return order.WAIT, fmt.Sprintf("WAIT %d, %d, %d", prev2, prev, cur)
}
