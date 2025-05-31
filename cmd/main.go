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

	bf := botFather.GetBotFather()

	go runBothFather(bf)

	controller.Register()

	slog.Info("Server started on port: 8080")

}

func runBothFather(bf *botFather.BotFather) {
	leverage := 10.0
	takeProfit := 0.5
	stopLoss := 0.2
	takeProfit2 := 1.0
	stopLoss2 := 0.5

	sts := []strategy.Strategy{
		strategy.HAStrategy{},
		//strategy.HASmoothedStrategy{},
		//strategy.HAEMAStrategy{},
		//strategy.BBHAStrategy{},
		//strategy.BBHA2Strategy{},
		strategy.BBHA3{},
		//strategy.HASmoothedEMAStrategy{},
		//strategy.Random{},
	}

	tfs := []timeframe.Timeframe{
		timeframe.MINUTE_1,
	}

	smbs := []symbol.Symbol{
		symbol.SOLUSDT,
	}

	for _, tf := range tfs {
		for _, st := range sts {
			for _, smb := range smbs {
				b := bot.NewBot2(tf, st, smb, leverage, takeProfit, stopLoss)
				b2 := bot.NewBot2(tf, st, smb, leverage, takeProfit2, stopLoss2)
				_, err := service.SaveBot(b)
				if err != nil {
					slog.Debug("Failed to save bot: ", err)
				}
				_, err = service.SaveBot(b2)
				if err != nil {
					slog.Debug("Failed to save bot: ", err)
				}
			}
		}
	}

	bots := service.GetAllBots(map[string]interface{}{})

	for i := range bots {
		bf.AddBot(&bots[i])
		if bots[i].InPos {
			bf.IncreaseTotalBotsInOrder()
		}
	}

	go bf.CheckAndStartMonitoring()
	bf.Start()
}
