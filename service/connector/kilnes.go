package connector

import (
	"encoding/json"
	"fmt"
	"github.com/berserkkv/trader/model"
	"io"
	"net/http"
	"strconv"
	"time"
)

func FetchKlines(symbol string, interval string, limit int) ([]model.Candle, error) {
	url := fmt.Sprintf("https://api.binance.com/api/v3/klines?symbol=%s&interval=%s&limit=%d", symbol, interval, limit)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var raw [][]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	candles := make([]model.Candle, 0, len(raw))

	for _, k := range raw {

		openTime := time.UnixMilli(int64(k[0].(float64)))

		open, _ := strconv.ParseFloat(k[1].(string), 64)
		high, _ := strconv.ParseFloat(k[2].(string), 64)
		low, _ := strconv.ParseFloat(k[3].(string), 64)
		closePrice, _ := strconv.ParseFloat(k[4].(string), 64)
		volume, _ := strconv.ParseFloat(k[5].(string), 64)

		candles = append(candles, model.Candle{
			Time:   openTime,
			Open:   open,
			High:   high,
			Low:    low,
			Close:  closePrice,
			Volume: volume,
		})
	}

	return candles, nil
}

func main() {
	candles, err := FetchKlines("SOLUSDT", "15m", 200)
	if err != nil {
		panic(err)
	}
	for _, c := range candles {
		fmt.Printf("%s - O:%.2f H:%.2f L:%.2f C:%.2f\n", c.Time.Format(time.RFC822), c.Open, c.High, c.Low, c.Close)
	}
}
