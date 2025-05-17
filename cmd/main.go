package main

import (
	"github.com/berserkkv/trader/bot"
	"github.com/berserkkv/trader/bot/botFather"
	"github.com/berserkkv/trader/database"
	"github.com/berserkkv/trader/model/enum/symbol"
	"github.com/berserkkv/trader/model/enum/timeframe"
	"github.com/berserkkv/trader/service"
	strategy "github.com/berserkkv/trader/strategy/bbha"
	"github.com/berserkkv/trader/tools/config"
	logger "github.com/berserkkv/trader/tools/log"
	"log/slog"
	"time"
)

func main() {
	// init
	cnf := config.Load()
	logger.Init(cnf.Logger.Level, cnf.Env)

	database.ConnectDB()
	time.Sleep(1 * time.Second)

	slog.Info("Server started on port: 8080")

	bf := botFather.GetBotFather()

	bbha := strategy.BBHAStrategy{}

	b := bot.NewBot(timeframe.MINUTE_5, bbha, symbol.SOLUSDT)
	b2 := bot.NewBot(timeframe.MINUTE_1, bbha, symbol.BTCUSDT)

	_, err := service.CreateBot(b)
	if err != nil {
		slog.Error(err.Error())
	}
	_, err = service.CreateBot(b2)

	bots := service.GetAllBots()

	for i := range bots {
		bf.AddBot(&bots[i])
	}

	bf.Start()

}
