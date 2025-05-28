package service

import (
	"github.com/berserkkv/trader/database"
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/model/enum/symbol"
	"github.com/berserkkv/trader/model/enum/timeframe"
	"github.com/berserkkv/trader/repository"
	"gorm.io/gorm"
	"log/slog"
	"time"
)

var db *gorm.DB = database.DB
var r = repository.NewStatisticsRepository(db)

type Result struct {
	Sol15 []model.Statistics
	Sol1  []model.Statistics
	Eth15 []model.Statistics
	Eth1  []model.Statistics
}

func GetAllStatistics() (Result, error) {
	sol15 := []model.Statistics{}
	sol1 := []model.Statistics{}
	eth15 := []model.Statistics{}
	eth1 := []model.Statistics{}

	sts, err := r.GetAllStatistics()
	if err != nil {
		slog.Error("GetAllStatistics err:", err)
	}
	for _, st := range sts {
		if st.Symbol == symbol.SOLUSDT {
			if st.Timeframe == timeframe.MINUTE_15 {
				sol15 = append(sol15, st)
			} else if st.Timeframe == timeframe.MINUTE_1 {
				sol1 = append(sol1, st)
			}
		} else if st.Symbol == symbol.ETHUSDT {
			if st.Timeframe == timeframe.MINUTE_15 {
				eth15 = append(eth15, st)
			} else if st.Timeframe == timeframe.MINUTE_1 {
				eth1 = append(eth1, st)
			}
		}
	}
	res := Result{
		Sol15: sol15,
		Sol1:  sol1,
		Eth15: eth15,
		Eth1:  eth1,
	}
	return res, nil
}

func GetStatisticsBySymbol(symbol string) ([]model.Statistics, error) {
	return r.GetStatisticsBySymbol(symbol)
}

func GetStatisticsByTimeframe(timeframe string) ([]model.Statistics, error) {
	return r.GetStatisticsByTimeframe(timeframe)
}

func GetStatisticsByStrategy(strategyName string) ([]model.Statistics, error) {
	return r.GetStatisticsByStrategy(strategyName)
}

func SaveStatistics() {
	bots := repository.GetAllBots(nil)
	sts := make([]model.Statistics, 0)
	for _, b := range bots {
		newSt := model.Statistics{
			StrategyName: b.StrategyName,
			Symbol:       b.Symbol,
			Timeframe:    b.Timeframe,
			Balance:      b.CurrentCapital + b.OrderCapital,
			Time:         time.Now(),
		}
		sts = append(sts, newSt)
	}

	err := r.SaveAllStatistics(&sts)
	if err != nil {
		slog.Error("save statistics error", err)
	}
}
