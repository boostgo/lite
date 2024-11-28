package to

import (
	"time"
)

// Time convert value to time.Time object by provided format.
// If converting return error returns "zero time".
// "zero time" is empty time.Time object.
func Time(value any, format string) time.Time {
	parsedTime, err := time.Parse(format, String(value))
	if err != nil {
		return zeroTime()
	}

	return parsedTime
}

func zeroTime() time.Time {
	return time.Time{}
}
