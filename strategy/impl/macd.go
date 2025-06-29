package strategyImpl

import (
	"fmt"
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/model/enum/order"
	"github.com/berserkkv/trader/service/tools"
	"github.com/berserkkv/trader/ta"
	"log/slog"
)

type Macd struct{}

func (m *Macd) Name() string {
	return "Macd"
}

func (m *Macd) Run(candles []model.Candle) (order.Command, string) {
	if len(candles) == 0 {
		return order.WAIT, "candles is empty"
	}
	macdLine, signalLine, histogram := ta.MACDSlice(tools.GetClosePrices(candles))

	n := len(macdLine)
	macd := macdLine[n-1]
	signal := signalLine[n-1]
	hist := histogram[n-1]

	info := fmt.Sprintf("macd=%.2f, siganl=%.2f, hist=%.2f", macd, signal, hist)

	if macd > 0 && signal < 0 && hist > 0 {
		slog.Info("Long ", "macd", fmt.Sprintf("%.2f", macd), "signal", fmt.Sprintf("%.3f", signal), "histogram", hist)
		return order.LONG, "Long " + info
	} else if macd < 0 && signal > 0 && hist < 0 {
		slog.Info("Short ", "macd", macd, "signal", signal, "histogram", hist)

		return order.SHORT, "Short " + info
	}
	return order.WAIT, "Wait " + info
}
