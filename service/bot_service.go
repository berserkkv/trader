package service

import (
	"errors"
	"github.com/berserkkv/trader/bot"
	"github.com/berserkkv/trader/bot/botFather"
	"github.com/berserkkv/trader/model/enum/symbol"
	"github.com/berserkkv/trader/model/enum/timeframe"
	"github.com/berserkkv/trader/repository"
	"github.com/berserkkv/trader/strategy"
	"log/slog"
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

	return repository.CreateBot(bot)
}

func GetAllBots() []bot.Bot {
	bots := repository.GetAllBots()
	for i := range bots {
		bots[i].Strategy = strategy.GetStrategy(bots[i].StrategyName)
	}
	return bots
}

func UpdateAllBots(bots []*bot.Bot) []error {
	return repository.UpdateAllBots(bots)
}

func UpdateBot(bot *bot.Bot) (*bot.Bot, error) {
	return repository.UpdateBot(bot)
}

func StopBot(botId int64) error {
	b := botFather.GetBotFather().Bots()[botId]

	if b == nil {
		slog.Error("error stopping bot, bot is nil", "id", botId)
		return errors.New("error stopping bot, bot is nil")
	}
	if b.InPos {
		slog.Error("error stopping bot, bot is in position", "id", botId)
		return errors.New("error stopping bot, bot is in position")
	}

	b.IsNotActive = true

	_, err := repository.UpdateBot(b)
	if err != nil {
		slog.Error("error updating bot table", "id", botId, "error", err)
		return err
	}

	return nil
}

func StartBot(botId int64) error {
	b := botFather.GetBotFather().Bots()[botId]
	if b == nil {
		slog.Error("Error starting bot, bot is nil", "id", botId)
		return errors.New("error starting bot, bot is nil")
	}

	if b.CurrentCapital <= 10 {
		slog.Error("error starting bot, bot's capital <= 10", "id", botId)
		return errors.New("error starting bot, bot's capital <= 10")
	}

	b.IsNotActive = false

	_, err := repository.UpdateBot(b)
	if err != nil {
		slog.Error("error starting bot", "id", botId, "error", err)
		return err
	}

	return nil
}

func CreateBot(tradingSymbol, strategyName, tradingTimeFrame string, capital float64) (*bot.Bot, error) {

	tradingStrategy := strategy.GetStrategy(strategyName)
	tf, err := timeframe.Parse(tradingTimeFrame)
	if err != nil {
		return nil, err
	}
	smb, err := symbol.Parse(tradingSymbol)
	if err != nil {
		return nil, err
	}

	b := bot.NewBot(tf, tradingStrategy, smb, capital)

	savedBot, err := repository.CreateBot(b)
	if err != nil {
		slog.Error("error saving bot", "error", err)
	}

	botFather.GetBotFather().AddBot(savedBot)

	return savedBot, err
}
