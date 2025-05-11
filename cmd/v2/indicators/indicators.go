package indicators

func SMA(candles []Candle, period int) []float64 {
	if len(candles) < period {
		return nil
	}
	smas := make([]float64, 0, len(candles)-period+1)
	for i := 0; i <= len(candles)-period; i++ {
		sum := 0.0
		for j := 0; j < period; j++ {
			sum += candles[i+j].Close
		}
		smas = append(smas, sum/float64(period))
	}
	return smas
}

func RSI(candles []Candle, period int) []float64 {
	if len(candles) <= period {
		return nil
	}

	rsi := make([]float64, 0, len(candles)-period)
	var gains, losses float64

	for i := 1; i <= period; i++ {
		change := candles[i].Close - candles[i-1].Close
		if change >= 0 {
			gains += change
		} else {
			losses -= change
		}
	}

	avgGain := gains / float64(period)
	avgLoss := losses / float64(period)

	if avgLoss == 0 {
		rsi = append(rsi, 100.0)
	} else {
		rs := avgGain / avgLoss
		rsi = append(rsi, 100.0-(100.0/(1+rs)))
	}

	for i := period + 1; i < len(candles); i++ {
		change := candles[i].Close - candles[i-1].Close
		if change >= 0 {
			avgGain = (avgGain*(float64(period-1)) + change) / float64(period)
			avgLoss = (avgLoss * float64(period-1)) / float64(period)
		} else {
			avgGain = (avgGain * float64(period-1)) / float64(period)
			avgLoss = (avgLoss*(float64(period-1)) - change) / float64(period)
		}

		if avgLoss == 0 {
			rsi = append(rsi, 100.0)
		} else {
			rs := avgGain / avgLoss
			rsi = append(rsi, 100.0-(100.0/(1+rs)))
		}
	}

	return rsi
}

func Volume(candles []Candle) []float64 {
	volumes := make([]float64, len(candles))
	for i, c := range candles {
		volumes[i] = c.Volume
	}
	return volumes
}
