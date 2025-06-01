package strategy

import (
	"fmt"
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/model/enum/order"
	"github.com/berserkkv/trader/ta"
	"log/slog"
)

type HAStrategy struct{}

func (s HAStrategy) Name() string {
	return "HA"
}

func (s HAStrategy) Start(candles []model.Candle) (order.Command, string) {
	ha := ta.CalculateHeikinAshi(candles)

	changed, color := ta.DetectHeikinAshiColorChange(ha)

	slog.Debug(s.Name(), "changed", changed, "color", color)

	info := fmt.Sprintf("HAColor=%s, HAChanged=%t", color, changed)

	if changed {
		if color == "green" {
			return order.LONG, "LONG " + info
		} else if color == "red" {
			return order.SHORT, "SHORT " + info
		}
	}
	return order.WAIT, "WAIT " + info
}
