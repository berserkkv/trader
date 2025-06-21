package main

import (
	"fmt"
	"github.com/berserkkv/trader/backtest"
	"github.com/berserkkv/trader/bot"
	"github.com/berserkkv/trader/model/enum/symbol"
	"github.com/berserkkv/trader/model/enum/timeframe"
	"github.com/berserkkv/trader/strategy"
	logger "github.com/berserkkv/trader/tools/log"
)

func main() {
	logger.Init("info", "local")

	capital := 100.0
	leverage := 10.0
	takeProfit := 1.0
	stopLoss := 0.3
	st := &strategy.HAStrategy{}
	tf := timeframe.HOUR_1
	smb := symbol.BTCUSDT

	bf := backtest.NewBacktestBotFather(200)
	b := bot.NewBot(tf, st, smb, capital, leverage, takeProfit, stopLoss)

	bf.LoadData(getName(b))

	bf.Randomize()
	startIndex := bf.Index

	bf.Run(b)

	fmt.Println(b.String())
	fmt.Println((bf.Index - startIndex) / 24)

}

func getName(bot *bot.Bot) string {
	data1 := "_2020.08.11_2025.06.21.csv"
	data2 := "_2022.08.11_2025.06.21.csv"

	if bot.Timeframe == timeframe.HOUR_1 {
		return string(bot.Symbol) + "_" + string(bot.Timeframe) + data1
	}
	if bot.Timeframe == timeframe.MINUTE_15 || bot.Timeframe == timeframe.MINUTE_5 {
		return string(bot.Symbol) + "_" + string(bot.Timeframe) + data2
	}
	return ""
}
