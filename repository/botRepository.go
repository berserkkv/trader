package repository

import (
	"github.com/berserkkv/trader/bot"
	"github.com/berserkkv/trader/database"
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/strategy"
	"log/slog"
)

func GetAllBots(fields map[string]interface{}) []bot.Bot {
	var bots []bot.Bot
	database.DB.
		Where(fields).
		Order("is_not_active").
		Find(&bots)
	return bots
}

func GetBotById(id int64) *bot.Bot {
	var b bot.Bot
	database.DB.First(&b, id)
	b.Strategy = strategy.GetStrategy(b.StrategyName)
	return &b
}

func CreateBot(bot *bot.Bot) (*bot.Bot, error) {
	result := database.DB.Create(&bot)
	if result.Error != nil {
		slog.Error("Failed to create bot", "error", result.Error)
		return nil, result.Error
	}
	return bot, nil
}

func UpdateBot(bot *bot.Bot) (*bot.Bot, error) {
	result := database.DB.Save(&bot)
	if result.Error != nil {
		slog.Error("Failed to update bot", "error", result.Error)
		return nil, result.Error
	}
	return bot, nil
}

func UpdateAllBots(bots []*bot.Bot) []error {
	var errs []error
	for _, b := range bots {
		if err := database.DB.Save(b).Error; err != nil {
			slog.Error("Failed to update bot", "name", b.Name, "error", err)
			errs = append(errs, err)
		}
	}

	return errs
}

func DeleteBotById(id int64) {
	if err := database.DB.Where("bot_id = ?", id).Delete(&model.Order{}).Error; err != nil {
		slog.Error("Failed to delete orders", "error", err)
		return
	}

	if err := database.DB.Delete(&bot.Bot{}, id).Error; err != nil {
		slog.Error("Failed to delete bot", "error", err)
	}
}
