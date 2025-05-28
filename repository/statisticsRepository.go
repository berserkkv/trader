package repository

import (
	"github.com/berserkkv/trader/database"
	"github.com/berserkkv/trader/model"
	"gorm.io/gorm"
)

type StatisticsRepository struct {
	db *gorm.DB
}

func NewStatisticsRepository(db *gorm.DB) *StatisticsRepository {
	return &StatisticsRepository{db: db}
}

func (r *StatisticsRepository) GetAllStatistics() ([]model.Statistics, error) {
	var stats []model.Statistics
	err := database.DB.Find(&stats).Error
	return stats, err
}

func (r *StatisticsRepository) GetStatisticsByTimeframe(tf string) ([]model.Statistics, error) {
	var stats []model.Statistics
	err := database.DB.Where("timeframe = ?", tf).Find(&stats).Error
	return stats, err
}

func (r *StatisticsRepository) GetStatisticsByStrategy(strategy string) ([]model.Statistics, error) {
	var stats []model.Statistics
	err := database.DB.Where("strategy_name = ?", strategy).Find(&stats).Error
	return stats, err
}

func (r *StatisticsRepository) GetStatisticsBySymbol(symbol string) ([]model.Statistics, error) {
	var stats []model.Statistics
	err := database.DB.Where("symbol = ?", symbol).Find(&stats).Error
	return stats, err
}

func (r *StatisticsRepository) SaveAllStatistics(stats *[]model.Statistics) error {
	if len(*stats) == 0 {
		return nil
	}

	result := database.DB.Create(stats)
	return result.Error
}
