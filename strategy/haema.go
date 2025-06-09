package strategy

import (
	"fmt"
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/model/enum/order"
	"github.com/berserkkv/trader/ta"
	"log/slog"
)

type HAEMAStrategy struct{}

func (s HAEMAStrategy) Name() string {
	return "HAEMA"
}

// Enters LONG on Heikin Ashi bullish color change if price is above EMA-40.
// Enters SHORT on bearish color change if price is below EMA-40.
// Otherwise, waits — filters trades using trend confirmation.

func (s HAEMAStrategy) Start(candles []model.Candle) (order.Command, string) {
	ema40 := ta.EMA(candles, 40)

	ha := ta.CalculateHeikinAshi(candles)

	changed, color := ta.DetectHeikinAshiColorChange(ha)

	info := fmt.Sprintf("price=%.4f, ema=%.4f, HAColor=%s, HAChanged=%t",
		candles[len(candles)-1].Close, ema40[len(ema40)-1], color, changed)

	slog.Debug(s.Name(), "changed", changed, "color", color, "ema", ema40[len(ema40)-1], "price", candles[len(candles)-1].Close)
	if changed {
		if color == "green" && candles[len(candles)-1].Close > ema40[len(ema40)-1] {
			return order.LONG, "LONG " + info
		} else if color == "red" && candles[len(candles)-1].Close < ema40[len(ema40)-1] {
			return order.SHORT, "SHORT " + info
		}
	}

	return order.WAIT, "WAIT " + info
}
