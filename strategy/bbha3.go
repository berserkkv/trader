package strategy

import (
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/model/enum/order"
	"github.com/berserkkv/trader/ta"
	"log/slog"
)

type BBHA3 struct{}

func (BBHA3) Name() string {
	return "BBHA3"
}

func (s BBHA3) Start(candles []model.Candle) order.Command {
	bb20 := ta.CalculateBollingerPercentB(candles, 20)
	ha := ta.CalculateHeikinAshi(candles)
	changed, color := ta.DetectHeikinAshiColorChange(ha)

	slog.Debug(s.Name(), "bb", bb20[len(bb20)-1], "bb-2", bb20[len(bb20)-2], "HeikinAshi", color, "HA changed", changed)

	if changed && color == "green" && bb20[len(bb20)-1] > 0.5 {
		return order.LONG
	} else if changed && color == "red" && bb20[len(bb20)-1] < 0.5 {
		return order.SHORT
	}

	return order.WAIT
}
