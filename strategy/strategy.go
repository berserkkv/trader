package strategy

import (
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/model/enum/order"
	"log/slog"
)

type Strategy interface {
	Name() string
	Start(candles []model.Candle) (order.Command, string)
}

func GetStrategy(name string) Strategy {
	switch name {
	case "BBHA":
		return &BBHAStrategy{}
	case "HASmoothed":
		return &HASmoothedStrategy{}
	case "HA":
		return &HAStrategy{}
	case "HAEMA":
		return &HAEMAStrategy{}
	case "HASmoothedEMA":
		return &HASmoothedEMAStrategy{}
	case "BBHA2":
		return &BBHA2Strategy{}
	case "BBHA3":
		return &BBHA3{}
	case "RANDOM":
		return &Random{}
	case "S":
		return &Supertrend{}
	case "S2":
		return &Supertrend2{}
	default:
		slog.Error("Strategy not found", "name", name)
		return nil
	}
}
