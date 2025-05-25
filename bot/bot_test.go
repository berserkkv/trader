package bot_test

import (
	"github.com/berserkkv/trader/bot"
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/model/enum/order"
	"github.com/berserkkv/trader/model/enum/symbol"
	"github.com/berserkkv/trader/model/enum/timeframe"
	"github.com/berserkkv/trader/service/calculator"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type MockStrategy struct{}

func (m *MockStrategy) Name() string {
	return "MOCK"
}

func (m *MockStrategy) Start(candles []model.Candle) order.Command {
	return order.LONG
}

func TestNewBot(t *testing.T) {
	st := &MockStrategy{}
	b := bot.NewBot(timeframe.MINUTE_1, st, symbol.BTCUSDT, 1000.0, 10.0, 0.5, 0.3)

	assert.Equal(t, "MOCK_1m_BTCUSDT", b.Name)
	assert.Equal(t, 1000.0, b.InitialCapital)
	assert.Equal(t, b.CurrentCapital, b.InitialCapital)
	assert.Equal(t, 10.0, b.Leverage)
	assert.Equal(t, 0.5, b.TakeProfit)
	assert.Equal(t, 0.3, b.StopLoss)
	assert.False(t, b.InPos)
}

func TestCanOpenPosition(t *testing.T) {
	st := &MockStrategy{}
	b := bot.NewBot(timeframe.MINUTE_1, st, symbol.BTCUSDT, 1000.0, 10.0, 0.5, 0.3)

	// Bot is active and not in position
	err := b.CanOpenPosition()
	assert.Nil(t, err)

	// Bot is already in position
	b.InPos = true
	err = b.CanOpenPosition()
	assert.Error(t, err)
	b.InPos = false

	// Bot is not active
	b.IsNotActive = true
	err = b.CanOpenPosition()
	assert.Error(t, err)
	b.IsNotActive = false

	// Insufficient capital
	b.CurrentCapital = 5
	err = b.CanOpenPosition()
	assert.Error(t, err)
}

func TestShiftStopLoss(t *testing.T) {
	tests := []struct {
		name             string
		bot              bot.Bot
		expectedStopLoss float64
	}{
		{
			name: "ROE < 0.2 → no shift",
			bot: bot.Bot{
				Roe:             1.0,
				Leverage:        10.0,
				OrderEntryPrice: 100.0,
				OrderStopLoss:   95.0,
				OrderType:       order.LONG,
			},
			expectedStopLoss: 95.0,
		},
		{
			name: "LONG → shift stop loss up",
			bot: bot.Bot{
				Roe:             4.0,
				Leverage:        10.0,
				OrderEntryPrice: 100.0,
				OrderStopLoss:   95.0,
				OrderType:       order.LONG,
			},
			expectedStopLoss: 100.2, // 100 * (1 + (4/100)/2)
		},
		{
			name: "SHORT → shift stop loss down",
			bot: bot.Bot{
				Roe:             4.0,
				Leverage:        10.0,
				OrderEntryPrice: 100.0,
				OrderStopLoss:   105.0,
				OrderType:       order.SHORT,
			},
			expectedStopLoss: 99.8,
		},
		{
			name: "LONG → new stop loss < current → no shift",
			bot: bot.Bot{
				Roe:             4.0,
				Leverage:        10.0,
				OrderEntryPrice: 100.0,
				OrderStopLoss:   101.0,
				OrderType:       order.LONG,
			},
			expectedStopLoss: 101.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.bot.ShiftStopLoss()
			if tt.bot.OrderStopLoss != tt.expectedStopLoss {
				t.Errorf("expected %f, got %f", tt.expectedStopLoss, tt.bot.OrderStopLoss)
			}
		})
	}
}

func TestShouldClosePosition(t *testing.T) {
	tests := []struct {
		name       string
		bot        bot.Bot
		curPrice   float64
		shouldExit bool
	}{
		// LONG cases
		{
			name: "LONG - hit take profit",
			bot: bot.Bot{
				OrderType:       order.LONG,
				OrderTakeProfit: 110.0,
				OrderStopLoss:   90.0,
			},
			curPrice:   110.0,
			shouldExit: true,
		},
		{
			name: "LONG - hit stop loss",
			bot: bot.Bot{
				OrderType:       order.LONG,
				OrderTakeProfit: 110.0,
				OrderStopLoss:   90.0,
			},
			curPrice:   89.5,
			shouldExit: true,
		},
		{
			name: "LONG - still in range",
			bot: bot.Bot{
				OrderType:       order.LONG,
				OrderTakeProfit: 110.0,
				OrderStopLoss:   90.0,
			},
			curPrice:   100.0,
			shouldExit: false,
		},

		// SHORT cases
		{
			name: "SHORT - hit take profit",
			bot: bot.Bot{
				OrderType:       order.SHORT,
				OrderTakeProfit: 90.0,
				OrderStopLoss:   110.0,
			},
			curPrice:   90.0,
			shouldExit: true,
		},
		{
			name: "SHORT - hit stop loss",
			bot: bot.Bot{
				OrderType:       order.SHORT,
				OrderTakeProfit: 90.0,
				OrderStopLoss:   110.0,
			},
			curPrice:   111.0,
			shouldExit: true,
		},
		{
			name: "SHORT - still in range",
			bot: bot.Bot{
				OrderType:       order.SHORT,
				OrderTakeProfit: 90.0,
				OrderStopLoss:   110.0,
			},
			curPrice:   100.0,
			shouldExit: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.bot.ShouldClosePosition(tt.curPrice)
			if got != tt.shouldExit {
				t.Errorf("expected %v, got %v", tt.shouldExit, got)
			}
		})
	}
}

func TestUpdatePnlAndRoe(t *testing.T) {
	b := &bot.Bot{
		OrderEntryPrice:          100.0,
		Leverage:                 10,
		OrderType:                order.LONG,
		OrderCapitalWithLeverage: 1000.0,
		OrderQuantity:            10.0,
	}
	curPrice := 101.0

	// Expected values using same logic as your calculator (simple test)
	expectedRoe := calculator.CalculateRoe(b.OrderEntryPrice, curPrice, b.Leverage, b.OrderType)
	expectedPnl := calculator.CalculatePNL(curPrice, b.OrderCapitalWithLeverage, b.OrderQuantity, b.OrderType)

	b.UpdatePnlAndRoe(curPrice)
	if b.Roe != expectedRoe {
		t.Errorf("expected ROE %.2f, got %.2f", expectedRoe, b.Roe)
	}
	if b.Pnl != expectedPnl {
		t.Errorf("expected PNL %.2f, got %.2f", expectedPnl, b.Pnl)
	}
}

type mockConnector struct{}

func (mockConnector) GetPrice(smb symbol.Symbol) float64 {
	return 10.0 // mock price
}

func (mockConnector) GetCandles(smb symbol.Symbol, interval timeframe.Timeframe, limit int) []model.Candle {
	return []model.Candle{}
}

func TestOpenPosition(t *testing.T) {
	b := &bot.Bot{
		Name:           "TestBot",
		Symbol:         "BTCUSDT",
		Connector:      &mockConnector{},
		CurrentCapital: 100,
		Leverage:       10,
		OrderType:      order.LONG,
		StopLoss:       1,
		TakeProfit:     2,
	}

	err := b.OpenPosition(order.LONG)

	assert.NoError(t, err)
	assert.True(t, b.InPos)
	assert.Equal(t, 10.0, b.OrderEntryPrice)
	assert.Equal(t, 99.98, b.OrderQuantity)
	assert.Equal(t, 9.9, b.OrderStopLoss)
	assert.Equal(t, 10.2, b.OrderTakeProfit)
	assert.InDelta(t, 999.8, b.OrderCapitalWithLeverage, 1e-6)
	assert.Equal(t, 0.02, b.OrderFee)
	assert.WithinDuration(t, time.Now(), b.OrderCreatedTime, time.Second)
	assert.Equal(t, order.LONG, b.OrderType)
}
