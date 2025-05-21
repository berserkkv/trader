package strategy

import (
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/model/enum/order"
	"github.com/berserkkv/trader/ta"
	"log/slog"
)

type BBHAStrategy struct{}

func (s BBHAStrategy) Name() string {
	return "BBHA"
}

func (s BBHAStrategy) Start(candles []model.Candle) order.Command {

	bb20 := ta.CalculateBollingerPercentB(candles, 20)
	ha := ta.CalculateHeikinAshi(candles)
	changed, lastColor := ta.DetectHeikinAshiColorChange(ha)

	// BB %B crosses above 0.3 (from ≤0.3 to >0.3)
	bbCrossUp := bb20[len(bb20)-2] <= 0.3 && bb20[len(bb20)-1] > 0.3

	// BB %B crosses below 0.7 (from ≥0.7 to <0.7)
	bbCrossDown := bb20[len(bb20)-2] >= 0.7 && bb20[len(bb20)-1] < 0.7

	slog.Info(s.Name(), "bb", bb20[len(bb20)-1], "bb-2", bb20[len(bb20)-2], "HeikinAshi", lastColor, "HA changed", changed)

	if bbCrossUp && changed && lastColor == "green" {
		return order.LONG
	}

	if bbCrossDown && changed && lastColor == "red" {
		return order.SHORT
	}

	return order.WAIT
}
