package strategyImpl

import (
	"fmt"
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/model/enum/order"
	"github.com/berserkkv/trader/ta"
)

type TripleCandleSlma struct{}

func (*TripleCandleSlma) Name() string {
	return "3CanSlma"
}

func (s *TripleCandleSlma) Run(candles []model.Candle) (order.Command, string) {
	if len(candles) < 3 {
		return order.WAIT, "Not enough candles"
	}
	cur := ta.CandleColor(candles[len(candles)-1])
	prev := ta.CandleColor(candles[len(candles)-2])
	prev2 := ta.CandleColor(candles[len(candles)-3])

	closePrices := make([]float64, len(candles))
	for i, c := range candles {
		closePrices[i] = c.Close
	}
	slma := ta.SLMA(closePrices, 20)

	lastPrice := candles[len(candles)-2].Close

	info := fmt.Sprintf("price=%.2f, slma=%.2f", lastPrice, slma)

	if cur == 1 && prev == 1 && prev2 == 1 && slma < lastPrice {
		return order.LONG, "LONG all 3 green " + info
	} else if cur == -1 && prev == -1 && prev2 == -1 && slma > lastPrice {
		return order.SHORT, "SHORT all 3 red " + info
	}

	return order.WAIT, fmt.Sprintf("WAIT %d, %d, %d", prev2, prev, cur) + info
}
