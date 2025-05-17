package strategy

import (
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/model/enum/order"
	strategy "github.com/berserkkv/trader/strategy/bbha"
	"log/slog"
)

type Strategy interface {
	Name() string
	Start(candles []model.Candle) order.Command
}

func GetStrategy(name string) Strategy {
	switch name {
	case "BBHA":
		return &strategy.BBHAStrategy{}
	default:
		slog.Error("Strategy not found", "name", name)
		return nil
	}
}
