package main

import (
	"github.com/berserkkv/trader/service/connector"
	bb "github.com/berserkkv/trader/strategy"
	"log"
)

func main() {
	candles, err := connector.FetchKlines("SOLUSDT", "15m", 1000)
	if err != nil {
		log.Fatal(err)
	}
	bb.EvaluateSignals(candles)

}
