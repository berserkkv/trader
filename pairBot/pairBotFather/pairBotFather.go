package pairBotFather

import (
	"fmt"
	"github.com/berserkkv/trader/model/enum/order"
	"github.com/berserkkv/trader/model/enum/timeframe"
	"github.com/berserkkv/trader/pairBot"
	"github.com/berserkkv/trader/repository/pairBotRepository"

	"github.com/berserkkv/trader/service/tools"
	"log/slog"
	"sync"
	"time"
)

var (
	instance *PairBotFather
	once     sync.Once
)

type PairBotFather struct {
	bots              map[int64]*pairBot.PairBot
	totalBotsInOrder  int64
	monitoringRunning bool
	mu                sync.Mutex
}

func (bf *PairBotFather) Start() {
	for {
		tools.WaitUntilNextAlignedTick(60 * time.Second)

		runTime := time.Now()
		minute := runTime.Minute()
		hour := runTime.Hour()

		bf.runBots(minute, hour)
		pairBotRepository.UpdateAllBots(bf.Bots())
	}
}

func (bf *PairBotFather) runBots(minute int, hour int) {
	time.Sleep(time.Duration(5) * time.Second)
	for _, b := range bf.Bots() {

		if b == nil || b.IsNotActive || b.InPos {
			slog.Info("Bot skipped", "bot", b)
			continue
		}

		if b.CurrentCapital1 <= 85 || b.CurrentCapital2 <= 85 {
			b.IsNotActive = true
			continue
		}

		switch b.Timeframe {
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

func (bf *PairBotFather) runStrategy(b *pairBot.PairBot) {

	cmd := b.ShouldOpenPosition()

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

func GetPairBotFather() *PairBotFather {
	once.Do(func() {
		instance = &PairBotFather{
			bots: make(map[int64]*pairBot.PairBot),
		}
	})
	return instance
}

func (bf *PairBotFather) AddBot(bot *pairBot.PairBot) {
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

	bf.bots[bot.Id] = bot
	slog.Info("bot added successfully to BotFather", "name", bot.Name)
}

func (bf *PairBotFather) IncreaseTotalBotsInOrder() {
	bf.mu.Lock()
	bf.totalBotsInOrder += 1
	bf.mu.Unlock()
}

func (bf *PairBotFather) DecreaseTotalBotsInOrder() {
	bf.mu.Lock()
	bf.totalBotsInOrder -= 1
	bf.mu.Unlock()
}

func (bf *PairBotFather) RemoveBot(id int64) {
	delete(bf.bots, id)
}

func (bf *PairBotFather) Bots() []*pairBot.PairBot {
	return mapToSlice(bf.bots)
}

func mapToSlice(m map[int64]*pairBot.PairBot) []*pairBot.PairBot {
	bots := make([]*pairBot.PairBot, 0, len(m))
	for _, b := range m {
		bots = append(bots, b)
	}
	return bots
}

func (bf *PairBotFather) StopBot(id int64) (*pairBot.PairBot, error) {
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

func (bf *PairBotFather) StartBot(id int64) (*pairBot.PairBot, error) {
	b, ok := bf.bots[id]
	if !ok {
		return nil, fmt.Errorf("bot with id %d not found", id)
	}
	if b.CurrentCapital1 <= 80 {
		return nil, fmt.Errorf("bot with id %d has capital less than 80", id)
	}

	b.IsNotActive = false
	return b, nil
}
