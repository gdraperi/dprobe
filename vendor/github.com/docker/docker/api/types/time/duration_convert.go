package time

import (
	"strconv"
	"time"
)

// DurationToSecondsString converts the specified duration to the number
// seconds it represents, formatted as a string.
func DurationToSecondsString(duration time.Duration) string ***REMOVED***
	return strconv.FormatFloat(duration.Seconds(), 'f', 0, 64)
***REMOVED***
