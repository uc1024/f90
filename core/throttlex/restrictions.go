package throttlex

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/uc1024/f90/core/throttlex/script"
)

const (
	/*
		--   @Last visit            上次访问时间
		--   @Overclock             被限制时调用的次数
	*/
	FieldLastVisit = "last_visit"
	FieldOverclock = "overclock"

	// default
	DefaultAction   = "default_action"
	DefaultInterval = 300
	DefaultName     = "RESTRICTIONS"
	/*
		status
		-- 1 允许访问
		-- 2 时间间隔限制中
	*/
	Allow        = 1
	TimeFail     = 2
	FrequentFail = 3
)

// //go:embed script/allow.lua
// var allowScript string

type (
	RestrictionsOptions func(*restrictionsOptions)

	restrictionsOptions struct {
		name     string
		action   string // * 执行的动作
		interval int64  // * 访问key的间隔(毫秒)

	}

	Restrictions struct {
		options *restrictionsOptions
		store   *redis.Client
	}

	RestrictionLimit struct {
		LastVisit time.Time
		Overclock int
	}
)

func SetRestrictionsAction(value string) RestrictionsOptions {
	return func(ro *restrictionsOptions) {
		ro.action = value
	}
}

func SetRestrictionsInterval(millisecond int64) RestrictionsOptions {
	return func(ro *restrictionsOptions) {
		ro.interval = millisecond
	}
}

func SetRestrictionsName(value string) RestrictionsOptions {
	return func(ro *restrictionsOptions) {
		ro.name = value
	}
}

func NewRestrictions(cli *redis.Client,
	opts ...RestrictionsOptions) *Restrictions {

	options := &restrictionsOptions{
		name:     DefaultName,
		action:   DefaultAction,
		interval: DefaultInterval,
	}

	for _, v := range opts {
		v(options)
	}

	return &Restrictions{
		store:   cli,
		options: options,
	}
}

func (r *Restrictions) Key(key string) string {
	return fmt.Sprintf("%s:%s:%s", DefaultName, r.options.action, key)
}

// * 判断是否允许
func (r *Restrictions) Allow(ctx context.Context, key string) (status int, err error) {

	// * 当前访问时间
	ARGV_1 := time.Now().UnixMilli()
	ARGV_2 := r.options.interval

	status, err = script.RestrictionsAllowScript.Run(ctx, r.store,
		[]string{r.Key(key)}, ARGV_1, ARGV_2).Int()

	if err != nil {
		return
	}

	return
}
