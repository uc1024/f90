package executors

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/uc1024/f90/core/timex"
)

const threshold = 10

type container struct {
	interval time.Duration
	tasks    []int
	execute  func(tasks any)
}

func newContainer(interval time.Duration, execute func(tasks any)) *container {
	return &container{
		interval: interval,
		execute:  execute,
	}
}

func (c *container) AddTask(task any) bool {
	c.tasks = append(c.tasks, task.(int))
	return len(c.tasks) > threshold
}

func (c *container) Execute(tasks any) {
	if c.execute != nil {
		c.execute(tasks)
	} else {
		time.Sleep(c.interval)
	}
}

func (c *container) RemoveAll() any {
	tasks := c.tasks
	c.tasks = nil
	return tasks
}

func TestPeriodicalExecutor_Sync(t *testing.T) {
	var done int32
	exec := NewPeriodicalExecutor(time.Second, newContainer(time.Millisecond*500, nil))
	exec.Sync(func() {
		atomic.AddInt32(&done, 1)
	})
	assert.Equal(t, int32(1), atomic.LoadInt32(&done))
}

func TestPeriodicalExecutor_QuitGoroutine(t *testing.T) {
	ticker := timex.NewFakeTicker()
	exec := NewPeriodicalExecutor(time.Millisecond, newContainer(time.Millisecond, nil))
	exec.newTicker = func(d time.Duration) timex.Ticker {
		return ticker
	}
	routines := runtime.NumGoroutine()
	exec.Add(1)
	ticker.Tick()
	ticker.Wait(time.Millisecond * idleRound * 2)
	ticker.Tick()
	ticker.Wait(time.Millisecond * idleRound)
	assert.Equal(t, routines, runtime.NumGoroutine())
	// proc.Shutdown()
}

func TestPeriodicalExecutor_Bulk(t *testing.T) {
	ticker := timex.NewFakeTicker()
	var vals []int
	// avoid data race
	var lock sync.Mutex
	exec := NewPeriodicalExecutor(time.Millisecond, newContainer(time.Millisecond, func(tasks any) {
		t := tasks.([]int)
		for _, each := range t {
			lock.Lock()
			vals = append(vals, each)
			lock.Unlock()
		}
	}))
	exec.newTicker = func(d time.Duration) timex.Ticker {
		return ticker
	}
	for i := 0; i < threshold*10; i++ {
		if i%threshold == 5 {
			time.Sleep(time.Millisecond * idleRound * 2)
		}
		exec.Add(i)
	}
	ticker.Tick()
	ticker.Wait(time.Millisecond * idleRound * 2)
	ticker.Tick()
	ticker.Tick()
	ticker.Wait(time.Millisecond * idleRound)
	var expect []int
	for i := 0; i < threshold*10; i++ {
		expect = append(expect, i)
	}

	lock.Lock()
	assert.EqualValues(t, expect, vals)
	lock.Unlock()
}

func TestPeriodicalExecutor_Panic(t *testing.T) {
	// avoid data race
	var lock sync.Mutex
	ticker := timex.NewFakeTicker()

	var (
		executedTasks []int
		expected      []int
	)
	executor := NewPeriodicalExecutor(time.Millisecond, newContainer(time.Millisecond, func(tasks any) {
		tt := tasks.([]int)
		lock.Lock()
		executedTasks = append(executedTasks, tt...)
		lock.Unlock()
		if tt[0] == 0 {
			panic("test")
		}
	}))
	executor.newTicker = func(duration time.Duration) timex.Ticker {
		return ticker
	}
	for i := 0; i < 30; i++ {
		executor.Add(i)
		expected = append(expected, i)
	}
	ticker.Tick()
	ticker.Tick()
	time.Sleep(time.Millisecond)
	lock.Lock()
	assert.Equal(t, expected, executedTasks)
	lock.Unlock()
}

func TestPeriodicalExecutor_FlushPanic(t *testing.T) {
	var (
		executedTasks []int
		expected      []int
		lock          sync.Mutex
	)
	executor := NewPeriodicalExecutor(time.Millisecond, newContainer(time.Millisecond, func(tasks any) {
		tt := tasks.([]int)
		lock.Lock()
		executedTasks = append(executedTasks, tt...)
		lock.Unlock()
		if tt[0] == 0 {
			panic("flush panic")
		}
	}))
	for i := 0; i < 8; i++ {
		executor.Add(i)
		expected = append(expected, i)
	}
	executor.Flush()
	lock.Lock()
	assert.Equal(t, expected, executedTasks)
	lock.Unlock()
}
