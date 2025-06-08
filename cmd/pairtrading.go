package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"time"

	"github.com/berserkkv/trader/model/enum/symbol"
	"github.com/berserkkv/trader/model/enum/timeframe"
)

type KlineData [][]interface{}

type ZScorePoint struct {
	Time   time.Time `json:"time"`
	ZScore float64   `json:"z_score"`
}

func fetchKlines(symbol string, interval string, limit int) ([]float64, []time.Time, error) {
	url := fmt.Sprintf("https://api.binance.com/api/v3/klines?symbol=%s&interval=%s&limit=%d", symbol, interval, limit)

	resp, err := http.Get(url)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	var rawData KlineData
	if err := json.Unmarshal(body, &rawData); err != nil {
		return nil, nil, err
	}

	var closes []float64
	var times []time.Time
	for _, k := range rawData {
		price := k[4].(string)
		openTime := int64(k[0].(float64)) // milliseconds

		var f float64
		fmt.Sscanf(price, "%f", &f)

		closes = append(closes, f)
		times = append(times, time.UnixMilli(openTime))
	}
	return closes, times, nil
}

func calculateZScore(slice []float64) float64 {
	n := float64(len(slice))
	if n == 0 {
		return 0
	}
	var sum, mean, stddev float64

	for _, v := range slice {
		sum += v
	}
	mean = sum / n

	for _, v := range slice {
		stddev += math.Pow(v-mean, 2)
	}
	stddev = math.Sqrt(stddev / n)

	latest := slice[len(slice)-1]
	if stddev == 0 {
		return 0
	}
	return (latest - mean) / stddev
}

func rollingZScores(spread []float64, times []time.Time, window int) []ZScorePoint {
	var zscores []ZScorePoint
	for i := window; i <= len(spread); i++ {
		subSpread := spread[i-window : i]
		z := calculateZScore(subSpread)
		zscores = append(zscores, ZScorePoint{
			Time:   times[i-1], // end of window
			ZScore: z,
		})
	}
	return zscores
}

func GetZscores(symbol1, symbol2 symbol.Symbol, interval timeframe.Timeframe, limit int) ([]ZScorePoint, error) {
	const windowSize = 20

	pricesA, times, err := fetchKlines(string(symbol1), string(interval), limit)
	if err != nil {
		return nil, err
	}
	pricesB, _, err := fetchKlines(string(symbol2), string(interval), limit)
	if err != nil {
		return nil, err
	}

	if len(pricesA) != len(pricesB) {
		return nil, fmt.Errorf("mismatched price lengths")
	}

	var spread []float64
	for i := range pricesA {
		spread = append(spread, math.Log(pricesA[i])-math.Log(pricesB[i]))
	}

	zscores := rollingZScores(spread, times, windowSize)

	for _, z := range zscores {
		fmt.Printf("Time: %s, Z-score: %.2f\n", z.Time.Format("2006-01-02 15:04"), z.ZScore)
	}

	return zscores, nil
}
