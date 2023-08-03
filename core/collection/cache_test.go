package collection

import (
	"testing"
	"time"

	"github.com/uc1024/f90/core/syncx"
)

func TestCache(t *testing.T) {
	c, _ := NewCache(time.Second*10, SetCacheLimit(100))
	c.SetWithExpire("test", 100, time.Second*5)
	v, _ := c.Get("test")
	t.Log("value", v)

	x := syncx.NewCond()
	x.WaitWithTimeout(8 * time.Second)
	v1, _ := c.Get("test")
	t.Log("11s value1", v1)
}
