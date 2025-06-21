package backtest

import (
	"github.com/berserkkv/trader/bot"
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/model/enum/order"
	"github.com/berserkkv/trader/service/tools"
	"log/slog"
	"math/rand"
	"time"
)

type BacktestBotFather struct {
	candles []model.Candle
	length  int
	Index   int
}

func NewBacktestBotFather(length int) *BacktestBotFather {
	return &BacktestBotFather{
		length: length,
	}
}

func (b *BacktestBotFather) Run(bot *bot.Bot) {
	for b.next() {
		c := b.getCandles()
		if len(c) == 0 {
			break
		}
		if bot.CurrentCapital > bot.InitialCapital*10 {
			slog.Info("Excellent")
			return
		}

		if bot.InPos {
			if b.stopLossInRange(bot.OrderStopLoss, c[len(c)-1], bot.OrderType) {
				_, err := bot.ClosePosition(bot.OrderStopLoss)
				if err != nil {
					slog.Error("Error closing position", "err", err)
				}
			}

			if b.takeProfitInRange(bot.OrderTakeProfit, c[len(c)-1], bot.OrderType) {
				_, err := bot.ClosePosition(bot.OrderTakeProfit)
				if err != nil {
					slog.Error("Error closing position", "err", err)
				}
			}
		} else {
			if bot.CurrentCapital <= 95 {
				slog.Info("Bot Stopped", "capital", bot.CurrentCapital)
				return
			}
			command, _ := bot.Strategy.Run(c)

			if command == order.LONG || command == order.SHORT {
				err := bot.OpenPosition(command)
				if err != nil {
					slog.Error("Error opening position", "error", err)
					return
				}
			}
		}

	}
}

func (b *BacktestBotFather) LoadData(name string) error {
	candles, err := tools.LoadCandlesFromCSV(name)
	if err != nil {
		slog.Error("Error loading candles from CSV", "error", err)
		return err
	}
	b.candles = candles
	return nil
}

func (b *BacktestBotFather) getCandles() []model.Candle {
	end := b.Index + b.length
	if end > len(b.candles) {
		slog.Info("End of candles")
		return []model.Candle{}
	}
	return b.candles[b.Index:end]
}

func (b *BacktestBotFather) next() bool {
	if b.Index+1 >= len(b.candles) {
		return false
	}
	b.Index++
	return true
}

func (b *BacktestBotFather) stopLossInRange(stopLoss float64, c model.Candle, command order.Command) bool {
	if command == order.LONG {
		return c.Close <= stopLoss || c.Open <= stopLoss || c.High <= stopLoss || c.Low <= stopLoss
	} else if command == order.SHORT {
		return c.Close >= stopLoss || c.Open >= stopLoss || c.High >= stopLoss || c.Low >= stopLoss
	}
	return false
}

func (b *BacktestBotFather) takeProfitInRange(p float64, candle model.Candle, command order.Command) bool {
	if command == order.LONG {
		return candle.High >= p || candle.Open >= p || candle.Close >= p || candle.Low >= p
	} else if command == order.SHORT {
		return candle.Low <= p || candle.Open <= p || candle.Close <= p || candle.High <= p
	}
	return false
}

func (b *BacktestBotFather) Randomize() {
	rand.Seed(time.Now().UnixNano())

	b.Index = rand.Intn(len(b.candles))
	slog.Info("Randomized", "Index", b.Index)
}
