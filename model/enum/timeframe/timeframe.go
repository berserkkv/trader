package timeframe

import "fmt"

type Timeframe string

const (
	MINUTE_1  Timeframe = "1m"
	MINUTE_5  Timeframe = "5m"
	MINUTE_15 Timeframe = "15m"
	MINUTE_30 Timeframe = "30m"
	HOUR_1    Timeframe = "1h"
	HOUR_4    Timeframe = "4h"
	DAY       Timeframe = "1D"
)

func Parse(s string) (Timeframe, error) {
	switch s {
	case string(MINUTE_1), string(MINUTE_5), string(MINUTE_15), string(MINUTE_30), string(HOUR_1), string(HOUR_4), string(DAY):
		return Timeframe(s), nil
	default:
		return "", fmt.Errorf("invalid timeframe: %s", s)
	}
}
