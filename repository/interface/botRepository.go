package repository

import "github.com/berserkkv/trader/bot"

type BotRepository interface {
	GetAll(field map[string]interface{}) []*bot.Bot
	GetByID(id int64) *bot.Bot
	Update(bot *bot.Bot) (*bot.Bot, error)
	Create(bot *bot.Bot) (*bot.Bot, error)
	Delete(id int64)
	UpdateAll(bots []bot.Bot) []error
}
