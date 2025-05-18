package main

import (
	"github.com/berserkkv/trader/bot"
	"github.com/berserkkv/trader/bot/botFather"
	"github.com/berserkkv/trader/controller"
	"github.com/berserkkv/trader/database"
	"github.com/berserkkv/trader/model/enum/symbol"
	"github.com/berserkkv/trader/model/enum/timeframe"
	"github.com/berserkkv/trader/service"
	strategy "github.com/berserkkv/trader/strategy/bbha"
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

	bbha := strategy.BBHAStrategy{}

	b := bot.NewBot(timeframe.MINUTE_5, bbha, symbol.SOLUSDT, 1000)
	b2 := bot.NewBot(timeframe.MINUTE_1, bbha, symbol.BTCUSDT, 1000)

	_, err := service.SaveBot(b)
	if err != nil {
		slog.Error(err.Error())
	}
	_, err = service.SaveBot(b2)

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
