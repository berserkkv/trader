package repository

import (
	"github.com/berserkkv/trader/database"
	"github.com/berserkkv/trader/model"
)

func GelAllOrders() []model.Order {
	var orders []model.Order
	database.DB.Find(&orders)
	return orders
}

func GetAllOrdersByBotID(botId int64) []model.Order {
	var orders []model.Order
	database.DB.Where("bot_id = ?", botId).Find(&orders)
	return orders
}

func CreateOrder(order model.Order) model.Order {
	database.DB.Create(&order)
	return order
}

func UpdateOrder(order model.Order) model.Order {
	database.DB.Save(&order)
	return order
}
