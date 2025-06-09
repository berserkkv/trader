package strategy

import (
	"fmt"
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/model/enum/order"
	"github.com/berserkkv/trader/ta"
	"log/slog"
)

type BBHA2Strategy struct{}

func (s BBHA2Strategy) Name() string {
	return "BBHA2"
}

func (s BBHA2Strategy) Start(candles []model.Candle) (order.Command, string) {
	bb20 := ta.BollingerPercentB(candles, 20)
	ha := ta.CalculateHeikinAshi(candles)
	changed, lastColor := ta.DetectHeikinAshiColorChange(ha)

	info := fmt.Sprintf("bb=%.4f, HAColor=%s, HAChanged=%t",
		bb20[len(bb20)-1], lastColor, changed)

	slog.Debug(s.Name(), "bb", bb20[len(bb20)-1], "bb-2", bb20[len(bb20)-2], "HeikinAshi", lastColor, "HA changed", changed)

	if changed && lastColor == "green" && bb20[len(bb20)-1] < 0.5 {
		return order.LONG, "LONG " + info
	}

	if changed && lastColor == "red" && bb20[len(bb20)-1] > 0.5 {
		return order.SHORT, "SHORT " + info
	}

	return order.WAIT, "WAIT " + info
}
