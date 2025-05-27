package botFather

import (
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/repository"
	"log/slog"
	"time"
)

func (bf *BotFather) CheckAndStartMonitoring() {
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

func (bf *BotFather) monitorPosition() {
	var closedOrders []model.Order

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

			curPrice := b.Connector.GetPrice(b.Symbol)

			if b.ShouldClosePosition(curPrice) {
				closedOrder, err := b.ClosePosition(curPrice)
				if err != nil {
					slog.Error("Can't close position", "error", err.Error(), "botName", b.Name)
					continue
				}

				bf.DecreaseTotalBotsInOrder()

				closedOrders = append(closedOrders, closedOrder)
				slog.Debug("Position closed, bot ready for new orders", "botName", b.Name)
			} else {
				b.UpdatePnlAndRoe(curPrice)
				b.GridOrderMonitor(curPrice)
				b.ShiftStopLoss()
				b.OrderScannedTime = time.Now()
				_, err := repository.UpdateBot(b)
				if err != nil {
					slog.Error("Error updating bot", "error", err.Error(), "botName", b.Name)
				}
			}

		}

		for i := range closedOrders {
			_, err := repository.UpdateBot(bf.bots[closedOrders[i].BotID])
			if err != nil {
				slog.Error("Can't update database after closing order", "error", err.Error(), "botName", bf.bots[closedOrders[i].BotID].Name)
			}
			repository.CreateOrder(closedOrders[i])
		}
		closedOrders = closedOrders[:0]
	}
}
