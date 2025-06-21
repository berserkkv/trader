package main

import (
	"context"
	"encoding/csv"
	"log"
	"os"
	"time"

	binance "github.com/adshao/go-binance/v2"
)

func main() {
	client := binance.NewClient("", "") // no API key needed for public data

	symbol := "XRPUSDT"
	interval := "1h"
	limit := 1000 // max per request

	// Set start date (Binance listing date for SOL)
	startTime := time.Date(2020, 8, 11, 0, 0, 0, 0, time.UTC).UnixMilli()
	endTime := time.Now().UnixMilli()

	start := time.UnixMilli(startTime).Format("2006.01.02")
	end := time.UnixMilli(endTime).Format("2006.01.02")
	// Create file
	file, err := os.Create(symbol + "_" + interval + "_" + start + "_" + end + ".csv")
	if err != nil {
		log.Fatal("Error creating file:", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"Time", "Open", "High", "Low", "Close", "Volume"})

	for startTime < endTime {
		klines, err := client.NewKlinesService().
			Symbol(symbol).
			Interval(interval).
			Limit(limit).
			StartTime(startTime).
			Do(context.Background())

		if err != nil {
			log.Fatalf("Error fetching data: %v", err)
		}

		if len(klines) == 0 {
			break
		}

		for _, k := range klines {
			writer.Write([]string{
				time.UnixMilli(k.OpenTime).Format(time.RFC3339),
				k.Open,
				k.High,
				k.Low,
				k.Close,
				k.Volume,
			})
		}

		// Move start time forward to next candle
		startTime = klines[len(klines)-1].OpenTime + 1
		log.Printf("Fetched up to: %s", time.UnixMilli(startTime).Format(time.RFC3339))

		// Sleep to avoid hitting API rate limits
		time.Sleep(500 * time.Millisecond)
	}

	log.Println("Download complete! Data saved to solusdt_all.csv")
}
