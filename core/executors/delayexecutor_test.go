package executors

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/uc1024/f90/core/syncx"
	"github.com/uc1024/f90/core/timex"
)

func TestDelayExecutor_Trigger(t *testing.T) {
	start := timex.Now()
	delay := time.Second * 1
	cond := syncx.NewCond()
	delayExecutor := NewDelayExecutor(func() {
		assert.True(t, timex.Since(start) >= time.Second)
		cond.Signal()
	}, delay)
	delayExecutor.Trigger()
	cond.Wait()
}

func TestDelayExecutor_ManyTrigger(t *testing.T) {
	delay := time.Millisecond * 1
	cond := syncx.NewCond()
	value := int32(0)
	delayExecutor := NewDelayExecutor(func() {
		atomic.AddInt32(&value, 1)
	}, delay)
	delayExecutor.Trigger()
	delayExecutor.Trigger()
	cond.WaitWithTimeout(delay * 3)
	delayExecutor.Trigger()
	cond.WaitWithTimeout(delay * 2)
	assert.Equal(t, int32(2), atomic.LoadInt32(&value))
}
