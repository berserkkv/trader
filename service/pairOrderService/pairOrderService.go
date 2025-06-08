package pairOrderService

import (
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/repository/pairBotRepository"
	"github.com/berserkkv/trader/repository/pairOrderRepository"
)

func GetAllOrders() []model.PairOrder {
	orders := pairOrderRepository.GetAllOrders("Desc")
	if len(orders) > 20 {
		return orders[:20]
	}
	return orders
}

func GetOrderStatistics(botId int64) []model.Statistics {
	o := pairOrderRepository.GetAllOrdersByBotID(botId)
	if len(o) == 0 {
		return []model.Statistics{}
	}
	s := make([]model.Statistics, 0)
	s = append(s, model.Statistics{
		Pnl:  o[0].ProfitLoss1 - o[0].Fee1 + o[0].ProfitLoss2 - o[0].Fee2,
		Time: o[0].ClosedTime,
	})
	for i := 1; i < len(o); i++ {
		statistic := model.Statistics{
			Pnl:  o[i].ProfitLoss1 + s[i-1].Pnl - o[i].Fee1 + o[1].ProfitLoss2 - o[i].Fee2,
			Time: o[i].ClosedTime,
		}
		s = append(s, statistic)
	}
	return s
}

func GetAllOrderStatistics() map[string][]model.Statistics {
	orders := pairOrderRepository.GetAllOrders("asc")
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
			Pnl:  o.ProfitLoss1 + lastPnl - o.Fee1 + o.ProfitLoss2 - o.Fee2,
			Time: o.ClosedTime,
		})
	}
	res := make(map[string][]model.Statistics)

	bots := pairBotRepository.GetAllBots(nil)

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

func CreateOrder(order model.PairOrder) model.PairOrder {
	return pairOrderRepository.CreateOrder(order)
}

func GetOrdersByBotId(botId int64) []model.PairOrder {
	return pairOrderRepository.GetAllOrdersByBotIDDesc(botId)
}
func UpdateOrder(order model.PairOrder) model.PairOrder {
	return pairOrderRepository.UpdateOrder(order)
}
