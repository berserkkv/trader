package tools

import (
	"encoding/csv"
	"github.com/berserkkv/trader/model"
	"log"
	"log/slog"
	"os"
	"strconv"
	"time"
)

func WaitUntilNextAlignedTick(interval time.Duration) {
	now := time.Now()

	// Find how far we are into the interval
	elapsed := time.Duration(now.Minute())*time.Minute +
		time.Duration(now.Second())*time.Second +
		time.Duration(now.Nanosecond())

	// Compute time until next interval
	wait := interval - (elapsed % interval)
	if wait <= 0 {
		wait = interval
	}

	slog.Debug("Waiting until next aligned tick.\n", "wait", wait)

	time.Sleep(wait)
}

func GetClosePrices(candles []model.Candle) []float64 {
	closePrices := make([]float64, len(candles))
	for i, candle := range candles {
		closePrices[i] = candle.Close
	}
	return closePrices
}

func GetOpenPrices(candles []model.Candle) []float64 {
	openPrices := make([]float64, len(candles))
	for i, candle := range candles {
		openPrices[i] = candle.Open
	}
	return openPrices
}

func LoadCandlesFromCSV(path string) ([]model.Candle, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var candles []model.Candle

	for i, row := range records {
		if i == 0 {
			continue // skip header
		}

		t, err := time.Parse(time.RFC3339, row[0])
		if err != nil {
			log.Printf("skipping row %d: invalid time %v", i, err)
			continue
		}

		open, _ := strconv.ParseFloat(row[1], 64)
		high, _ := strconv.ParseFloat(row[2], 64)
		low, _ := strconv.ParseFloat(row[3], 64)
		closePrice, _ := strconv.ParseFloat(row[4], 64)
		vol, _ := strconv.ParseFloat(row[5], 64)

		candle := model.Candle{
			Time:   t,
			Open:   open,
			High:   high,
			Low:    low,
			Close:  closePrice,
			Volume: vol,
		}

		candles = append(candles, candle)
	}

	return candles, nil
}
