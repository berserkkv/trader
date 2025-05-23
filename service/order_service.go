package service

import (
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/repository"
)

func GetAllOrders() []model.Order {
	orders := repository.GetAllOrders()
	if len(orders) > 20 {
		return orders[:20]
	}
	return orders
}

func CreateOrder(order model.Order) model.Order {
	return repository.CreateOrder(order)
}

func GetOrdersByBotId(botId int64) []model.Order {
	return repository.GetAllOrdersByBotID(botId)
}
func UpdateOrder(order model.Order) model.Order {
	return repository.UpdateOrder(order)
}
