package strategy

import (
	"fmt"
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/model/enum/order"
	"github.com/berserkkv/trader/ta"
	"log/slog"
)

type BBHAStrategy struct{}

func (s BBHAStrategy) Name() string {
	return "BBHA"
}

func (s BBHAStrategy) Start(candles []model.Candle) (order.Command, string) {

	bb20 := ta.CalculateBollingerPercentB(candles, 20)
	ha := ta.CalculateHeikinAshi(candles)
	changed, lastColor := ta.DetectHeikinAshiColorChange(ha)

	// BB %B crosses above 0.3 (from ≤0.3 to >0.3)
	bbCrossUp := bb20[len(bb20)-2] <= 0.3 && bb20[len(bb20)-1] > 0.3

	// BB %B crosses below 0.7 (from ≥0.7 to <0.7)
	bbCrossDown := bb20[len(bb20)-2] >= 0.7 && bb20[len(bb20)-1] < 0.7

	info := fmt.Sprintf("bbLast=%.4f, bbPrev=%.4f, HAColor=%s, HAChanged=%t",
		bb20[len(bb20)-1], bb20[len(bb20)-2], lastColor, changed)

	slog.Debug(s.Name(), "bb", bb20[len(bb20)-1], "bb-2", bb20[len(bb20)-2], "HeikinAshi", lastColor, "HA changed", changed)

	if bbCrossUp && changed && lastColor == "green" {
		return order.LONG, "LONG " + info
	}

	if bbCrossDown && changed && lastColor == "red" {
		return order.SHORT, "SHORT " + info
	}

	return order.WAIT, "WAIT" + info
}
