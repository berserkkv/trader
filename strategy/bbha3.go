package strategy

import (
	"fmt"
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/model/enum/order"
	"github.com/berserkkv/trader/ta"
	"log/slog"
)

type BBHA3 struct{}

func (BBHA3) Name() string {
	return "BBHA3"
}

func (s BBHA3) Start(candles []model.Candle) (order.Command, string) {
	bb20 := ta.BollingerPercentB(candles, 20)
	ha := ta.CalculateHeikinAshi(candles)
	changed, color := ta.DetectHeikinAshiColorChange(ha)

	slog.Debug(s.Name(), "bb", bb20[len(bb20)-1], "bb-2", bb20[len(bb20)-2], "HeikinAshi", color, "HA changed", changed)

	info := fmt.Sprintf("bb=%.2f, HAColor=%s, HAChanged=%t", bb20[len(bb20)-1], color, changed)

	if changed && color == "green" && bb20[len(bb20)-1] > 0.5 {
		return order.LONG, "LONG " + info
	} else if changed && color == "red" && bb20[len(bb20)-1] < 0.5 {
		return order.SHORT, "SHORT " + info
	}

	return order.WAIT, "WAIT " + info
}
