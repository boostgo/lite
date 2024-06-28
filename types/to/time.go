package to

import (
	"time"
)

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
