package main

import (
	"fmt"
	"github.com/berserkkv/trader/backtest"
	"github.com/berserkkv/trader/bot"
	"github.com/berserkkv/trader/model/enum/symbol"
	"github.com/berserkkv/trader/model/enum/timeframe"
	"github.com/berserkkv/trader/strategy"
	strategyImpl "github.com/berserkkv/trader/strategy/impl"
	logger "github.com/berserkkv/trader/tools/log"
)

func main() {
	logger.Init("info", "local")
	strategies := []strategy.Strategy{
		strategyImpl.HAStrategy{},
		strategyImpl.HASmoothedStrategy{},
		strategyImpl.HAEMAStrategy{},
		strategyImpl.BBHAStrategy{},
		strategyImpl.BBHA2Strategy{},
		strategyImpl.BBHA3{},
		strategyImpl.HASmoothedEMAStrategy{},
		strategyImpl.Random{},
		strategyImpl.Supertrend{},
		strategyImpl.Supertrend2{},
		&strategyImpl.HaSlma{},
		&strategyImpl.TripleCandle{},
		&strategyImpl.SlmaBb{},
		&strategyImpl.TripleCandleSlma{},
		&strategyImpl.SlmaEmaBb{},
		&strategyImpl.SlmaEmaBb2{},
		&strategyImpl.BbHaSlma{},
		&strategyImpl.Macd{},
		&strategyImpl.EmaMacd{},
	}

	capital := 100.0
	leverage := 15.0
	takeProfit := 0.7
	stopLoss := 0.2
	st := strategies[18]
	tf := timeframe.MINUTE_5
	smb := symbol.SOLUSDT

	bf := backtest.NewBacktestBotFather(202)
	b := bot.NewBot(tf, st, smb, capital, leverage, takeProfit, stopLoss)

	bf.LoadData(getName(b))

	bf.Randomize()
	startIndex := bf.Index

	bf.Run(b)

	length := bf.Index - startIndex
	days := 0
	switch b.Timeframe {
	case timeframe.HOUR_1:
		days = length / 24
	case timeframe.MINUTE_15:
		days = length / 96
	case timeframe.MINUTE_5:
		days = length / 288
	}
	info := fmt.Sprintf("Name: %s, Capital: %.2f, Pnl: %d/%d, Days: %d\n", b.Name, b.CurrentCapital+b.OrderCapital, b.TotalWins, b.TotalLosses, days)
	bf.PrintChart(bf.Capital, info)

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
