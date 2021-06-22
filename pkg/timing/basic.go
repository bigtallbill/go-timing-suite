package timing

import (
	"fmt"
	"log"
	"time"
)

// LogDuration is a convenience method for use in defer statements.
// Usage is `defer timing.LogDuration("Op took", time.Now())`.
// If one or many log.Logger instances are provided, then they will be used
// instead of the default system logger.
func LogDuration(message string, start time.Time, loggers ...*log.Logger) {
	var msg = fmt.Sprintf("%s: %s", message, time.Since(start))

	if len(loggers) == 0 {
		log.Print(msg)
	} else {
		for _, logger := range loggers {
			logger.Print(msg)
		}
	}
}
