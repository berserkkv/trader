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
	"time"
)

func main() {
	cnf := config.Load()
	logger.Init(cnf.Logger.Level, cnf.Env)

	database.ConnectDB()

	bf := botFather.GetBotFather()

	go runBothFather(bf)

	controller.Register()

	slog.Info("Server started on port: 8080")

}

func runBothFather(bf *botFather.BotFather) {
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
		strategy.Random{},
	}

	tfs := []timeframe.Timeframe{
		timeframe.MINUTE_1,
		timeframe.MINUTE_15,
	}

	smbs := []symbol.Symbol{
		symbol.SOLUSDT,
		symbol.ETHUSDT,
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

	bots := service.GetAllBots(map[string]interface{}{})

	go StartStatisticsScheduler()

	for i := range bots {
		bf.AddBot(&bots[i])
		if bots[i].InPos {
			bf.IncreaseTotalBotsInOrder()
		}
	}

	go bf.CheckAndStartMonitoring()
	bf.Start()
}

func StartStatisticsScheduler() {
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for {
			service.SaveStatistics()
			<-ticker.C // wait for next tick
		}
	}()
}
