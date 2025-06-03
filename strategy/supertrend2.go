package strategy

import (
	"fmt"
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/model/enum/order"
	"github.com/berserkkv/trader/ta"
)

type Supertrend2 struct{}

func (Supertrend2) Name() string {
	return "S2"
}

func (Supertrend2) Start(candles []model.Candle) (order.Command, string) {
	n := len(candles) - 1
	if n <= 0 {
		return order.WAIT, "not enough candles"
	}

	s1, d1 := ta.Supertrend(candles, 10, 2.0)
	s2, d2 := ta.Supertrend(candles, 12, 3.0)

	info := fmt.Sprintf(
		"s1=%.2f(d1=%d), s2=%.2f(d2=%d)",
		s1[n], d1[n], s2[n], d2[n],
	)

	// Check for first transition to all green (long)
	if d1[n] == 1 && d2[n] == 1 &&
		(d1[n-1] != 1 || d2[n-1] != 1) {
		return order.LONG, "LONG " + info
	}

	// Check for first transition to all red (short)
	if d1[n] == -1 && d2[n] == -1 &&
		(d1[n-1] != -1 || d2[n-1] != -1) {
		return order.SHORT, "SHORT " + info
	}

	return order.WAIT, "WAIT " + info

}
