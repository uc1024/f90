package throttlex

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"github.com/uc1024/f90/core/throttlex/script"
)

type LimitSender struct {
	client  *redis.Client
	options *SendLimitOptions
}

type SendLimitOptions struct {
	Count     int64         `json:"count"`      // * 每天的次数
	Period    time.Duration `json:"period"`     // * 每次的间隔秒
	KeyPrefix string        `json:"key_prefix"` // * key前缀
}

type SendLimitResult struct {
	Count       int64 `json:"count"`        // * 每天的次数
	Period      int64 `json:"period"`       // * 每次的间隔秒
	LastTime    int64 `json:"last_time"`    // * 最后一次发送时间
	SendedCount int64 `json:"sended_count"` // * 当天已发送次数
	WaitTime    int64 `json:"wait_time"`    // * 等待时间
}

type SendLimitOptionsFunc func(*SendLimitOptions)

func NewLimitSender(client *redis.Client, opt ...SendLimitOptionsFunc) (*LimitSender, error) {
	options := &SendLimitOptions{
		Count:     10,
		Period:    time.Second * 60,
		KeyPrefix: "limit_send",
	}
	for _, o := range opt {
		o(options)
	}

	// * 预加载
	script.GetLimitSendScript(context.Background(), client)

	return &LimitSender{client, options}, nil
}

func (ls *LimitSender) Key(s string) string {
	return fmt.Sprintf("%s:%s", ls.options.KeyPrefix, s)
}

func (ls *LimitSender) SendLimit(key string) (*SendLimitResult, error) {
	return ls.SendLimitCtx(context.Background(), key)
}

func (ls *LimitSender) SendLimitCtx(ctx context.Context, key string) (*SendLimitResult, error) {

	res, err := script.LimitSendScript.EvalSha(ctx, ls.client, []string{ls.Key(key)},
		ls.options.Period.Seconds(),
		ls.options.Count,
		time.Now().Add(24*time.Hour).Unix(), // * 默认的过期时间
	).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to execute limit_send_script: %v", err)
	}

	values, ok := res.([]interface{})
	if !ok || len(values) != 2 {
		return nil, fmt.Errorf("unexpected result format from limit_send_script")

	}

	count, ok := values[0].(int64)
	if !ok {
		return nil, fmt.Errorf("unexpected count type from limit_send_script")
	}

	json_str, ok := values[1].(string)

	if !ok {
		return nil, fmt.Errorf("unexpected json bytes type from limit_send_script %T", values[1])
	}

	var limit_ary []string
	if err := json.Unmarshal([]byte(json_str), &limit_ary); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json from limit_send_script: %v", err)
	}

	var send_limit = make(map[string]string)
	for i := 0; i < len(limit_ary)-1; i += 2 {
		send_limit[limit_ary[i]] = limit_ary[i+1]
	}
	// ( "count", "last_time"))

	result := SendLimitResult{
		Count:       ls.options.Count,
		Period:      int64(ls.options.Period.Seconds()),
		SendedCount: cast.ToInt64(send_limit["count"]),
		LastTime:    cast.ToInt64(send_limit["last_time"]),
		WaitTime:    count,
	}

	return &result, nil
}

func (ls *LimitSender) GetLimitWithCtx(ctx context.Context, key string) (*SendLimitResult, error) {

	result, err := ls.client.HGetAll(ctx, ls.Key(key)).Result()

	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return &SendLimitResult{
			Count:       ls.options.Count,
			Period:      int64(ls.options.Period.Seconds()),
			SendedCount: 0,
			LastTime:    0,
			WaitTime:    0,
		}, nil
	}
	var send_limit = make(map[string]string)
	for i := 0; i < len(result)-1; i += 2 {
		ii := strconv.Itoa(i)
		iii := strconv.Itoa(i + 1)
		send_limit[result[ii]] = result[iii]

	}

	info := &SendLimitResult{
		Count:       ls.options.Count,
		Period:      int64(ls.options.Period.Seconds()),
		SendedCount: cast.ToInt64(send_limit["count"]),
		LastTime:    cast.ToInt64(send_limit["last_time"]),
		WaitTime:    0,
	}

	if ls.options.Count > cast.ToInt64(send_limit["count"]) {
	}

	return info, nil
}
