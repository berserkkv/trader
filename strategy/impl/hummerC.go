package strategyImpl

import (
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/model/enum/order"
	"github.com/berserkkv/trader/ta"
)

type HummerC struct{}

func (*HummerC) Name() string {
	return "HammerC"
}

func (s *HummerC) Run(candles []model.Candle) (order.Command, string) {
	if len(candles) == 0 {
		return order.WAIT, "no candles"
	}
	n := len(candles)

	if ta.IsRedCandle(&candles[n-4]) && ta.IsRedCandle(&candles[n-5]) && ta.IsBullishHammer(&candles[n-3]) && ta.IsGreenCandle(&candles[n-2]) && ta.IsGreenCandle(&candles[n-1]) {
		return order.LONG, "Long"
	}
	if ta.IsGreenCandle(&candles[n-4]) && ta.IsGreenCandle(&candles[n-5]) && ta.IsBearishHammer(&candles[n-3]) && ta.IsRedCandle(&candles[n-2]) && ta.IsRedCandle(&candles[n-1]) {
		return order.SHORT, "Short"
	}
	return order.WAIT, "Wait"
}
