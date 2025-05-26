package main

import (
	"github.com/berserkkv/trader/bot"
	"github.com/berserkkv/trader/bot/botFather"
	"github.com/berserkkv/trader/controller"
	"github.com/berserkkv/trader/database"
	"github.com/berserkkv/trader/model/enum/symbol"
	"github.com/berserkkv/trader/model/enum/timeframe"
	"github.com/berserkkv/trader/service"
	"github.com/berserkkv/trader/strategy"
	"github.com/berserkkv/trader/tools/config"
	logger "github.com/berserkkv/trader/tools/log"
	"log/slog"
)

func main() {
	cnf := config.Load()
	logger.Init(cnf.Logger.Level, cnf.Env)

	database.ConnectDB()

	go runBothFather()
	controller.Register()

	slog.Info("Server started on port: 8080")

}

func runBothFather() {
	bf := botFather.GetBotFather()

	capital := 100.0
	leverage := 10.0
	takeProfit := 1.0
	stopLoss := 0.75

	sts := []strategy.Strategy{
		strategy.HAStrategy{},
		strategy.HASmoothedStrategy{},
		strategy.HAEMAStrategy{},
		strategy.BBHAStrategy{},
		strategy.BBHA2Strategy{},
		strategy.HASmoothedEMAStrategy{},
		strategy.BBHA3{},
	}
	tfs := []timeframe.Timeframe{
		timeframe.MINUTE_1,
		timeframe.MINUTE_5,
		timeframe.MINUTE_15,
	}

	smbs := []symbol.Symbol{
		symbol.SOLUSDT,
		symbol.ETHUSDT,
		//symbol.BTCUSDT,
		symbol.BNBUSDT,
	}

	for _, tf := range tfs {
		for _, st := range sts {
			for _, smb := range smbs {
				b := bot.NewBot(tf, st, smb, capital, leverage, takeProfit, stopLoss)

				_, err := service.SaveBot(b)
				if err != nil {
					slog.Debug("Failed to save bot: ", err)
				}
			}

		}
	}

	bots := service.GetAllBots()

	for i := range bots {
		bf.AddBot(&bots[i])
		if bots[i].InPos {
			bf.IncreaseTotalBotsInOrder()
		}
	}

	go bf.CheckAndStartMonitoring()
	bf.Start()
}
