package main

import (
	"github.com/berserkkv/trader/bot"
	"github.com/berserkkv/trader/bot/botFather"
	"github.com/berserkkv/trader/container"
	"github.com/berserkkv/trader/pairBot"
	"github.com/berserkkv/trader/pairBot/pairBotFather"
	"github.com/berserkkv/trader/router"
	"github.com/berserkkv/trader/service/pairBotService"
	strategyImpl "github.com/berserkkv/trader/strategy/impl"

	"github.com/berserkkv/trader/model/enum/symbol"
	"github.com/berserkkv/trader/model/enum/timeframe"
	service "github.com/berserkkv/trader/service/interface"
	"github.com/berserkkv/trader/strategy"
	"log/slog"
)

func main() {
	cnt := container.New()

	go runBothFather(cnt.BotFather, cnt.BotService)

	router.Register(cnt.BotController, cnt.OrderController)
}

func runBothFather(bf *botFather.BotFather, service service.BotService) {
	capital := 100.0
	leverage := 10.0
	takeProfit := 1.0
	stopLoss := 0.15

	sts := []strategy.Strategy{
		//strategy.HAStrategy{},
		//strategy.HASmoothedStrategy{},
		//strategy.HAEMAStrategy{},
		//strategy.BBHAStrategy{},
		//strategy.BBHA2Strategy{},
		//strategy.BBHA3{},
		//strategy.HASmoothedEMAStrategy{},
		//strategy.Random{},
		//strategy.Supertrend{},
		//strategy.Supertrend2{},
		//&strategy.HaSlma{},
		//&strategy.TripleCandle{},
		//&strategy.SlmaBb{},
		//&strategy.TripleCandleSlma{},
		&strategy.SlmaEmaBb{},
		&strategyImpl.SlmaEmaBb2{},
	}

	tfs := []timeframe.Timeframe{
		timeframe.MINUTE_1,
		timeframe.MINUTE_15,
		timeframe.HOUR_1,
	}

	smbs := []symbol.Symbol{
		symbol.SOLUSDT,
	}

	for _, tf := range tfs {
		for _, st := range sts {
			for _, smb := range smbs {
				b := bot.NewBot(tf, st, smb, capital, leverage, takeProfit, stopLoss)
				_, err := service.Create(b)
				if err != nil {
					slog.Debug("Failed to save bot: ", err)
				}
			}
		}
	}

	bots := service.GetAll(map[string]interface{}{})

	for i := range bots {
		bf.AddBot(bots[i])
		if bots[i].InPos {
			bf.IncreaseTotalBotsInOrder()
		}
	}

	go bf.CheckAndStartMonitoring()
	bf.Start()
}

func runPairBots() {
	bf := pairBotFather.GetPairBotFather()

	capital := 100.0
	leverage := 10.0
	takeProfit := 1.5
	stopLoss := 0.5

	tfs := []timeframe.Timeframe{
		timeframe.MINUTE_1,
		timeframe.MINUTE_15,
		//timeframe.HOUR_1,
	}

	smbs := [][]symbol.Symbol{
		//[]symbol.Symbol{symbol.BTCUSDT, symbol.ETHUSDT},
		//[]symbol.Symbol{symbol.BTCUSDT, symbol.SOLUSDT},
		[]symbol.Symbol{symbol.ETHUSDT, symbol.SOLUSDT},
		[]symbol.Symbol{symbol.SOLUSDT, symbol.AVAXUSDT},
	}

	for _, tf := range tfs {
		for _, smb := range smbs {
			b := pairBot.NewPairBot(smb[0], smb[1], tf, capital, leverage, takeProfit, stopLoss)
			_, err := pairBotService.SaveBot(b)
			if err != nil {
				slog.Debug("Failed to save bot: ", err)
			}
		}
	}

	bots := pairBotService.GetAllBots(map[string]interface{}{})

	for i := range bots {
		bf.AddBot(&bots[i])
		if bots[i].InPos {
			bf.IncreaseTotalBotsInOrder()
		}
	}

	go bf.CheckAndStartMonitoring()
	bf.Start()
}
