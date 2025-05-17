package ta

import (
	"github.com/berserkkv/trader/model"
	"math"
)

func CalculateBollingerPercentB(candles []model.Candle, period int) []float64 {
	n := len(candles)
	result := make([]float64, n)

	for i := period - 1; i < n; i++ {
		var sum float64
		for j := i - period + 1; j <= i; j++ {
			sum += candles[j].Close
		}
		sma := sum / float64(period)

		var variance float64
		for j := i - period + 1; j <= i; j++ {
			diff := candles[j].Close - sma
			variance += diff * diff
		}
		stddev := math.Sqrt(variance / float64(period))

		upper := sma + 2*stddev
		lower := sma - 2*stddev

		percentB := (candles[i].Close - lower) / (upper - lower)
		result[i] = percentB
	}
	return result
}

func CalculateHeikinAshi(candles []model.Candle) []model.HeikinAshi {
	n := len(candles)
	result := make([]model.HeikinAshi, n)

	for i := 0; i < n; i++ {
		c := candles[i]
		ha := model.HeikinAshi{}
		ha.Close = (c.Open + c.High + c.Low + c.Close) / 4

		if i == 0 {
			ha.Open = (c.Open + c.Close) / 2
		} else {
			ha.Open = (result[i-1].Open + result[i-1].Close) / 2
		}

		ha.High = math.Max(c.High, math.Max(ha.Open, ha.Close))
		ha.Low = math.Min(c.Low, math.Min(ha.Open, ha.Close))

		if ha.Close >= ha.Open {
			ha.Color = "green"
		} else {
			ha.Color = "red"
		}
		result[i] = ha
	}
	return result
}

func DetectHeikinAshiColorChange(has []model.HeikinAshi) (changed bool, lastColor string) {
	if len(has) < 2 {
		return false, ""
	}
	prev := has[len(has)-2].Color
	curr := has[len(has)-1].Color
	return prev != curr, curr
}
