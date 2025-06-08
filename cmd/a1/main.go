package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
)

type KlineData [][]interface{}

func fetchKlines(symbol string, interval string, limit int) ([]float64, error) {
	url := fmt.Sprintf("https://api.binance.com/api/v3/klines?symbol=%s&interval=%s&limit=%d", symbol, interval, limit)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var rawData KlineData
	if err := json.Unmarshal(body, &rawData); err != nil {
		return nil, err
	}

	var closes []float64
	for _, k := range rawData {
		price := k[4].(string) // close price
		var f float64
		fmt.Sscanf(price, "%f", &f)
		closes = append(closes, f)
	}
	return closes, nil
}

func calculateZScore(spread []float64) float64 {
	n := float64(len(spread))
	var sum, mean, stddev float64

	for _, v := range spread {
		sum += v
	}
	mean = sum / n

	for _, v := range spread {
		stddev += math.Pow(v-mean, 2)
	}
	stddev = math.Sqrt(stddev / n)

	latest := spread[len(spread)-1]
	return (latest - mean) / stddev
}

func main() {
	pricesA, _ := fetchKlines("ETHUSDT", "1h", 100)
	pricesB, _ := fetchKlines("BTCUSDT", "1h", 100)

	var spread []float64
	for i := range pricesA {
		spread = append(spread, math.Log(pricesA[i])-math.Log(pricesB[i]))
	}

	z := calculateZScore(spread)
	fmt.Printf("Current Spread Z-Score: %.2f\n", z)
}
