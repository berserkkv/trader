package calculator

import (
	"github.com/berserkkv/trader/constant"
	"github.com/berserkkv/trader/model/enum/order"
	"math"
)

func CalculateStopLoss(price float64, percent float64, orderType order.Command) float64 {
	if orderType == order.SHORT {
		return price + (price * percent / 100)
	}
	return price - (price * percent / 100)
}

func CalculateTakeProfit(price float64, percent float64, orderType order.Command) float64 {
	if orderType == order.SHORT {
		return price - (price * percent / 100)
	}
	return price + (price * percent / 100)
}

func CalculateBuyQuantity(price float64, capital float64) float64 {
	if price == 0 {
		return 0
	}
	return capital / price
}

func CalculateTakerFee(price float64) float64 {
	return constant.TakerFeeRate * price
}

func CalculateMakerFee(price float64) float64 {
	return price * constant.MakerFeeRate
}

func CalculatePNL(price float64, capital float64, quantity float64, orderType order.Command) float64 {
	if orderType == order.LONG {
		return (price * quantity) - capital
	} else if orderType == order.SHORT {
		return capital - (price * quantity)
	}
	return 0
}
func CalculateRoe(entryPrice, exitPrice, leverage float64, orderType order.Command) float64 {
	if orderType == order.LONG {
		return ((exitPrice - entryPrice) / entryPrice) * 100 * leverage
	} else if orderType == order.SHORT {
		return ((entryPrice - exitPrice) / entryPrice) * 100 * leverage
	}
	return 0
}

func CalculatePairTradingSpread(prices1, prices2 []float64) float64 {
	var spread []float64
	for i := range prices1 {
		spread = append(spread, math.Log(prices1[i])-math.Log(prices2[i]))
	}

	n := float64(len(spread))
	var sum, mean, stddev float64

	for _, v := range spread {
		sum += v
	}
	mean = sum / n

	for _, v := range spread {
		stddev += math.Pow(v-mean, 2)
	}
	stddev = math.Sqrt(stddev / n)

	latest := spread[len(spread)-1]

	return (latest - mean) / stddev
}
