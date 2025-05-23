package calculator

const (
	takerFeeRate float64 = 0.0004 // 0.04%
	makerFeeRate float64 = 0.0002 // 0.02%
)

func CalculateStopLossWithPercent(price float64, percent float64, isShort bool) float64 {
	if isShort {
		return price + (price * percent / 100)
	}
	return price - (price * percent / 100)
}

func CalculateTakeProfitWithPercent(price float64, percent float64, isShort bool) float64 {
	if isShort {
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
	return takerFeeRate * price
}

func CalculateMakerFee(price float64) float64 {
	return price * makerFeeRate
}

func CalculatePNLForLong(price float64, capital float64, quantity float64) float64 {
	return (price * quantity) - capital
}

func CalculatePNLForShort(price float64, capital float64, quantity float64) float64 {
	return capital - (price * quantity)
}

func CalculatePNLPercentForLong(entryPrice float64, exitPrice float64) float64 {
	return ((exitPrice - entryPrice) / entryPrice) * 100
}

func CalculatePNLPercentForShort(entryPrice float64, exitPrice float64) float64 {
	return ((entryPrice - exitPrice) / entryPrice) * 100
}
