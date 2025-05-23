package strategy

import (
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/model/enum/order"
	"github.com/berserkkv/trader/ta"
	"log/slog"
)

type BBHA2Strategy struct{}

func (s BBHA2Strategy) Name() string {
	return "BBHA2"
}

func (s BBHA2Strategy) Start(candles []model.Candle) order.Command {
	bb20 := ta.CalculateBollingerPercentB(candles, 20)
	ha := ta.CalculateHeikinAshi(candles)
	changed, lastColor := ta.DetectHeikinAshiColorChange(ha)

	slog.Debug(s.Name(), "bb", bb20[len(bb20)-1], "bb-2", bb20[len(bb20)-2], "HeikinAshi", lastColor, "HA changed", changed)

	if changed && lastColor == "green" && bb20[len(bb20)-1] < 0.5 {
		return order.LONG
	}

	if changed && lastColor == "red" && bb20[len(bb20)-1] > 0.5 {
		return order.SHORT
	}

	return order.WAIT
}
