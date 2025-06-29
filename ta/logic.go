package ta

import (
	"github.com/berserkkv/trader/model"
	"math"
)

func IsBullishHammer(c *model.Candle) bool {
	if c == nil {
		return false
	}

	body := math.Abs(c.Close - c.Open)
	totalLength := c.High - c.Low

	if totalLength == 0 {
		return false
	}

	bodyTop := max(c.Close, c.Open)
	bodyBottom := min(c.High, c.Low)

	upperWick := c.High - bodyTop
	lowerWick := bodyBottom - c.Low

	return lowerWick >= 2*body && upperWick <= body*0.3
}

func IsBearishHammer(c *model.Candle) bool {
	if c == nil {
		return false
	}

	body := math.Abs(c.Close - c.Open)
	totalLength := c.High - c.Low

	if totalLength == 0 {
		return false
	}

	bodyTop := max(c.Close, c.Open)
	bodyBottom := min(c.High, c.Low)

	upperWick := c.High - bodyTop
	lowerWick := bodyBottom - c.Low

	return upperWick >= body*2 && lowerWick <= 0.3*body
}

func IsGreenCandle(c *model.Candle) bool {
	if c == nil {
		return false
	}
	return c.Close > c.Open
}

func IsRedCandle(c *model.Candle) bool {
	if c == nil {
		return false
	}
	return c.Close < c.Open
}
