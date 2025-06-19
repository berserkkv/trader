package serviceImpl

import (
	"errors"
	"github.com/berserkkv/trader/bot"
	"github.com/berserkkv/trader/bot/botFather"
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/model/enum/symbol"
	"github.com/berserkkv/trader/model/enum/timeframe"
	repository "github.com/berserkkv/trader/repository/interface"
	"github.com/berserkkv/trader/service/connector"
	service "github.com/berserkkv/trader/service/interface"
	"github.com/berserkkv/trader/strategy"
	"log/slog"
)

type BotServiceImpl struct {
	r  repository.BotRepository
	bf *botFather.BotFather
}

func NewBotService(repo repository.BotRepository, bf *botFather.BotFather) *BotServiceImpl {
	return &BotServiceImpl{r: repo, bf: bf}
}

func (s *BotServiceImpl) GetAll(field map[string]interface{}) []*bot.Bot {
	bots := s.r.GetAll(field)

	for i := range bots {
		bots[i].Strategy = strategy.GetStrategy(bots[i].StrategyName)
		bots[i].Connector = &connector.BinanceConnector{}
	}
	return bots
}

func (s *BotServiceImpl) GetById(id int64) *bot.Bot {
	return s.r.GetByID(id)
}

func (s *BotServiceImpl) Create(bot *bot.Bot) (*bot.Bot, error) {
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

	return s.r.Create(bot)
}

func (s *BotServiceImpl) Update(bot *bot.Bot) (*bot.Bot, error) {
	if bot == nil {
		return nil, errors.New("bot not updated, nil bot")
	}

	updatedBot, err := s.r.Update(bot)
	if err != nil {
		return nil, err
	}
	return updatedBot, nil
}

func (s *BotServiceImpl) Delete(id int64) error {
	s.r.Delete(id)
	s.bf.DeleteBot(id)
	return nil
}

func (s *BotServiceImpl) UpdateAll(bots []*bot.Bot) []error {
	return s.r.UpdateAll(bots)
}

func (s *BotServiceImpl) Stop(id int64) error {
	b, err := s.bf.StopBot(id)

	if err != nil {
		slog.Error("error stopping bot", "id", id)
		return err
	}

	_, err = s.r.Update(b)
	if err != nil {
		slog.Error("error updating bot table", "id", id, "error", err)
		return err
	}

	return nil
}

func (s *BotServiceImpl) Start(id int64) error {
	b, err := s.bf.StartBot(id)
	if err != nil {
		slog.Error("Error starting bot", "id", id)
		return err
	}

	_, err = s.r.Update(b)
	if err != nil {
		slog.Error("error starting bot", "id", id, "error", err)
		return err
	}

	return nil
}

func (s *BotServiceImpl) ClosePosition(id int64) error {
	s.bf.ClosePosition(id)

	slog.Info("closed position", "id", id)

	return nil
}

func (s *BotServiceImpl) CreateAndAddToBotFather(botReq *model.CreateBotRequest) (*bot.Bot, error) {
	if botReq == nil {
		return nil, errors.New("bot not saved, nil botRequest")
	}

	tradingStrategy := strategy.GetStrategy(botReq.Strategy)
	if tradingStrategy == nil {
		return nil, errors.New("strategy not found")
	}

	tf, err := timeframe.Parse(botReq.Timeframe)
	if err != nil {
		return nil, err
	}
	smb, err := symbol.Parse(botReq.Symbol)
	if err != nil {
		return nil, err
	}

	b := bot.NewBot(tf, tradingStrategy, smb, botReq.Capital, botReq.Leverage, botReq.TakeProfit, botReq.StopLoss)

	savedBot, err := s.r.Create(b)
	if err != nil {
		slog.Error("error saving bot", "error", err)
	}

	s.bf.AddBot(savedBot)

	return savedBot, err

}

var _ service.BotService = (*BotServiceImpl)(nil)
