package connector

import (
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/model/enum/symbol"
	"github.com/berserkkv/trader/model/enum/timeframe"
)

type Connector interface {
	GetPrice(symbol.Symbol) float64
	GetCandles(symbol.Symbol, timeframe.Timeframe, int) []model.Candle
}
