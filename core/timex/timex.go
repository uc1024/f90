package timex

import (
	"errors"
	"time"
)

const (
	TimeLayout = "2006-01-02 15:04:05"
)

// * Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
func ParseDuration(s string) time.Duration {
	if in, err := time.ParseDuration(s); err == nil {
		return in
	} else {
		return time.Duration(0)
	}
}

// * u eq second
func Uint64TimeFormat(u uint64) string {
	return time.Unix(int64(u), 0).Format(TimeLayout)
}

// IsCurrentTimeWithinInterval checks if the current time falls within the interval defined by the given Unix timestamp and offset.
func IsCurrentTimeWithinInterval(unixTimestamp int64, offset time.Duration) error {
	// Calculate the start and end times of the interval in milliseconds
	intervalStart := time.UnixMilli(unixTimestamp).UnixMilli()
	intervalEnd := time.UnixMilli(unixTimestamp).Add(offset).UnixMilli()
	currentTime := time.Now().UnixMilli()

	// Check if the current time falls outside the interval
	if !(currentTime <= intervalEnd && currentTime >= intervalStart) {
		return errors.New("time out")
	}
	return nil
}
