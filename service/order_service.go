package service

import (
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/repository"
)

func GetAllOrders() []model.Order {
	orders := repository.GetAllOrders("Desc")
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

func GetAllOrderStatistics() map[string][]model.Statistics {
	orders := repository.GetAllOrders("asc")
	if len(orders) == 0 {
		return map[string][]model.Statistics{}
	}

	statsMap := make(map[int64][]model.Statistics)

	for _, o := range orders {
		lastPnl := 0.0
		if _, exists := statsMap[o.BotID]; exists {
			lastPnl = statsMap[o.BotID][len(statsMap[o.BotID])-1].Pnl
		}
		statsMap[o.BotID] = append(statsMap[o.BotID], model.Statistics{
			Pnl:  o.ProfitLoss + lastPnl,
			Time: o.ClosedTime,
		})
	}
	res := make(map[string][]model.Statistics)

	bots := repository.GetAllBots(nil)

	for _, b := range bots {
		if b.IsNotActive {
			continue
		}
		if _, exists := statsMap[b.Id]; exists {
			res[b.Name] = statsMap[b.Id]
		}
	}

	return res
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
