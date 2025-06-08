package pairBotRepository

import (
	"github.com/berserkkv/trader/database"
	"github.com/berserkkv/trader/pairBot"
	"log/slog"
)

func GetAllBots(fields map[string]interface{}) []pairBot.PairBot {
	var bots []pairBot.PairBot
	database.DB.
		Where(fields).
		Order("is_not_active").
		Find(&bots)
	return bots
}

func GetBotById(id int64) *pairBot.PairBot {
	var b pairBot.PairBot
	database.DB.First(&b, id)
	return &b
}

func CreateBot(bot *pairBot.PairBot) (*pairBot.PairBot, error) {
	result := database.DB.Create(&bot)
	if result.Error != nil {
		slog.Error("Failed to create bot", "error", result.Error)
		return nil, result.Error
	}
	return bot, nil
}

func UpdateBot(bot *pairBot.PairBot) (*pairBot.PairBot, error) {
	result := database.DB.Save(&bot)
	if result.Error != nil {
		slog.Error("Failed to update bot", "error", result.Error)
		return nil, result.Error
	}
	return bot, nil
}

func UpdateAllBots(bots []*pairBot.PairBot) []error {
	var errs []error
	for _, b := range bots {
		if err := database.DB.Save(b).Error; err != nil {
			slog.Error("Failed to update bot", "name", b.Name, "error", err)
			errs = append(errs, err)
		}
	}

	return errs
}
