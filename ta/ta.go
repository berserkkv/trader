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

func DetectHeikinAshiColorChange(has []model.HeikinAshi) (changed bool, lastColor string) {
	if len(has) < 2 {
		return false, ""
	}
	prev := has[len(has)-2].Color
	curr := has[len(has)-1].Color
	return prev != curr, curr
}
