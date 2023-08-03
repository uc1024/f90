package throttlex

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	// Unknown means not initialized state.
	Unknown = iota
	// Allowed means allowed state.
	Allowed
	// HitQuota means this request exactly hit the quota.
	HitQuota
	// OverQuota means passed the quota.
	OverQuota

	internalOverQuota = 0
	internalAllowed   = 1
	internalHitQuota  = 2
)

var (
	// ErrUnknownCode is an error that represents unknown status code.
	ErrUnknownCode = errors.New("unknown status code")

	// to be compatible with aliyun redis, we cannot use `local key = KEYS[1]` to reuse the key
	periodScript = redis.NewScript(`local limit = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local current = redis.call("INCRBY", KEYS[1], 1)
if current == 1 then
    redis.call("expire", KEYS[1], window)
end
if current < limit then
    return 1
elseif current == limit then
    return 2
else
    return 0
end`)
)

type (
	// PeriodOption defines the method to customize a PeriodLimit.
	PeriodOption func(l *PeriodLimit)

	// A PeriodLimit is used to limit requests during a period of time.
	PeriodLimit struct {
		period     int // * 时间段长度，单位为秒
		quota      int // * 时间段内最多请求次数
		limitStore *redis.Client
		keyPrefix  string // * 存储 Redis key 的前缀
		align      bool   // * 是否对齐时间段开始时间
	}
)

// NewPeriodLimit 返回一个配置好的 PeriodLimit 实例
func NewPeriodLimit(period, quota int, limitStore *redis.Client, keyPrefix string,
	opts ...PeriodOption) *PeriodLimit {
	limiter := &PeriodLimit{
		period:     period,
		quota:      quota,
		limitStore: limitStore,
		keyPrefix:  keyPrefix,
	}

	for _, opt := range opts {
		opt(limiter)
	}

	return limiter
}

// Take requests a permit, it returns the permit state.
func (h *PeriodLimit) Take(key string) (int, error) {
	return h.TakeCtx(context.Background(), key)
}

// TakeCtx requests a permit with context, it returns the permit state.
func (h *PeriodLimit) TakeCtx(ctx context.Context, key string) (int, error) {
	code, err := periodScript.Run(ctx, h.limitStore, []string{h.keyPrefix + key}, []string{
		strconv.Itoa(h.quota),
		strconv.Itoa(h.calcExpireSeconds()),
	}).Int()
	if err != nil {
		return Unknown, err
	}

	switch code {
	case internalOverQuota:
		return OverQuota, nil
	case internalAllowed:
		return Allowed, nil
	case internalHitQuota:
		return HitQuota, nil
	default:
		return Unknown, ErrUnknownCode
	}
}

func (h *PeriodLimit) calcExpireSeconds() int {
	// 如果 align 为 true，则计算距离当前时间段结束还有多长时间
	if h.align {
		now := time.Now()
		_, offset := now.Zone()
		unix := now.Unix() + int64(offset)
		return h.period - int(unix%int64(h.period))
	}
	// 如果 align 为 false，则直接返回时间段长度
	return h.period

}

// Align 返回一个 PeriodOption，用于对 limiter 进行对齐配置
func Align() PeriodOption {
	return func(l *PeriodLimit) {
		l.align = true
	}
}
