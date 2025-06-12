package pairBotFather

import (
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/repository/pairBotRepository"
	"github.com/berserkkv/trader/repository/pairOrderRepository"
	"log/slog"
	"time"
)

func (bf *PairBotFather) CheckAndStartMonitoring() {
	bf.mu.Lock()
	defer bf.mu.Unlock()

	if bf.totalBotsInOrder > 0 && !bf.monitoringRunning {
		bf.monitoringRunning = true
		go func() {
			bf.monitorPosition()
			bf.mu.Lock()
			bf.monitoringRunning = false
			bf.mu.Unlock()
		}()
	}
}

func (bf *PairBotFather) monitorPosition() {
	var closedOrders []model.PairOrder

	for {
		time.Sleep(10 * time.Second)

		bf.mu.Lock()
		total := bf.totalBotsInOrder
		if total <= 0 {
			bf.monitoringRunning = false
			bf.totalBotsInOrder = 0
			bf.mu.Unlock()
			slog.Info("No bots in order, stopping MonitorPosition")
			return
		}
		bf.mu.Unlock()

		for _, b := range bf.bots {
			if !b.InPos {
				continue
			}

			curPrice1 := b.Connector.GetPrice(b.Symbol1)
			curPrice2 := b.Connector.GetPrice(b.Symbol2)

			b.UpdateZScore()
			b.UpdatePnlAndRoe(curPrice1, curPrice2)

			if b.ShouldClosePosition() || b.Roe1+b.Roe2 <= -10 {
				closedOrder, err := b.ClosePosition(curPrice1, curPrice2)
				if err != nil {
					slog.Error("Can't close position", "error", err.Error(), "botName", b.Name)
					continue
				}

				bf.DecreaseTotalBotsInOrder()

				closedOrders = append(closedOrders, closedOrder)
				slog.Debug("Position closed, bot ready for new orders", "botName", b.Name)
			} else {

				b.OrderScannedTime = time.Now()
				_, err := pairBotRepository.UpdateBot(b)
				if err != nil {
					slog.Error("Error updating bot", "error", err.Error(), "botName", b.Name)
				}
			}

		}

		for i := range closedOrders {
			_, err := pairBotRepository.UpdateBot(bf.bots[closedOrders[i].BotID])
			if err != nil {
				slog.Error("Can't update database after closing order", "error", err.Error(), "botName", bf.bots[closedOrders[i].BotID].Name)
			}
			pairOrderRepository.CreateOrder(closedOrders[i])
		}
		closedOrders = closedOrders[:0]
	}
}

func (bf *PairBotFather) ClosePosition(botId int64) {
	b := bf.bots[botId]
	curPrice1 := b.Connector.GetPrice(b.Symbol1)
	curPrice2 := b.Connector.GetPrice(b.Symbol2)
	closedOrder, err := b.ClosePosition(curPrice1, curPrice2)
	if err != nil {
		slog.Error("Error closing position", "error", err.Error(), "botName", b.Name)
	}

	bf.DecreaseTotalBotsInOrder()

	_, err = pairBotRepository.UpdateBot(b)
	if err != nil {
		slog.Error("Can't update database after closing order", "error", err.Error(), "botName", b.Name)
	}
	pairOrderRepository.CreateOrder(closedOrder)

}
