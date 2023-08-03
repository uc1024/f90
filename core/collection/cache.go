package collection

import (
	"container/list"

	"log"
	"sync"
	"time"

	"github.com/uc1024/f90/core/mathx"
	"github.com/uc1024/f90/core/syncx"
)

var (
	default_cache_expire   = time.Minute * 30
	default_cache_slots    = 300 // * 默认一轮的槽位
	default_cache_interval = time.Minute
)

const (
	cache_time_offset = 0.01
)

type (
	cacheOption func(cache *Cache)

	Cache struct {
		name        string
		lock        sync.Mutex
		data        map[string]interface{}
		expire      time.Duration
		timingWheel *TimingWheel
		lruCache    ILru
		barrier     syncx.ShareResults
		unstable    mathx.Unstable
		execute     func(k, v interface{})
	}
)

func NewCache(expire time.Duration, opt ...cacheOption) (cache *Cache, err error) {
	cache = &Cache{
		data:     make(map[string]interface{}),
		expire:   expire,
		barrier:  syncx.NewShareCall(),
		unstable: mathx.NewUnstable(cache_time_offset),
		lruCache: &emptyLru{},
	}

	// * 过期处理
	cache.execute = func(k, v interface{}) {
		key, ok := k.(string)
		if !ok {
			return
		}
		cache.Del(key)
	}

	for _, o := range opt {
		o(cache)
	}

	wheel, err := NewTimingWheel(time.Second, default_cache_slots, cache.execute)

	if err != nil {
		return
	}

	cache.timingWheel = wheel

	return
}

func (c *Cache) Del(key string) {
	c.lock.Lock()
	delete(c.data, key)    // * 删除数据
	c.lruCache.Remove(key) // * 从策略中移除key
	c.lock.Unlock()
	c.timingWheel.RemoveTimer(key) // * 缓存失效定时器移除
}

func (c *Cache) Get(key string) (value interface{}, b bool) {
	return c.doGet(key)
}

func (c *Cache) Set(key string, value interface{}) {
	c.SetWithExpire(key, value, c.expire)
}

func (c *Cache) SetWithExpire(k string, v interface{}, expire time.Duration) {
	c.lock.Lock()
	_, cover := c.data[k] // * 判断是覆盖还是新增
	c.data[k] = v
	c.lruCache.Add(k)
	c.lock.Unlock()

	t := c.unstable.AroundDuration(expire)
	if cover {
		c.timingWheel.MoveTimer(k, t)
	} else {
		c.timingWheel.SetTimer(k, v, t)
	}
}

func (c *Cache) doGet(key string) (interface{}, bool) {
	defer c.lock.Unlock()
	c.lock.Lock()
	value, ok := c.data[key]
	if ok {
		return value, ok
	}
	return nil, ok
}

func (c *Cache) Take(key string,
	fetch func() (interface{}, error)) (value interface{}, err error) {

	value, ok := c.doGet(key)

	if ok {
		return
	}

	value, err = c.barrier.Do(key, func() (interface{}, error) {
		// * 在调用fetch前在尝试一次是否存在缓存
		if val, ok := c.doGet(key); ok {
			return val, nil
		}
		results, e := fetch()
		if e != nil {
			return nil, e
		}
		c.Set(key, results)
		return results, nil
	})

	return
}

func (c *Cache) Refresh(k string, v interface{}) {
	c.RefreshWithExpire(k, v, c.expire)
}

// * 刷新一定是在之前的定时器跑完了继续推入任务
func (c *Cache) RefreshWithExpire(k string, v interface{},
	expire time.Duration) {

	c.lock.Lock()
	c.data[k] = v
	c.lruCache.Add(k)
	c.lock.Unlock()

	t := c.unstable.AroundDuration(expire)
	if err := c.timingWheel.SetTimer(k, v, t); err != nil {
		log.Println(err)
	}
}

func (c *Cache) onExile(key string) {
	// * 不要加锁。实际上已经锁了
	// * @Del中上锁了才调用的
	delete(c.data, key)
	c.timingWheel.RemoveTimer(key)
}

func (c *Cache) Size() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return len(c.data)
}

// ------------------------CacheOption-----------------------------------------

// * 设置缓存策略
func SetCacheLimit(i int) cacheOption {
	return func(cache *Cache) {
		if i > 0 {
			cache.lruCache = NewKeyLru(i, cache.onExile)
		}
	}
}

func SetCacheName(name string) cacheOption {
	return func(cache *Cache) {
		cache.name = name
	}
}

func SetCacheExecute(fn func(k, v interface{})) cacheOption {
	return func(cache *Cache) {
		cache.execute = fn
	}
}

// -----------------------------------------------------------------------------

type (
	// * lru 缓存淘汰策略
	ILru interface {
		Add(key string) int
		Remove(key string) int
	}

	// * key调用策略 如果超出限制则丢弃最尾部的缓存
	// & 统计key的调用,使用最少得key被优先移除
	keyLru struct {
		limit      int
		exileOrder *list.List               // * 移除顺序,尾部优先移除
		elements   map[string]*list.Element // * 储存key
		onExile    func(key string)         // * 缓存key移除通知
	}
)

func NewKeyLru(limit int, on_exile func(key string)) ILru {
	return &keyLru{
		limit:      limit,
		exileOrder: list.New(),
		elements:   make(map[string]*list.Element),
		onExile:    on_exile,
	}
}

// * key 调用++
func (lru *keyLru) Add(key string) int {

	if elem, ok := lru.elements[key]; ok {
		// *  如果元素存在,放在最前面.
		lru.exileOrder.MoveToFront(elem)
		return 0
	}

	defer lru.verifySize()
	// * 新增key处理默认放在最前面
	lru.elements[key] = lru.exileOrder.PushFront(key)
	return 1
}

// * 移除
func (lru *keyLru) Remove(key string) int {
	if elem, ok := lru.elements[key]; ok {
		lru.removeElements(elem)
		return 1
	}
	return 0
}

// * 移除一个元素
func (lru *keyLru) removeElements(e *list.Element) {
	lru.exileOrder.Remove(e)
	key := e.Value.(string)
	delete(lru.elements, key)
	lru.onExile(key)
}

// * 检查长度超出则销毁尾部key
func (lru *keyLru) verifySize() {
	if lru.limit < lru.exileOrder.Len() {
		lru.removeOldest()
	}
}

// * 移除尾部
func (lru *keyLru) removeOldest() {
	elem := lru.exileOrder.Back()
	if elem != nil {
		lru.removeElements(elem)
	}
}

// -----------------------------------------------------------------------------
// * 空实现策略
type emptyLru struct{}

func (lru *emptyLru) Add(key string) int    { return 1 }
func (lru *emptyLru) Remove(key string) int { return 1 }
