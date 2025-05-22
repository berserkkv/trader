package strategy

import (
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/model/enum/order"
	"github.com/berserkkv/trader/ta"
	"log/slog"
)

type HASmoothedStrategy struct{}

func (s HASmoothedStrategy) Name() string {
	return "HASmoothed"
}

func (s HASmoothedStrategy) Start(candles []model.Candle) order.Command {
	ha := ta.CalculateSmoothedHeikinAshi(candles, 3)

	changed, color := ta.DetectHeikinAshiColorChange(ha)

	slog.Debug(s.Name(), "changed", changed, "color", color)

	if changed {
		if color == "green" {
			return order.LONG
		} else if color == "red" {
			return order.SHORT
		}
	}

	return order.WAIT
}
