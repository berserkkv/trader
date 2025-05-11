package trtools

import (
	"github.com/berserkkv/trader/model"
	"math"
)

func GetClose(candles *[]model.Candle) []float64 {
	closes := make([]float64, len(*candles))
	for i, c := range *candles {
		closes[i] = c.Close
	}
	return closes
}

func GetOpen(candles *[]model.Candle) []float64 {
	opens := make([]float64, len(*candles))
	for i, c := range *candles {
		opens[i] = c.Open
	}
	return opens
}

func GetHigh(candles *[]model.Candle) []float64 {
	highs := make([]float64, len(*candles))
	for i, c := range *candles {
		highs[i] = c.High
	}
	return highs
}

func GetLow(candles *[]model.Candle) []float64 {
	lows := make([]float64, len(*candles))
	for i, c := range *candles {
		lows[i] = c.Low
	}
	return lows
}

// Converts normal candles to Heikin-Ashi candles
func ComputeHeikinAshi(candles []model.Candle) []model.HeikinAshi {
	n := len(candles)
	if n == 0 {
		return nil
	}

	ha := make([]model.HeikinAshi, n)
	for i, c := range candles {
		haClose := (c.Open + c.High + c.Low + c.Close) / 4

		var haOpen float64
		if i == 0 {
			haOpen = (c.Open + c.Close) / 2 // first candle uses normal values
		} else {
			haOpen = (ha[i-1].Open + ha[i-1].Close) / 2
		}

		haHigh := math.Max(c.High, math.Max(haOpen, haClose))
		haLow := math.Min(c.Low, math.Min(haOpen, haClose))

		ha[i] = model.HeikinAshi{
			Open:  haOpen,
			Close: haClose,
			High:  haHigh,
			Low:   haLow,
		}
	}
	return ha
}
