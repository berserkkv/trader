package strategy

import (
	"fmt"
	"github.com/berserkkv/trader/model"
	"github.com/berserkkv/trader/model/enum/order"
	"testing"
)

func TestHA_SLMAStrategy_LongSignal(t *testing.T) {
	candles := make([]model.Candle, 26)

	// Fill with neutral candles
	for i := 0; i < 20; i++ {
		candles[i] = model.Candle{
			Open:  100,
			High:  101,
			Low:   99,
			Close: 100,
		}
	}

	// Force 3 candles that produce red HA
	candles[20] = model.Candle{Open: 105, High: 106, Low: 101, Close: 102}
	candles[21] = model.Candle{Open: 102, High: 103, Low: 98, Close: 99}
	candles[22] = model.Candle{Open: 99, High: 100, Low: 95, Close: 96}

	// Make the 4th candle push HA into red territory
	candles[23] = model.Candle{Open: 96, High: 97, Low: 92, Close: 97} // still red

	// Now flip it: strong bullish candle to trigger green HA
	candles[24] = model.Candle{Open: 93, High: 105, Low: 92, Close: 104}
	candles[25] = model.Candle{Open: 93, High: 105, Low: 92, Close: 104}

	// Create strategy
	strat := &HA_SLMA{}

	// First call: should detect HA color change
	cmd, info := strat.Start(candles)
	fmt.Println("CALL 1:", cmd, info)

	if cmd != order.LONG {
		t.Errorf("Expected LONG command, got %v", cmd)
	}
}
