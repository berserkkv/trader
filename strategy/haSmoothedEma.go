package strategy

import (
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/model/enum/order"
	"github.com/berserkkv/trader/ta"
	"log/slog"
)

type HASmoothedEMAStrategy struct{}

func (s HASmoothedEMAStrategy) Name() string {
	return "HASmoothedEMA"
}

// Enters LONG on Heikin Ashi bullish color change if price is above EMA-40.
// Enters SHORT on bearish color change if price is below EMA-40.
// Otherwise, waits — filters trades using trend confirmation.

func (s HASmoothedEMAStrategy) Start(candles []model.Candle) order.Command {
	ema40 := ta.CalculateEMA(candles, 40)

	ha := ta.CalculateSmoothedHeikinAshi(candles, 3)

	changed, color := ta.DetectHeikinAshiColorChange(ha)

	slog.Debug(s.Name(), "changed", changed, "color", color, "ema", ema40[len(ema40)-1], "price", candles[len(candles)-1].Close)
	if changed {
		if color == "green" && candles[len(candles)-1].Close > ema40[len(ema40)-1] {
			return order.LONG
		} else if color == "red" && candles[len(candles)-1].Close < ema40[len(ema40)-1] {
			return order.SHORT
		}
	}

	return order.WAIT
}
