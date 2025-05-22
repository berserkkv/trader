package botFather

import (
	"fmt"
	"github.com/berserkkv/trader/bot"
	"github.com/berserkkv/trader/model/enum/order"
	"github.com/berserkkv/trader/model/enum/timeframe"
	"github.com/berserkkv/trader/repository"

	"github.com/berserkkv/trader/service/connector"
	"github.com/berserkkv/trader/service/tools"
	"log/slog"
	"sync"
	"time"
)

var (
	instance *BotFather
	once     sync.Once
)

type BotFather struct {
	bots              map[int64]*bot.Bot
	totalBotsInOrder  int64
	monitoringRunning bool
	mu                sync.Mutex
}

func (bf *BotFather) Start() {
	for {
		tools.WaitUntilNextAlignedTick(60 * time.Second)

		runTime := time.Now()
		minute := runTime.Minute()
		hour := runTime.Hour()

		bf.runBots(minute, hour)
		repository.UpdateAllBots(bf.Bots())
	}
}

func (bf *BotFather) runBots(minute int, hour int) {
	for _, b := range bf.Bots() {

		if b == nil || b.IsNotActive || b.InPos {
			slog.Info("Bot skipped", "bot", b)
			continue
		}

		switch b.TimeFrame {
		case timeframe.MINUTE_1:
			bf.runStrategy(b)

		case timeframe.MINUTE_5:
			if minute%5 == 0 {
				bf.runStrategy(b)
			}
		case timeframe.MINUTE_15:
			if minute%15 == 0 {
				bf.runStrategy(b)
			}

		case timeframe.MINUTE_30:
			if minute%30 == 0 {
				bf.runStrategy(b)
			}

		case timeframe.HOUR_1:
			if minute == 0 {
				bf.runStrategy(b)
			}

		case timeframe.DAY:
			if hour == 0 && minute == 0 {
				bf.runStrategy(b)
			}

		default:

		}
	}
}

func (bf *BotFather) runStrategy(b *bot.Bot) {
	candles := connector.GetKlines(b.Symbol, b.TimeFrame, 50)

	slog.Debug("Fetched klines from API", "length", len(candles))

	cmd := b.Strategy.Start(candles)

	slog.Info("Scanned", "command", cmd, "bot", b)

	switch cmd {
	case order.LONG, order.SHORT:
		err := b.OpenPosition(cmd)
		if err != nil {
			slog.Error("Error opening position", "error", err)
			return
		}

		bf.IncreaseTotalBotsInOrder()

		bf.CheckAndStartMonitoring()

	case order.WAIT:
		slog.Debug("No signal yet", "name", b.Name)
	default:
		slog.Debug("Order command not identified", "name", b.Name, "command", cmd)
	}

	b.LastScanned = time.Now()

}

func GetBotFather() *BotFather {
	once.Do(func() {
		instance = &BotFather{
			bots: make(map[int64]*bot.Bot),
		}
	})
	return instance
}

func (bf *BotFather) AddBot(bot *bot.Bot) {
	if bot == nil {
		slog.Error("bot not added to BotFather, bot is nil")
		return
	}
	if bot.Id == 0 {
		slog.Error("bot not added to BotFather, bot id is 0")
		return
	}
	if _, exists := bf.bots[bot.Id]; exists {
		slog.Error("bot not added to BotFather, bot with id already exists", "botId", bot.Id)
		return
	}
	if bot.Strategy == nil {
		slog.Error("bot not added to BotFather, bot strategy is nil")
	}

	bf.bots[bot.Id] = bot
	slog.Info("bot added successfully to BotFather", "name", bot.Name)
}

func (bf *BotFather) IncreaseTotalBotsInOrder() {
	bf.mu.Lock()
	bf.totalBotsInOrder += 1
	bf.mu.Unlock()
}

func (bf *BotFather) DecreaseTotalBotsInOrder() {
	bf.mu.Lock()
	bf.totalBotsInOrder -= 1
	bf.mu.Unlock()
}

func (bf *BotFather) RemoveBot(id int64) {
	delete(bf.bots, id)
}

func (bf *BotFather) Bots() []*bot.Bot {
	return mapToSlice(bf.bots)
}

func mapToSlice(m map[int64]*bot.Bot) []*bot.Bot {
	bots := make([]*bot.Bot, 0, len(m))
	for _, b := range m {
		bots = append(bots, b)
	}
	return bots
}

func (bf *BotFather) StopBot(id int64) (*bot.Bot, error) {
	b, ok := bf.bots[id]
	if !ok {
		return nil, fmt.Errorf("bot with id %d not found", id)
	}
	if b.InPos {
		return nil, fmt.Errorf("bot with id %d is in position", id)
	}
	b.IsNotActive = true
	return b, nil
}

func (bf *BotFather) StartBot(id int64) (*bot.Bot, error) {
	b, ok := bf.bots[id]
	if !ok {
		return nil, fmt.Errorf("bot with id %d not found", id)
	}
	if b.CurrentCapital <= 10 {
		return nil, fmt.Errorf("bot with id %d has capital less than 10", id)
	}

	b.IsNotActive = false
	return b, nil
}
