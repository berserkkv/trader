package strategy

import (
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/model/enum/order"
	strategyImpl "github.com/berserkkv/trader/strategy/impl"
	"log/slog"
)

type Strategy interface {
	Name() string
	Run(candles []model.Candle) (order.Command, string)
}

func GetStrategy(name string) Strategy {
	switch name {
	case "BBHA":
		return &strategyImpl.BBHAStrategy{}
	case "HASmoothed":
		return &strategyImpl.HASmoothedStrategy{}
	case "HA":
		return &strategyImpl.HAStrategy{}
	case "HAEMA":
		return &strategyImpl.HAEMAStrategy{}
	case "HASmoothedEMA":
		return &strategyImpl.HASmoothedEMAStrategy{}
	case "BBHA2":
		return &strategyImpl.BBHA2Strategy{}
	case "BBHA3":
		return &strategyImpl.BBHA3{}
	case "RANDOM":
		return &strategyImpl.Random{}
	case "S":
		return &strategyImpl.Supertrend{}
	case "S2":
		return &strategyImpl.Supertrend2{} // 10
	case "HaSlma":
		return &strategyImpl.HaSlma{}
	case "3Can":
		return &strategyImpl.TripleCandle{}
	case "SlmaBb":
		return &strategyImpl.SlmaBb{}
	case "3CanSlma":
		return &strategyImpl.TripleCandleSlma{}
	case "SlmaEmaBb":
		return &strategyImpl.SlmaEmaBb{}
	case "SlmaEmaBb2":
		return &strategyImpl.SlmaEmaBb2{}
	case "BbHaSlma":
		return &strategyImpl.BbHaSlma{}
	case "Macd":
		return &strategyImpl.Macd{}
	case "EmaMacd":
		return &strategyImpl.EmaMacd{}
	case "HammerC":
		return &strategyImpl.HummerC{}
	default:
		slog.Error("Strategy not found", "name", name)
		return nil
	}
}
