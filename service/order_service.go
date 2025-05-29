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

func GetOrderStatistics(botId int64) []model.Statistics {
	o := repository.GetAllOrdersByBotID(botId)

	s := make([]model.Statistics, 0)
	for _, order := range o {
		statistic := model.Statistics{
			Pnl:  order.ProfitLossPercent,
			Time: order.ClosedTime,
		}
		s = append(s, statistic)
	}
	return s
}

func CreateOrder(order model.Order) model.Order {
	return repository.CreateOrder(order)
}

func GetOrdersByBotId(botId int64) []model.Order {
	return repository.GetAllOrdersByBotIDDesc(botId)
}
func UpdateOrder(order model.Order) model.Order {
	return repository.UpdateOrder(order)
}
