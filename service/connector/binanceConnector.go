package connector

import (
	"encoding/json"
	"fmt"
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/model/enum/symbol"
	"github.com/berserkkv/trader/model/enum/timeframe"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type priceResponse struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

type BinanceConnector struct{}

func (BinanceConnector) GetPrice(smb symbol.Symbol) float64 {
	url := "https://fapi.binance.com/fapi/v2/ticker/price?symbol=" + string(smb)

	resp, err := http.Get(url)

	if err != nil {
		slog.Error("Failed to get price: %v", err)
	}
	defer resp.Body.Close()

	var price priceResponse
	if err := json.NewDecoder(resp.Body).Decode(&price); err != nil {
		slog.Error("Failed to decode response: %v", err)
	}

	pr, err := strconv.ParseFloat(price.Price, 64)
	if err != nil {
		slog.Error("Failed to parse price: %v", err)
	}

	slog.Debug("Returning price", "price", pr)
	return pr

}

func (BinanceConnector) GetCandles(smb symbol.Symbol, interval timeframe.Timeframe, limit int) []model.Candle {

	url := fmt.Sprintf("https://fapi.binance.com/fapi/v1/klines?symbol=%s&interval=%s&limit=%d", string(smb), interval, limit)
	slog.Debug("GetCandles", "url", url)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		slog.Error("failed to fetch klines", "error", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		var errResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errResp)
		slog.Error("API error", "error", errResp)
		return nil
	}

	var rawKlines [][]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&rawKlines); err != nil {
		slog.Error("failed to decode response", "error", err)
		return nil
	}

	candles := make([]model.Candle, len(rawKlines))
	for i, k := range rawKlines {
		open, _ := strconv.ParseFloat(k[1].(string), 64)
		high, _ := strconv.ParseFloat(k[2].(string), 64)
		low, _ := strconv.ParseFloat(k[3].(string), 64)
		closePrice, _ := strconv.ParseFloat(k[4].(string), 64)
		volume, _ := strconv.ParseFloat(k[5].(string), 64)

		candles[i] = model.Candle{
			Time:   time.UnixMilli(int64(k[0].(float64))),
			Open:   open,
			High:   high,
			Low:    low,
			Close:  closePrice,
			Volume: volume,
		}
	}

	if len(candles) == 0 {
		return candles
	}

	return candles[:len(candles)-1]
}
