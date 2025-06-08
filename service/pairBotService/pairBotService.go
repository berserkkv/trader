package pairBotService

import (
	"errors"
	"github.com/berserkkv/trader/model/enum/symbol"
	"github.com/berserkkv/trader/model/enum/timeframe"
	"github.com/berserkkv/trader/pairBot"
	"github.com/berserkkv/trader/pairBot/pairBotFather"
	"github.com/berserkkv/trader/repository/pairBotRepository"
	"github.com/berserkkv/trader/service/connector"
	"log/slog"
)

func SaveBot(bot *pairBot.PairBot) (*pairBot.PairBot, error) {
	if bot == nil {
		return nil, errors.New("bot not saved, nil bot")
	}
	if bot.Name == "" {
		return nil, errors.New("bot not saved, name is empty")
	}
	if bot.Symbol1 == "" {
		return nil, errors.New("bot not saved, symbol1 is empty")
	}
	if bot.Symbol2 == "" {
		return nil, errors.New("bot not saved, symbol2 is empty")
	}

	return pairBotRepository.CreateBot(bot)
}

func GetAllBots(fields map[string]interface{}) []pairBot.PairBot {
	bots := pairBotRepository.GetAllBots(fields)
	for i := range bots {
		bots[i].Connector = connector.BinanceConnector{}
	}
	return bots
}

func GetBotById(id int64) *pairBot.PairBot {
	return pairBotRepository.GetBotById(id)
}

func UpdateAllBots(bots []*pairBot.PairBot) []error {
	return pairBotRepository.UpdateAllBots(bots)
}

func UpdateBot(bot *pairBot.PairBot) (*pairBot.PairBot, error) {
	return pairBotRepository.UpdateBot(bot)
}

func StopBot(botId int64) error {
	b, err := pairBotFather.GetPairBotFather().StopBot(botId)

	if err != nil {
		slog.Error("error stopping bot", "id", botId)
		return err
	}

	_, err = pairBotRepository.UpdateBot(b)
	if err != nil {
		slog.Error("error updating bot table", "id", botId, "error", err)
		return err
	}

	return nil
}

func StartBot(botId int64) error {
	b, err := pairBotFather.GetPairBotFather().StartBot(botId)
	if err != nil {
		slog.Error("Error starting bot", "id", botId)
		return err
	}

	_, err = pairBotRepository.UpdateBot(b)
	if err != nil {
		slog.Error("error starting bot", "id", botId, "error", err)
		return err
	}

	return nil
}

func ClosePosition(botId int64) {
	pairBotFather.GetPairBotFather().ClosePosition(botId)

	slog.Error("Closed Position", "id", botId)
}

func CreateBot(tradingSymbol1, tradingSymbol2, tradingTimeFrame string, capital, leverage, takeProfit, stopLoss float64) (*pairBot.PairBot, error) {

	tf, err := timeframe.Parse(tradingTimeFrame)
	if err != nil {
		return nil, err
	}
	smb1, err := symbol.Parse(tradingSymbol1)
	if err != nil {
		return nil, err
	}

	smb2, err := symbol.Parse(tradingSymbol2)
	if err != nil {
		return nil, err
	}

	b := pairBot.NewPairBot(smb1, smb2, tf, capital, leverage, takeProfit, stopLoss)

	savedBot, err := pairBotRepository.CreateBot(b)
	if err != nil {
		slog.Error("error saving bot", "error", err)
	}

	pairBotFather.GetPairBotFather().AddBot(savedBot)

	return savedBot, err
}
