package sqliteImpl

import (
	"github.com/berserkkv/trader/model"
	"gorm.io/gorm"
)

type OrderRepositoryImpl struct {
	db *gorm.DB
}

func NewOrderRepositoryImpl(db *gorm.DB) *OrderRepositoryImpl {
	return &OrderRepositoryImpl{db: db}
}

func (repo *OrderRepositoryImpl) GetAll(fields map[string]interface{}) []*model.Order {
	var orders []*model.Order
	repo.db.Where(fields).
		Find(&model.Order{})

	return orders
}

func (repo *OrderRepositoryImpl) GetByID(id int64) *model.Order {
	var o model.Order
	repo.db.First(&o, "id = ?", id)
	return &o
}

func (repo *OrderRepositoryImpl) GetAllByBotID(botId int64, sort string) []*model.Order {
	var orders []*model.Order
	repo.db.
		Where("bot_id = ?", botId).
		Order("created_time " + sort).
		Find(&orders)
	return orders
}

func (repo *OrderRepositoryImpl) Create(order *model.Order) (*model.Order, error) {
	err := repo.db.Create(&order).Error
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (repo *OrderRepositoryImpl) Update(order *model.Order) (*model.Order, error) {
	err := repo.db.Save(&order).Error
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (repo *OrderRepositoryImpl) Delete(id int64) {
	repo.db.Delete(&model.Order{}, id)
}

func (repo *OrderRepositoryImpl) DeleteAllByBotID(botId int64) {
	repo.db.Where("bot_id = ?", botId).Delete(&model.Order{})
}
