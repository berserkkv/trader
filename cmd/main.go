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

	bbha := strategy.BBHAStrategy{}
	//ha := strategy.HAStrategy{}
	haSmoothed := strategy.HASmoothedStrategy{}
	haema := strategy.HAEMAStrategy{}
	haSmoothedEma := strategy.HASmoothedEMAStrategy{}
	bbha2 := strategy.BBHA2Strategy{}

	haema1m := bot.NewBot(timeframe.MINUTE_1, haema, symbol.SOLUSDT, 100)
	haSmoothedEma1m := bot.NewBot(timeframe.MINUTE_1, haSmoothedEma, symbol.SOLUSDT, 100)
	bbha2_1m := bot.NewBot(timeframe.MINUTE_1, bbha2, symbol.SOLUSDT, 100)

	bbha2_5m := bot.NewBot(timeframe.MINUTE_5, bbha2, symbol.SOLUSDT, 100)

	haema15m := bot.NewBot(timeframe.MINUTE_15, haema, symbol.SOLUSDT, 100)
	haSmoothedEma15m := bot.NewBot(timeframe.MINUTE_15, haSmoothedEma, symbol.SOLUSDT, 100)
	bbha15m := bot.NewBot(timeframe.MINUTE_15, bbha, symbol.SOLUSDT, 100)
	haSmoothed15m := bot.NewBot(timeframe.MINUTE_15, haSmoothed, symbol.SOLUSDT, 100)
	bbha2_15m := bot.NewBot(timeframe.MINUTE_15, bbha2, symbol.SOLUSDT, 100)

	// 1 min
	_, err := service.SaveBot(haema1m)
	if err != nil {
		slog.Debug(err.Error())
	}

	_, err = service.SaveBot(haSmoothedEma1m)
	if err != nil {
		slog.Debug(err.Error())
	}

	_, err = service.SaveBot(bbha2_1m)
	if err != nil {
		slog.Debug(err.Error())
	}

	// 5 min
	_, err = service.SaveBot(bbha2_5m)
	if err != nil {
		slog.Debug(err.Error())
	}

	// 15 min
	_, err = service.SaveBot(haema15m)
	if err != nil {
		slog.Debug(err.Error())
	}

	_, err = service.SaveBot(haSmoothedEma15m)
	if err != nil {
		slog.Debug(err.Error())
	}

	_, err = service.SaveBot(bbha15m)
	if err != nil {
		slog.Debug(err.Error())
	}

	_, err = service.SaveBot(haSmoothed15m)
	if err != nil {
		slog.Debug(err.Error())
	}

	_, err = service.SaveBot(bbha2_15m)
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
