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

	//go runBothFather()
	controller.Register()

	slog.Info("Server started on port: 8080")

}

func runBothFather() {
	bf := botFather.GetBotFather()

	bbha := strategy.BBHAStrategy{}
	ha := strategy.HAStrategy{}
	haSmoothed := strategy.HASmoothedStrategy{}

	bbha1m := bot.NewBot(timeframe.MINUTE_1, bbha, symbol.SOLUSDT, 1000)
	ha1m := bot.NewBot(timeframe.MINUTE_1, ha, symbol.SOLUSDT, 1000)
	haSmoothed1m := bot.NewBot(timeframe.MINUTE_1, haSmoothed, symbol.SOLUSDT, 1000)

	bbha15m := bot.NewBot(timeframe.MINUTE_15, bbha, symbol.SOLUSDT, 1000)
	ha15m := bot.NewBot(timeframe.MINUTE_15, ha, symbol.SOLUSDT, 1000)
	haSmoothed15m := bot.NewBot(timeframe.MINUTE_15, haSmoothed, symbol.SOLUSDT, 1000)

	_, err := service.SaveBot(bbha1m)
	if err != nil {
		slog.Debug(err.Error())
	}

	_, err = service.SaveBot(ha1m)
	if err != nil {
		slog.Debug(err.Error())
	}

	_, err = service.SaveBot(haSmoothed1m)
	if err != nil {
		slog.Debug(err.Error())
	}

	_, err = service.SaveBot(bbha15m)
	if err != nil {
		slog.Debug(err.Error())
	}
	_, err = service.SaveBot(ha15m)
	if err != nil {
		slog.Debug(err.Error())
	}
	_, err = service.SaveBot(haSmoothed15m)
	if err != nil {
		slog.Debug(err.Error())
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
