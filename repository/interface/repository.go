package repository

import (
	"github.com/berserkkv/trader/bot"
	"github.com/berserkkv/trader/model"
)

type BotRepository interface {
	GetAll(field map[string]interface{}) []*bot.Bot
	GetByID(id int64) *bot.Bot
	Update(bot *bot.Bot) (*bot.Bot, error)
	Create(bot *bot.Bot) (*bot.Bot, error)
	Delete(id int64)
	UpdateAll(bots []bot.Bot) []error
}

type OrderRepository interface {
	GetAll(fields map[string]interface{}) []*model.Order
	GetByID(id int64) *model.Order
	GetAllByBotID(botId int64, sort string) []*model.Order
	Create(order *model.Order) (*model.Order, error)
	Update(order *model.Order) (*model.Order, error)
	Delete(id int64)
	DeleteAllByBotID(botId int64)
}
