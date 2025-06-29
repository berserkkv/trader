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
	"sort"
)

type BotServiceImpl struct {
	r  repository.BotRepository
	bf *botFather.BotFather
}

func NewBotService(repo repository.BotRepository, bf *botFather.BotFather) *BotServiceImpl {
	return &BotServiceImpl{r: repo, bf: bf}
}

func (s *BotServiceImpl) GetAll(field map[string]interface{}) []*bot.Bot {
	bots := s.bf.Bots()
	sort.Slice(bots, func(i, j int) bool {
		if bots[i].CurrentCapital+bots[i].OrderCapital != bots[j].CurrentCapital+bots[j].OrderCapital {
			return bots[i].CurrentCapital+bots[i].OrderCapital > bots[j].CurrentCapital+bots[j].OrderCapital
		}
		first, second := 0, 0
		switch bots[i].Timeframe {
		case timeframe.MINUTE_1:
			first = 1
		case timeframe.MINUTE_5:
			first = 5
		case timeframe.MINUTE_15:
			first = 15
		case timeframe.MINUTE_30:
			first = 30
		case timeframe.HOUR_1:
			first = 60
		case timeframe.DAY:
			first = 3600
		}
		switch bots[j].Timeframe {
		case timeframe.MINUTE_1:
			second = 1
		case timeframe.MINUTE_5:
			second = 5
		case timeframe.MINUTE_15:
			second = 15
		case timeframe.MINUTE_30:
			second = 30
		case timeframe.HOUR_1:
			second = 60
		case timeframe.DAY:
			second = 3600
		}
		return first < second
	})

	return bots
}

func (s *BotServiceImpl) GetAllFromRepo(field map[string]interface{}) []*bot.Bot {
	bots := s.r.GetAll(field)

	for i := range bots {
		bots[i].Strategy = strategy.GetStrategy(bots[i].StrategyName)
		bots[i].Connector = &connector.BinanceConnector{}
	}
	return bots
}

func (s *BotServiceImpl) GetById(id int64) *bot.Bot {
	b, err := s.bf.GetBot(id)
	if err != nil {
		slog.Error("Bot not found", "id", id)
	}
	return b
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

func (s *BotServiceImpl) UpdateWithRequest(botUpdateRequest *model.BotUpdateRequest) (*bot.Bot, error) {
	b, err := s.bf.GetBot(botUpdateRequest.Id)
	if err != nil {
		slog.Error("error getting bot", "id", botUpdateRequest.Id)
		return nil, err
	}

	b.TakeProfit = botUpdateRequest.TakeProfit
	b.StopLoss = botUpdateRequest.StopLoss
	b.IsTrailingStopActive = botUpdateRequest.IsTrailingStopActive

	_, err = s.r.Update(b)
	if err != nil {
		slog.Error("error updating bot table", "id", b.Id, "error", err)
		return nil, err
	}
	return b, nil
}

var _ service.BotService = (*BotServiceImpl)(nil)
