package service

import (
	"errors"
	"github.com/berserkkv/trader/bot"
	"github.com/berserkkv/trader/bot/botFather"
	"github.com/berserkkv/trader/model/enum/symbol"
	"github.com/berserkkv/trader/model/enum/timeframe"
	"github.com/berserkkv/trader/repository"
	"github.com/berserkkv/trader/service/connector"
	"github.com/berserkkv/trader/strategy"
	"log/slog"
	"sort"
)

func SaveBot(bot *bot.Bot) (*bot.Bot, error) {
	if bot == nil {
		return nil, errors.New("bot not saved, nil bot")
	}
	if bot.Name == "" {
		return nil, errors.New("bot not saved, name is empty")
	}
	if bot.Symbol == "" {
		return nil, errors.New("bot not saved, symbol is empty")
	}
	if bot.Strategy == nil {
		return nil, errors.New("bot not saved, strategy is nil")
	}

	return repository.CreateBot(bot)
}

func GetAllBots(fields map[string]interface{}) []bot.Bot {
	bots := repository.GetAllBots(fields)
	for i := range bots {
		bots[i].Strategy = strategy.GetStrategy(bots[i].StrategyName)
		bots[i].Connector = connector.BinanceConnector{}
	}
	return bots
}

func GetBotById(id int64) *bot.Bot {
	return repository.GetBotById(id)
}

func UpdateAllBots(bots []*bot.Bot) []error {
	return repository.UpdateAllBots(bots)
}

func UpdateBot(bot *bot.Bot) (*bot.Bot, error) {
	return repository.UpdateBot(bot)
}

func StopBot(botId int64) error {
	b, err := botFather.GetBotFather().StopBot(botId)

	if err != nil {
		slog.Error("error stopping bot", "id", botId)
		return err
	}

	_, err = repository.UpdateBot(b)
	if err != nil {
		slog.Error("error updating bot table", "id", botId, "error", err)
		return err
	}

	return nil
}

func StartBot(botId int64) error {
	b, err := botFather.GetBotFather().StartBot(botId)
	if err != nil {
		slog.Error("Error starting bot", "id", botId)
		return err
	}

	_, err = repository.UpdateBot(b)
	if err != nil {
		slog.Error("error starting bot", "id", botId, "error", err)
		return err
	}

	return nil
}

func CreateBot(tradingSymbol, strategyName, tradingTimeFrame string, capital, leverage, takeProfit, stopLoss float64) (*bot.Bot, error) {

	tradingStrategy := strategy.GetStrategy(strategyName)
	if tradingStrategy == nil {
		return nil, errors.New("strategy not found")
	}

	tf, err := timeframe.Parse(tradingTimeFrame)
	if err != nil {
		return nil, err
	}
	smb, err := symbol.Parse(tradingSymbol)
	if err != nil {
		return nil, err
	}

	b := bot.NewBot(tf, tradingStrategy, smb, capital, leverage, takeProfit, stopLoss)

	savedBot, err := repository.CreateBot(b)
	if err != nil {
		slog.Error("error saving bot", "error", err)
	}

	botFather.GetBotFather().AddBot(savedBot)

	return savedBot, err
}

type Stats struct {
	StrategyName string
	Total        float64
	M1           float64
	M5           float64
	M15          float64
	BotCount     int
}

func GetBotStatistics() []Stats {
	bots := repository.GetAllBots(nil)
	statistics := map[string]*Stats{}

	for _, b := range bots {
		s, exists := statistics[b.StrategyName]
		if !exists {
			s = &Stats{StrategyName: b.StrategyName}
			statistics[b.StrategyName] = s
		}

		balance := b.CurrentCapital + b.OrderCapital
		s.Total += balance
		s.BotCount++

		switch b.Timeframe {
		case timeframe.MINUTE_1:
			s.M1 += balance
		case timeframe.MINUTE_5:
			s.M5 += balance
		case timeframe.MINUTE_15:
			s.M15 += balance
		}
	}

	// Convert map to slice
	result := make([]Stats, 0, len(statistics))
	for _, stat := range statistics {
		result = append(result, *stat)
	}

	// Sort by Total descending (example)
	sort.Slice(result, func(i, j int) bool {
		return result[i].Total > result[j].Total
	})

	return result
}
