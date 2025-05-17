package timeframe

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
