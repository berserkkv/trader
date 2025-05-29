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
	if len(o) == 0 {
		return []model.Statistics{}
	}
	s := make([]model.Statistics, 0)
	s = append(s, model.Statistics{
		Pnl:  o[0].ProfitLossPercent,
		Time: o[0].ClosedTime,
	})
	for i := 1; i < len(o); i++ {
		statistic := model.Statistics{
			Pnl:  o[i].ProfitLossPercent + s[i-1].Pnl,
			Time: o[i].ClosedTime,
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
