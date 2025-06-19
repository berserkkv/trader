package serviceImpl

import (
	"github.com/berserkkv/trader/model"
	repository "github.com/berserkkv/trader/repository/interface"
	service "github.com/berserkkv/trader/service/interface"
)

type OrderServiceImpl struct {
	r repository.OrderRepository
}

func NewOrderServiceImpl(r repository.OrderRepository) *OrderServiceImpl {
	return &OrderServiceImpl{r: r}
}

func (s *OrderServiceImpl) GetAll(fields map[string]interface{}) []*model.Order {
	return s.r.GetAll(fields)
}

func (s *OrderServiceImpl) GetByID(id int64) *model.Order {
	return s.r.GetByID(id)
}

func (s *OrderServiceImpl) GetAllByBotID(botID int64) []*model.Order {
	return s.r.GetAllByBotID(botID, "DESC")
}

func (s *OrderServiceImpl) GetStatisticsByBotID(botID int64) []*model.Statistics {
	o := s.r.GetAllByBotID(botID, "ASC")
	if len(o) == 0 {
		return []*model.Statistics{}
	}
	statistics := make([]*model.Statistics, 0)

	statistics = append(statistics, &model.Statistics{
		Pnl:  o[0].ProfitLoss - o[0].Fee,
		Time: o[0].ClosedTime,
	})
	for i := 1; i < len(o); i++ {
		statistic := model.Statistics{
			Pnl:  o[i].ProfitLoss + statistics[i-1].Pnl - o[i].Fee,
			Time: o[i].ClosedTime,
		}
		statistics = append(statistics, &statistic)
	}
	return statistics
}

func (s *OrderServiceImpl) GetAllStatistics() map[string][]*model.Statistics {
	//orders := s.r.GetAll(nil)
	//if len(orders) == 0 {
	//	return map[string][]*model.Statistics{}
	//}
	//
	//statsMap := make(map[int64][]model.Statistics)
	//
	//for _, o := range orders {
	//	lastPnl := 0.0
	//	if _, exists := statsMap[o.BotID]; exists {
	//		lastPnl = statsMap[o.BotID][len(statsMap[o.BotID])-1].Pnl
	//	}
	//	statsMap[o.BotID] = append(statsMap[o.BotID], model.Statistics{
	//		Pnl:  o.ProfitLoss + lastPnl - o.Fee,
	//		Time: o.ClosedTime,
	//	})
	//}
	//res := make(map[string][]model.Statistics)
	//
	//bots := repository.GetAllBots(nil)
	//
	//for _, b := range bots {
	//	if b.IsNotActive {
	//		continue
	//	}
	//	if _, exists := statsMap[b.Id]; exists {
	//		res[b.Name] = statsMap[b.Id]
	//	}
	//}
	//
	//return res

	return nil
}

func (s *OrderServiceImpl) Create(order *model.Order) (*model.Order, error) {
	return s.r.Create(order)
}

func (s *OrderServiceImpl) Update(order *model.Order) (*model.Order, error) {
	return s.r.Update(order)
}

func (s *OrderServiceImpl) Delete(id int64) {
	s.r.Delete(id)
}

func (s *OrderServiceImpl) DeleteAllByBotID(botID int64) {
	s.r.DeleteAllByBotID(botID)
}

var _ service.OrderService = (*OrderServiceImpl)(nil)
