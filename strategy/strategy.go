package strategy

import (
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/model/enum/order"
	"log/slog"
)

type Strategy interface {
	Name() string
	Start(candles []model.Candle) order.Command
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
	default:
		slog.Error("Strategy not found", "name", name)
		return nil
	}
}
