package records

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

/*
	使用redis实现用户浏览记录
*/

const defaultPrefixKey = "BROWSE:HISTORY:"

type (
	BrowseHistoryOptions struct {
		PrefixKey string
		Expire    time.Duration
		Redis     *redis.Client
		Max       int
	}

	BrowseHistory struct {
		options BrowseHistoryOptions
		rds     *redis.Client
	}

	SetBrowseHistoryOptions func(*BrowseHistoryOptions)
)

func NewBrowseHistory(opts ...SetBrowseHistoryOptions) *BrowseHistory {
	options := BrowseHistoryOptions{
		PrefixKey: defaultPrefixKey,
		Max:       30,
	}
	for _, v := range opts {
		v(&options)
	}
	return &BrowseHistory{
		options: options,
		rds:     options.Redis,
	}
}

func (history *BrowseHistory) key(id string) string {
	return fmt.Sprintf("%s%s", history.options.PrefixKey, id)
}

func (history *BrowseHistory) Push(
	ctx context.Context, key, value string) (err error) {

	err = history.rds.LPush(ctx, history.key(key), value).Err()
	if err != nil {
		return
	}

	// * 修剪记录
	history.rds.LTrim(ctx, key, 0, int64(history.options.Max))
	return
}

func (history *BrowseHistory) Len(
	ctx context.Context, Key string) (int64, error) {

	return history.rds.LLen(ctx, history.key(Key)).Result()
}

func (history *BrowseHistory) Range(
	ctx context.Context, Key string) ([]string, error) {

	rsp := history.rds.LRange(ctx, history.key(Key),
		0, int64(history.options.Max))
	if rsp.Err() != nil {
		return nil, rsp.Err()
	}
	return rsp.Val(), nil
}
