package service

import (
	"errors"
	"fmt"
	"github.com/berserkkv/trader/bot"
	"github.com/berserkkv/trader/repository"
	"github.com/berserkkv/trader/strategy"
)

func CreateBot(bot *bot.Bot) (*bot.Bot, error) {
	if bot == nil {
		return nil, errors.New("Bot not saved, nil bot")
	}
	if bot.Name == "" {
		return nil, errors.New("Bot not saved, name is empty")
	}
	if bot.Symbol == "" {
		fmt.Println(bot.Symbol)
		return nil, errors.New("Bot not saved, symbol is empty")
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
