package botFather

import (
	"fmt"
	"github.com/berserkkv/trader/bot"
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/model/enum/order"
	"github.com/berserkkv/trader/model/enum/symbol"
	"github.com/berserkkv/trader/model/enum/timeframe"
	repository "github.com/berserkkv/trader/repository/interface"
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
	botRepo           repository.BotRepository
	orderRepo         repository.OrderRepository
	botsData          map[timeframe.Timeframe]map[symbol.Symbol]struct{}
	connector         connector.Connector
	smallestTimeFrame int
	mu                sync.Mutex
}

func (bf *BotFather) Start() {
	for {
		tools.WaitUntilNextAlignedTick(time.Duration(bf.smallestTimeFrame) * time.Second)

		runTime := time.Now()
		minute := runTime.Minute()
		hour := runTime.Hour()

		bf.runBots(minute, hour)
		//bf.botRepo.UpdateAll(bf.Bots())
	}
}

func (bf *BotFather) getCandles(timeFrame timeframe.Timeframe, candles map[string][]model.Candle) {
	for tf, smbSet := range bf.botsData {
		for smb := range smbSet {
			if tf == timeFrame {
				if _, ok := candles[string(tf)+string(smb)]; !ok {
					candles[string(tf)+string(smb)] = []model.Candle{}
					candles[string(tf)+string(smb)] = bf.connector.GetCandles(smb, tf, 202)
				}
			}
		}
	}
}

func (bf *BotFather) runBots(minute int, hour int) {
	time.Sleep(time.Duration(3) * time.Second)
	candles := map[string][]model.Candle{}

	if bf.botsData[timeframe.MINUTE_1] != nil {
		bf.getCandles(timeframe.MINUTE_1, candles)
	}

	if minute%5 == 0 {
		bf.getCandles(timeframe.MINUTE_5, candles)
	}
	if minute%15 == 0 {
		bf.getCandles(timeframe.MINUTE_15, candles)
	}
	if minute%30 == 0 {
		bf.getCandles(timeframe.MINUTE_30, candles)
	}
	if minute == 0 {
		bf.getCandles(timeframe.HOUR_1, candles)
	}
	if minute == 0 && hour == 0 {
		bf.getCandles(timeframe.DAY, candles)
	}
	for _, b := range bf.Bots() {

		if b == nil || b.IsNotActive || b.InPos {
			slog.Info("Bot skipped", "bot", b)
			continue
		}

		if b.CurrentCapital <= 85 {
			b.IsNotActive = true
			continue
		}

		switch b.Timeframe {
		case timeframe.MINUTE_1:
			bf.runStrategy(b, candles)

		case timeframe.MINUTE_5:
			if minute%5 == 0 {
				bf.runStrategy(b, candles)
			}
		case timeframe.MINUTE_15:
			if minute%15 == 0 {
				bf.runStrategy(b, candles)
			}

		case timeframe.MINUTE_30:
			if minute%30 == 0 {
				bf.runStrategy(b, candles)
			}

		case timeframe.HOUR_1:
			if minute == 0 {
				bf.runStrategy(b, candles)
			}

		case timeframe.DAY:
			if hour == 0 && minute == 0 {
				bf.runStrategy(b, candles)
			}

		default:

		}
	}
}

func (bf *BotFather) runStrategy(b *bot.Bot, candlesMap map[string][]model.Candle) {
	if b.Connector == nil {
		slog.Error("No connector found")
		return
	}
	candles := candlesMap[string(b.Timeframe)+string(b.Symbol)]

	if len(candles) == 0 {
		slog.Error("Not enough candles")
		return
	}

	cmd, info := b.Strategy.Run(candles)

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
	b.StrategyInfo = info
	b.LastScanned = time.Now()

}

func GetBotFather(botRepo repository.BotRepository, orderRepo repository.OrderRepository) *BotFather {
	once.Do(func() {
		instance = &BotFather{
			bots:      make(map[int64]*bot.Bot),
			botRepo:   botRepo,
			orderRepo: orderRepo,
			botsData:  make(map[timeframe.Timeframe]map[symbol.Symbol]struct{}),
			connector: connector.BinanceConnector{},
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
		return
	}

	bf.bots[bot.Id] = bot
	if _, ok := bf.botsData[bot.Timeframe]; !ok {
		bf.botsData[bot.Timeframe] = make(map[symbol.Symbol]struct{})
	}
	var tf int
	switch bot.Timeframe {
	case timeframe.MINUTE_1:
		tf = 60
	case timeframe.MINUTE_5:
		tf = 300
	case timeframe.MINUTE_15:
		tf = 900
	case timeframe.MINUTE_30:
		tf = 1800
	case timeframe.HOUR_1:
		tf = 3600
	}
	bf.botsData[bot.Timeframe][bot.Symbol] = struct{}{}
	if bf.smallestTimeFrame == 0 || (bf.smallestTimeFrame != 0 && bf.smallestTimeFrame > tf) {
		bf.smallestTimeFrame = tf
	}
	slog.Info("bot added successfully to BotFather", "name", bot.Name)
}

func (bf *BotFather) GetBot(id int64) (*bot.Bot, error) {
	b, ok := bf.bots[id]
	if !ok {
		return nil, fmt.Errorf("bot with id %d not found", id)
	}
	return b, nil
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

func (bf *BotFather) DeleteBot(id int64) {
	bf.mu.Lock()
	defer bf.mu.Unlock()

	bot, ok := bf.bots[id]
	if !ok {
		return // bot not found, nothing to delete
	}

	if bot.InPos {
		_, _ = bot.ClosePosition(100)
		bf.DecreaseTotalBotsInOrder()
	}

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
