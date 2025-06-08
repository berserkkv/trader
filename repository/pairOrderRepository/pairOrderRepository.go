package pairOrderRepository

import (
	"github.com/berserkkv/trader/database"
	"github.com/berserkkv/trader/model"
)

func GetAllOrders(sortOrder string) []model.PairOrder {
	var orders []model.PairOrder
	if sortOrder != "desc" {
		sortOrder = "asc"
	}
	database.DB.Order("closed_time " + sortOrder).Find(&orders)

	return orders
}

func GetAllOrdersByBotIDDesc(botId int64) []model.PairOrder {
	var orders []model.PairOrder
	database.DB.
		Where("bot_id = ?", botId).
		Order("created_time DESC").
		Find(&orders)
	return orders
}

func GetAllOrdersByBotID(botId int64) []model.PairOrder {
	var orders []model.PairOrder
	database.DB.
		Where("bot_id = ?", botId).
		Find(&orders)
	return orders
}

func CreateOrder(order model.PairOrder) model.PairOrder {
	database.DB.Create(&order)
	return order
}

func UpdateOrder(order model.PairOrder) model.PairOrder {
	database.DB.Save(&order)
	return order
}
