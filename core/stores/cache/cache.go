package cache

import (
	"context"
	"time"


)

type ( // Cache interface is used to define the cache implementation.
	Cache interface {
		// Del deletes cached values with keys.
		Del(keys ...string) error
		// DelCtx deletes cached values with keys.
		DelCtx(ctx context.Context, keys ...string) error
		// Get gets the cache with key and fills into v.
		Get(key string, val interface{}) error
		// GetCtx gets the cache with key and fills into v.
		GetCtx(ctx context.Context, key string, val interface{}) error
		// IsNotFound checks if the given error is the defined errNotFound.
		IsNotFound(err error) bool
		// Set sets the cache with key and v, using c.expiry.
		Set(key string, val interface{}) error
		// SetCtx sets the cache with key and v, using c.expiry.
		SetCtx(ctx context.Context, key string, val interface{}) error
		// SetWithExpire sets the cache with key and v, using given expire.
		SetWithExpire(key string, val interface{}, expire time.Duration) error
		// SetWithExpireCtx sets the cache with key and v, using given expire.
		SetWithExpireCtx(ctx context.Context, key string, val interface{}, expire time.Duration) error
		// Take takes the result from cache first, if not found,
		// query from DB and set cache using c.expiry, then return the result.
		Take(val interface{}, key string, query func(val interface{}) error) error
		// TakeCtx takes the result from cache first, if not found,
		// query from DB and set cache using c.expiry, then return the result.
		TakeCtx(ctx context.Context, val interface{}, key string, query func(val interface{}) error) error
		// TakeWithExpire takes the result from cache first, if not found,
		// query from DB and set cache using given expire, then return the result.
		TakeWithExpire(val interface{}, key string, query func(val interface{}, expire time.Duration) error) error
		// TakeWithExpireCtx takes the result from cache first, if not found,
		// query from DB and set cache using given expire, then return the result.
		TakeWithExpireCtx(ctx context.Context, val interface{}, key string,
			query func(val interface{}, expire time.Duration) error) error
	}
)


