package strategyImpl

import (
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/model/enum/order"
	"github.com/berserkkv/trader/service/tools"
	"github.com/berserkkv/trader/ta"
)

type BbHaSlma struct{}

func (s *BbHaSlma) Name() string {
	return "BbHaSlma"
}

func (s *BbHaSlma) Run(candles []model.Candle) (order.Command, string) {
	bb200 := ta.BollingerPercentBSlice(candles, 200)
	slma20 := ta.SLMA(tools.GetClosePrices(candles), 20)
	haSlice := ta.HeikinAshi(candles)
	ha := haSlice[len(haSlice)-1]
	bb200Cur := bb200[len(bb200)-1]
	bb200Prev := bb200[len(bb200)-2]

	if bb200Prev < 0 && bb200Cur > 0 && ha.Close > slma20 && ha.Close > ha.Open && ha.Open < slma20 {
		return order.LONG, ""
	} else if bb200Prev > 1 && bb200Cur < 1 && ha.Close < slma20 && ha.Close < ha.Open && ha.Open > slma20 {
		return order.SHORT, ""
	}
	return order.WAIT, ""
}
