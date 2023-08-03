package collection

import (
	"testing"
	"time"

	"github.com/uc1024/f90/core/syncx"
)

func TestTimingwheel(t *testing.T) {
	wheel, _ := NewTimingWheel(time.Second, 6, func(key, value interface{}) {
		t.Log(key)
		value.(func())()
	})
	wheel.SetTimer("test1", func() {
		t.Log("a1")
	}, time.Second*1)
	wheel.SetTimer("test2", func() {
		t.Log("a2")
	}, time.Second*10)
	wheel.SetTimer("test3", func() {
		t.Log("a3")
	}, time.Second*20)

	s := syncx.NewCond()
	s.WaitWithTimeout(30 * time.Second)
}
