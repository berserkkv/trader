package sqliteImpl

import (
	"github.com/berserkkv/trader/bot"
	"github.com/berserkkv/trader/database"
	"github.com/berserkkv/trader/model"
	repository "github.com/berserkkv/trader/repository/interface"
	"github.com/berserkkv/trader/strategy"
	"gorm.io/gorm"
	"log/slog"
)

type BotRepository struct {
	db *gorm.DB
}

func NewBotRepository(db *gorm.DB) *BotRepository {
	return &BotRepository{db: db}
}

func (repo *BotRepository) GetAll(fields map[string]interface{}) []*bot.Bot {
	var bots []*bot.Bot
	repo.db.
		Where(fields).
		Order("is_not_active").
		Find(&bots)
	return bots
}

func (repo *BotRepository) GetByID(id int64) *bot.Bot {
	var b bot.Bot
	repo.db.First(&b, id)
	b.Strategy = strategy.GetStrategy(b.StrategyName)
	return &b
}

func (repo *BotRepository) Update(bot *bot.Bot) (*bot.Bot, error) {
	result := repo.db.Save(&bot)
	if result.Error != nil {
		slog.Error("Failed to update bot", "error", result.Error)
		return nil, result.Error
	}
	return bot, nil
}

func (repo *BotRepository) Create(bot *bot.Bot) (*bot.Bot, error) {
	result := repo.db.Create(&bot)
	if result.Error != nil {
		slog.Error("Failed to create bot", "error", result.Error)
		return nil, result.Error
	}
	return bot, nil
}

func (repo *BotRepository) Delete(id int64) {
	if err := repo.db.Where("bot_id = ?", id).Delete(&model.Order{}).Error; err != nil {
		slog.Error("Failed to delete orders", "error", err)
		return
	}

	if err := database.DB.Delete(&bot.Bot{}, id).Error; err != nil {
		slog.Error("Failed to delete bot", "error", err)
	}
}

func (repo *BotRepository) UpdateAll(bots []*bot.Bot) []error {
	var errs []error
	for _, b := range bots {
		if err := repo.db.Save(b).Error; err != nil {
			slog.Error("Failed to update bot", "name", b.Name, "error", err)
			errs = append(errs, err)
		}
	}

	return errs
}

var _ repository.BotRepository = (*BotRepository)(nil)
