package ta

import (
	"github.com/berserkkv/trader/model"
	"math"
)

func BollingerPercentB(candles []model.Candle, period int) []float64 {
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

func EMA(candles []model.Candle, period int) []float64 {
	n := len(candles)
	result := make([]float64, n)

	if n == 0 || period <= 0 {
		return result
	}

	k := 2.0 / float64(period+1) // smoothing factor

	// Start by calculating the SMA for the first 'period' candles as initial EMA
	var sum float64
	for i := 0; i < period && i < n; i++ {
		sum += candles[i].Close
	}
	initialEMA := sum / float64(period)
	result[period-1] = initialEMA

	// Calculate EMA for remaining candles
	for i := period; i < n; i++ {
		prevEMA := result[i-1]
		close := candles[i].Close
		ema := (close-prevEMA)*k + prevEMA
		result[i] = ema
	}

	// Optional: set NaN for values before period
	for i := 0; i < period-1; i++ {
		result[i] = math.NaN()
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

func ATR(candles []model.Candle, period int) []float64 {
	tr := make([]float64, len(candles))
	for i := 1; i < len(candles); i++ {
		h_l := candles[i].High - candles[i].Low
		h_pc := math.Abs(candles[i].High - candles[i-1].Close)
		l_pc := math.Abs(candles[i].Low - candles[i-1].Close)
		tr[i] = math.Max(h_l, math.Max(h_pc, l_pc))
	}

	atr := make([]float64, len(candles))
	var sum float64
	for i := 1; i <= period && i < len(candles); i++ {
		sum += tr[i]
	}
	atr[period] = sum / float64(period)

	for i := period + 1; i < len(candles); i++ {
		atr[i] = (atr[i-1]*(float64(period-1)) + tr[i]) / float64(period)
	}
	return atr
}

// Supertrend calculation for slice of Candles
func Supertrend(candles []model.Candle, atrPeriod int, factor float64) ([]float64, []int) {
	atr := ATR(candles, atrPeriod)

	supertrend := make([]float64, len(candles))
	direction := make([]int, len(candles)) // 1 = uptrend, -1 = downtrend

	var finalUpperBand, finalLowerBand float64

	for i := atrPeriod; i < len(candles); i++ {
		hl2 := (candles[i].High + candles[i].Low) / 2
		upperBand := hl2 + factor*atr[i]
		lowerBand := hl2 - factor*atr[i]

		if i == atrPeriod {
			finalUpperBand = upperBand
			finalLowerBand = lowerBand
			direction[i] = 1 // start with uptrend
			supertrend[i] = finalLowerBand
			continue
		}

		// Adjust bands according to previous values and price
		if upperBand < finalUpperBand || candles[i-1].Close > finalUpperBand {
			finalUpperBand = upperBand
		}
		if lowerBand > finalLowerBand || candles[i-1].Close < finalLowerBand {
			finalLowerBand = lowerBand
		}

		// Determine trend direction
		if candles[i].Close > finalUpperBand {
			direction[i] = 1
		} else if candles[i].Close < finalLowerBand {
			direction[i] = -1
		} else {
			direction[i] = direction[i-1]
		}

		// Set Supertrend line
		if direction[i] == 1 {
			supertrend[i] = finalLowerBand
		} else {
			supertrend[i] = finalUpperBand
		}
	}

	return supertrend, direction
}

// CalculateSLMA calculates the Smoothed Linear Moving Average over a slice of float64
func SLMA(prices []float64, period int) []float64 {
	if len(prices) < period {
		return nil
	}

	slma := make([]float64, 0, len(prices)-period+1)
	weightSum := period * (period + 1) / 2

	for i := 0; i <= len(prices)-period; i++ {
		var weightedSum float64
		for j := 0; j < period; j++ {
			weight := j + 1
			weightedSum += prices[i+j] * float64(weight)
		}
		slma = append(slma, weightedSum/float64(weightSum))
	}

	return slma
}

func CandleColor(c model.Candle) int {
	if c.Close > c.Open {
		return 1
	} else if c.Close < c.Open {
		return -1
	}
	return 0
}
