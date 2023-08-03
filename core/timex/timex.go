package timex

import (
	"errors"
	"time"

	"github.com/spf13/cast"
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

// * 是否在当前时间区间
func CurrentTimeOut(unix int64, offset time.Duration) error {
	// * 时间区间
	rt := time.Unix(cast.ToInt64(unix), 0).Unix()
	ex := time.Unix(cast.ToInt64(unix), 0).Add(offset).Unix()
	now := time.Now().Unix()
	if !(now <= ex && now >= rt) {
		return errors.New("time out")
	}
	return nil
}
