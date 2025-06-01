package strategy

import (
	"fmt"
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/model/enum/order"
	"github.com/berserkkv/trader/ta"
	"log/slog"
)

type HASmoothedStrategy struct{}

func (s HASmoothedStrategy) Name() string {
	return "HASmoothed"
}

func (s HASmoothedStrategy) Start(candles []model.Candle) (order.Command, string) {
	ha := ta.CalculateSmoothedHeikinAshi(candles, 3)

	changed, color := ta.DetectHeikinAshiColorChange(ha)

	info := fmt.Sprintf("HAColor=%s, HAChanged=%t", color, changed)

	slog.Debug(s.Name(), "changed", changed, "color", color)

	if changed {
		if color == "green" {
			return order.LONG, "LONG " + info
		} else if color == "red" {
			return order.SHORT, "SHORT " + info
		}
	}

	return order.WAIT, "WAIT " + info
}
