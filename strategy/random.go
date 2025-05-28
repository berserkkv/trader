package strategy

import (
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/model/enum/order"
	"math/rand"
	"time"
)

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

type Random struct{}

func (Random) Name() string {
	return "RANDOM"
}

func (Random) Start(candles []model.Candle) order.Command {

	r := rnd.Intn(3)

	switch r {
	case 0:
		return order.WAIT
	case 1:
		return order.LONG
	case 2:
		return order.SHORT
	}
	return order.WAIT
}
