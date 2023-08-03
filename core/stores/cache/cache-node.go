package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"strings"
	"time"

	"github.com/uc1024/f90/core/mathx"
	"github.com/uc1024/f90/core/slogx"
	"github.com/uc1024/f90/core/syncx"
	"github.com/redis/go-redis/v9"
)

const (
	// * 缓存占位符
	notFoundPlaceholder = "*"
	// * 到期时间偏移
	expiryDeviation = 0.05
)

// indicates there is no such value associate with the key
var errPlaceholder = errors.New("placeholder")

type (
	cacheNode struct {
		rds            *redis.Client
		expiry         time.Duration // default expiry time
		notFoundExpiry time.Duration // * 占位锁过期时间
		barrier        syncx.ShareResults
		unstable       mathx.Unstable
		errNotFound    error
	}
)

func NewNode(rds *redis.Client, errNotFound error, opts ...Option) Cache {
	o := newOptions(opts...)
	return cacheNode{
		rds:            rds,
		expiry:         o.Expiry,
		notFoundExpiry: o.NotFoundExpiry,
		barrier:        syncx.NewShareCall(),
		unstable:       mathx.NewUnstable(expiryDeviation),
		errNotFound:    errNotFound,
	}
}

func (c cacheNode) IsNotFound(err error) bool {
	return errors.Is(err, c.errNotFound)
}

func (c cacheNode) doGetCache(ctx context.Context,
	key string, v interface{}) (err error) {

	result := c.rds.Get(ctx, key)

	value := ""
	if result.Err() != nil {
		if !errors.Is(result.Err(), redis.Nil) {
			return result.Err()
		}
	} else {
		value = result.Val()
	}

	if len(value) == 0 {
		return c.errNotFound
	}
	if value == notFoundPlaceholder {
		return errPlaceholder
	}

	err = c.processCache(ctx, key, value, v)

	return
}

func (c cacheNode) Take(val interface{}, key string, query func(val interface{}) error) error {
	return c.TakeCtx(context.Background(), val, key, query)
}

func (c cacheNode) TakeCtx(ctx context.Context, val interface{}, key string, query func(val interface{}) error) error {
	return c.doTake(ctx, val, key, query, func(v interface{}) error {
		return c.SetCtx(ctx, key, v)
	})
}

func (c cacheNode) TakeWithExpire(val interface{}, key string, query func(val interface{}, expire time.Duration) error) error {
	return c.TakeWithExpireCtx(context.Background(), val, key, query)
}

func (c cacheNode) TakeWithExpireCtx(ctx context.Context, val interface{}, key string, query func(val interface{}, expire time.Duration) error) error {
	cacheval := func(v interface{}) error {
		return c.SetWithExpireCtx(ctx, key, v, c.expiry)
	}
	call_query := func(v interface{}) error {
		return query(v, c.expiry)
	}
	return c.doTake(ctx, val, key, call_query, cacheval)
}

func (c cacheNode) doTake(ctx context.Context, v interface{}, key string,
	query func(v interface{}) error, cacheVal func(v interface{}) error) error {

	val, fresh, err := c.barrier.DoEx(key, func() (interface{}, error) {

		if err := c.doGetCache(ctx, key, v); err != nil {
			if err == errPlaceholder {
				return nil, c.errNotFound
			} else if err != c.errNotFound {
				return nil, err
			}

			if err = query(v); err == c.errNotFound {

				if err = c.setCacheWithNotFound(ctx, key); err != nil {
					fmt.Println(err)
				}

				return nil, c.errNotFound

			} else if err != nil {
				return nil, err
			}

			// call the callback setting cache value
			if err = cacheVal(v); err != nil {
				fmt.Println(err)
			}
		}

		return json.Marshal(v)
	})

	if err != nil {
		return err
	}

	// not yourself get but can use
	if fresh {
		return nil
	}

	return json.Unmarshal(val.([]byte), v)
}

// delete the selected cache
func (c cacheNode) Del(keys ...string) error {
	return c.DelCtx(context.Background(), keys...)
}

func (c cacheNode) DelCtx(ctx context.Context, keys ...string) (err error) {
	if len(keys) == 0 {
		return nil
	}
	for _, key := range keys {
		if err = c.rds.Del(ctx, key).Err(); err != nil {
			fmt.Printf("cache delete key:%s err:%s \n", key, err.Error())
		} else {
			slogx.Default.Error(nil, "failed to clear cache with keys: %s, error: %v", strings.Join(keys, ","), err)
			c.asyncRetryDelCache(keys...)
		}
	}
	return nil
}

func (c cacheNode) Get(key string, val interface{}) error {
	return c.GetCtx(context.Background(), key, val)
}

func (c cacheNode) GetCtx(ctx context.Context, key string, val interface{}) error {
	err := c.doGetCache(ctx, key, val)
	if err == errPlaceholder {
		return c.errNotFound
	}

	return err
}

func (c cacheNode) processCache(ctx context.Context,
	key, data string, v interface{}) error {

	err := json.Unmarshal([]byte(data), v)
	if err == nil {
		return nil
	}

	// 到这错误说明内容格式不正常
	// 删除异常数据
	slogx.Default.Error(nil, fmt.Sprintf("failed to unmarshal cache with key: %s, error: %v", key, err))

	// 尝试删除错误缓存
	c.DelCtx(ctx, key)
	return err
}

// setting occupy a seat
func (c cacheNode) setCacheWithNotFound(ctx context.Context,
	key string) (err error) {

	result := c.rds.SetNX(ctx, key, notFoundPlaceholder, c.notFoundExpiry)

	if result.Err() != nil {
		return result.Err()
	}
	return
}

func (c cacheNode) Set(key string,
	val interface{}) error {

	return c.SetWithExpireCtx(context.Background(),
		key, val, time.Duration(c.expiry))
}

func (c cacheNode) SetCtx(ctx context.Context, key string,
	val interface{}) error {

	return c.SetWithExpireCtx(ctx, key, val,
		time.Duration(c.expiry))
}

func (c cacheNode) SetWithExpire(key string,
	val interface{}, expire time.Duration) error {

	return c.SetWithExpireCtx(context.Background(), key, val, expire)
}

/*
只在键 key 不存在的情况下， 将键 key 的值设置为 value 。
若键 key 已经存在， 则 SETNX 命令不做任何动作。
*/
func (c cacheNode) SetWithExpireCtx(ctx context.Context, key string,
	val interface{}, expire time.Duration) (err error) {

	data, err := json.Marshal(val)
	if err != nil {
		return err
	}

	result := c.rds.SetNX(ctx, key, data, c.unstable.AroundDuration(expire))

	if !result.Val() {
		return fmt.Errorf("set cache key:%s value:%s err:%s", key, data, result.Err().Error())
	}

	if result.Err() != nil {
		return result.Err()
	}

	return nil

}

// delay delete caches
func (c cacheNode) asyncRetryDelCache(keys ...string) {
	AddCleanTask(func() error {
		result := c.rds.Del(context.Background(), keys...)
		return result.Err()
	}, keys...)
}
