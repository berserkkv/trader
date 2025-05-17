package tools

import (
	"log/slog"
	"time"
)

func WaitUntilNextAlignedTick(interval time.Duration) {
	now := time.Now()

	// Find how far we are into the interval
	elapsed := time.Duration(now.Minute())*time.Minute +
		time.Duration(now.Second())*time.Second +
		time.Duration(now.Nanosecond())

	// Compute time until next interval
	wait := interval - (elapsed % interval)
	if wait <= 0 {
		wait = interval
	}

	slog.Debug("Waiting until next aligned tick.\n", "wait", wait)

	time.Sleep(wait)
}
