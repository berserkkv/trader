package strategy

import (
	"fmt"
	"github.com/berserkkv/trader/model"
	"math"
)

// HeikinAshi returns HA candle components for a slice of candles
func HeikinAshi(candles []model.Candle) []model.Candle {
	ha := make([]model.Candle, len(candles))

	for i := range candles {
		c := candles[i]

		var haOpen, haClose float64

		if i == 0 {
			haOpen = (c.Open + c.Close) / 2
		} else {
			haOpen = ha[i-1].Close
		}

		haClose = c.Close

		haHigh := math.Max(c.High, math.Max(haOpen, haClose))
		haLow := math.Min(c.Low, math.Min(haOpen, haClose))

		ha[i] = model.Candle{
			Time:  c.Time,
			Open:  haOpen,
			Close: haClose,
			High:  haHigh,
			Low:   haLow,
		}
	}

	return ha
}

func linearRegression(prices []float64) (float64, float64) {
	n := float64(len(prices))
	var sumX, sumY, sumXY, sumXX float64

	// Calculate the necessary sums
	for i, price := range prices {
		x := float64(i + 1) // Use 1-based index for x (starting from 1)
		y := price

		sumX += x
		sumY += y
		sumXY += x * y
		sumXX += x * x
	}

	// Calculate the slope (m)
	m := (n*sumXY - sumX*sumY) / (n*sumXX - sumX*sumX)

	// Calculate the intercept (b)
	b := (sumY - m*sumX) / n

	return m, b
}

// Function to calculate Linear Regression with offset
func calculateLinregWithOffset(prices []float64, length int, offset int) []float64 {
	linreg := make([]float64, length)

	// Calculate the linreg values for the source series with the given length and offset
	for i := 0; i <= len(prices)-length; i++ {
		// Get the data points for the current window (length)
		window := prices[i : i+length]

		// Calculate the linear regression for the current window
		m, b := linearRegression(window)

		// Calculate the linreg value at the current point with the offset
		linregValue := b + m*float64(length-1-offset)

		// Store the linreg value
		linreg = append(linreg, linregValue)
	}

	return linreg
}

// BBPercent calculates Bollinger %B
func BBPercent(src []float64, length int, mult float64) []float64 {
	out := make([]float64, len(src))
	for i := length - 1; i < len(src); i++ {
		var sum, stddev float64
		for j := i - length + 1; j <= i; j++ {
			sum += src[j]
		}
		mean := sum / float64(length)
		for j := i - length + 1; j <= i; j++ {
			stddev += math.Pow(src[j]-mean, 2)
		}
		stddev = math.Sqrt(stddev / float64(length))
		upper := mean + mult*stddev
		lower := mean - mult*stddev
		out[i] = (src[i] - lower) / (upper - lower)
	}
	return out
}

func EvaluateSignals(candles []model.Candle) (buy, sell []bool) {
	ha := HeikinAshi(candles)
	lengthLSMA := 20
	offsetLSMA := 7
	lengthBB := 200
	mult := 1.0

	lowSrc := make([]float64, len(candles))
	closeSrc := make([]float64, len(candles))
	for i := range candles {
		lowSrc[i] = candles[i].Low
		closeSrc[i] = candles[i].Close
	}

	lsma := calculateLinregWithOffset(lowSrc, lengthLSMA, offsetLSMA)
	bbr := BBPercent(closeSrc, lengthBB, mult)

	buy = make([]bool, len(candles))
	sell = make([]bool, len(candles))

	for i := 1; i < len(candles); i++ {
		if i < lengthLSMA+offsetLSMA || i < lengthBB {
			continue
		}

		isGreen := ha[i].Close > ha[i].Open
		isRed := ha[i].Close < ha[i].Open

		// Matching Pine Script logic for buy/sell conditions
		buy[i] = ha[i].Close > lsma[i] &&
			isGreen &&
			bbr[i] <= 0 &&
			ha[i-1].Close < lsma[i-1] &&
			ha[i-1].Open < lsma[i-1]

		sell[i] = ha[i].Close < lsma[i] &&
			isRed &&
			bbr[i] >= 1 &&
			ha[i-1].Close > lsma[i-1] &&
			ha[i-1].Open > lsma[i-1]
		//fmt.Printf("%s = %.2f - %.2f\n", candles[i].Time, candles[i].Close, lsma[i])
		if buy[i] {
			fmt.Printf("buy %.2f, time %s\n", candles[i].Close, candles[i].Time)
		}
		if sell[i] {
			fmt.Printf("sell %.2f, time %s\n", candles[i].Close, candles[i].Time)
		}
	}

	return
}
