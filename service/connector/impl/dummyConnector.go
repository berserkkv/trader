package connectorImpl

import (
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/model/enum/symbol"
	"github.com/berserkkv/trader/model/enum/timeframe"
)

type DummyConnector struct {
}

func (c *DummyConnector) GetPrice(symbol.Symbol) float64 {

	return 0.0
}

func (c *DummyConnector) GetCandles(symbol.Symbol, timeframe.Timeframe, int) []model.Candle {
	return []model.Candle{}
}
