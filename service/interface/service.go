package service

import (
	"github.com/berserkkv/trader/bot"
	"github.com/berserkkv/trader/model"
)

type BotService interface {
	GetAll(fields map[string]interface{}) []*bot.Bot
	GetById(id int64) *bot.Bot
	Create(bot *bot.Bot) (*bot.Bot, error)
	Update(bot *bot.Bot) (*bot.Bot, error)
	Delete(id int64) error
	UpdateAll(bots []*bot.Bot) []error
	Stop(id int64) error
	Start(id int64) error
	ClosePosition(id int64) error
	CreateAndAddToBotFather(botReq *model.CreateBotRequest) (*bot.Bot, error)
}

type OrderService interface {
	GetAll(fields map[string]interface{}) []*model.Order
	GetByID(id int64) *model.Order
	GetAllByBotID(botID int64) []*model.Order
	GetStatisticsByBotID(botID int64) []*model.Statistics
	GetAllStatistics() map[string][]*model.Statistics
	Create(order *model.Order) (*model.Order, error)
	Update(order *model.Order) (*model.Order, error)
	Delete(id int64)
	DeleteAllByBotID(botID int64)
}
