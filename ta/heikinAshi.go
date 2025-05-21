package ta

import (
	"github.com/berserkkv/trader/model"
	"math"
)

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

func CalculateSmoothedHeikinAshi(candles []model.Candle, smoothingPeriod int) []model.HeikinAshi {
	n := len(candles)
	result := make([]model.HeikinAshi, n)
	rawHA := make([]model.HeikinAshi, n)

	// Step 1: Calculate raw Heikin-Ashi
	for i := 0; i < n; i++ {
		c := candles[i]
		ha := model.HeikinAshi{}
		ha.Close = (c.Open + c.High + c.Low + c.Close) / 4

		if i == 0 {
			ha.Open = (c.Open + c.Close) / 2
		} else {
			ha.Open = (rawHA[i-1].Open + rawHA[i-1].Close) / 2
		}

		ha.High = math.Max(c.High, math.Max(ha.Open, ha.Close))
		ha.Low = math.Min(c.Low, math.Min(ha.Open, ha.Close))

		rawHA[i] = ha
	}

	// Step 2: Smooth Open and Close with SMA
	for i := 0; i < n; i++ {
		start := i - smoothingPeriod + 1
		if start < 0 {
			start = 0
		}
		count := i - start + 1

		var sumOpen, sumClose float64
		for j := start; j <= i; j++ {
			sumOpen += rawHA[j].Open
			sumClose += rawHA[j].Close
		}

		smoothed := rawHA[i]
		smoothed.Open = sumOpen / float64(count)
		smoothed.Close = sumClose / float64(count)

		if smoothed.Close >= smoothed.Open {
			smoothed.Color = "green"
		} else {
			smoothed.Color = "red"
		}

		result[i] = smoothed
	}

	return result
}
