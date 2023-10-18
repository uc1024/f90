package executors

import (
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/uc1024/f90/core/syncx"
	"github.com/uc1024/f90/core/threadingx"
	"github.com/uc1024/f90/core/timex"
)

const idleRound = 10

type (
	// TaskContainer 接口定义了一个用于存储周期性执行任务的容器类型。
	TaskContainer interface {
		// AddTask 将任务添加到容器中。
		// 如果添加后需要刷新容器，则返回 true。
		AddTask(task any) bool
		// Execute 在刷新时处理容器中收集到的任务。
		Execute(tasks any)
		// RemoveAll 移除所有的任务，并返回它们。
		RemoveAll() any
	}

	// PeriodicalExecutor 是一个周期性执行任务的执行器。
	PeriodicalExecutor struct {
		commander   chan any                                  // 用于接收任务的通道
		interval    time.Duration                             // 执行任务的时间间隔
		container   TaskContainer                             // 任务容器
		waitGroup   sync.WaitGroup                            // 用于等待任务执行完成的等待组
		wgBarrier   syncx.Barrier                             // 用于避免在调用 wg.Add/Done/Wait(...) 时的竞态条件
		confirmChan chan placeholderType                      // 确认通道，用于在任务执行完成后发送确认信号
		inflight    int32                                     // 正在执行的任务数量
		guarded     bool                                      // 是否受保护，用于避免竞态条件
		newTicker   func(duration time.Duration) timex.Ticker // 创建新的定时器的函数
		lock        sync.Mutex                                // 用于保护关键区域的互斥锁
	}
)

// NewPeriodicalExecutor 使用给定的间隔和容器创建一个 PeriodicalExecutor。
func NewPeriodicalExecutor(interval time.Duration, container TaskContainer) *PeriodicalExecutor {
	executor := &PeriodicalExecutor{
		commander:   make(chan any, 1),          // 创建一个缓冲区大小为1的通道
		interval:    interval,                   // 设置执行任务的时间间隔
		container:   container,                  // 设置任务容器
		confirmChan: make(chan placeholderType), // 创建确认通道
		newTicker: func(d time.Duration) timex.Ticker {
			return timex.NewTicker(d) // 创建一个新的定时器
		},
	}
	// TODO 安全关闭
	// proc.AddShutdownListener(func() {
	// 	executor.Flush() // 在程序关闭时执行 Flush()
	// })

	return executor
}

// Add 将任务添加到 pe 中。
func (pe *PeriodicalExecutor) Add(task any) {
	if vals, ok := pe.addAndCheck(task); ok {
		pe.commander <- vals
		<-pe.confirmChan
	}
}

// Flush 强制 pe 执行任务。
func (pe *PeriodicalExecutor) Flush() bool {
	pe.enterExecution()
	return pe.executeTasks(func() any {
		pe.lock.Lock()
		defer pe.lock.Unlock()
		return pe.container.RemoveAll() // 移除并返回所有任务
	}())
}

// Sync 允许调用者在 pe 中以线程安全的方式运行 fn，特别是对于底层容器。
func (pe *PeriodicalExecutor) Sync(fn func()) {
	pe.lock.Lock()
	defer pe.lock.Unlock()
	fn() // 执行传入的函数
}

// Wait 等待执行完成。
func (pe *PeriodicalExecutor) Wait() {
	pe.Flush() // 执行 Flush()
	pe.wgBarrier.Guard(func() {
		pe.waitGroup.Wait() // 等待任务执行完成
	})
}

// addAndCheck 将任务添加到 pe 中，并检查是否需要执行任务。
func (pe *PeriodicalExecutor) addAndCheck(task any) (any, bool) {
	pe.lock.Lock()
	defer func() {
		if !pe.guarded {
			pe.guarded = true
			defer pe.backgroundFlush() // 在后台进行刷新
		}
		pe.lock.Unlock()
	}()

	if pe.container.AddTask(task) { // 尝试将任务添加到容器中
		atomic.AddInt32(&pe.inflight, 1) // 原子操作，增加正在执行的任务数量
		return pe.container.RemoveAll(), true
	}

	return nil, false
}

// backgroundFlush 在后台执行刷新。
func (pe *PeriodicalExecutor) backgroundFlush() {
	go func() {
		defer pe.Flush() // 在退出协程之前执行刷新以避免丢失任务

		ticker := pe.newTicker(pe.interval) // 创建定时器
		defer ticker.Stop()                 // 停止定时器

		var commanded bool
		last := timex.Now()
		for {
			select {
			case vals := <-pe.commander: // 从命令通道中接收任务
				commanded = true
				atomic.AddInt32(&pe.inflight, -1) // 减少正在执行的任务数量
				pe.enterExecution()               // 进入执行状态
				pe.confirmChan <- struct{}{}      // 发送确认信号
				pe.executeTasks(vals)             // 执行任务
				last = timex.Now()                // 更新上次执行的时间
			case <-ticker.Chan(): // 定时器触发
				if commanded {
					commanded = false
				} else if pe.Flush() { // 执行强制刷新
					last = timex.Now() // 更新上次执行的时间
				} else if pe.shallQuit(last) { // 是否应该退出
					return
				}
			}
		}
	}()
}

// doneExecution 标记执行完成。
func (pe *PeriodicalExecutor) doneExecution() {
	pe.waitGroup.Done()
}

// enterExecution 进入执行状态。
func (pe *PeriodicalExecutor) enterExecution() {
	pe.wgBarrier.Guard(func() {
		pe.waitGroup.Add(1) // 添加一个任务到等待组
	})
}

// executeTasks 执行任务。
func (pe *PeriodicalExecutor) executeTasks(tasks any) bool {
	defer pe.doneExecution()

	ok := pe.hasTasks(tasks) // 检查是否有任务
	if ok {
		threadingx.RunSafe(func() {
			pe.container.Execute(tasks) // 执行任务
		})
	}

	return ok
}

// hasTasks 检查是否有任务。
func (pe *PeriodicalExecutor) hasTasks(tasks any) bool {
	if tasks == nil {
		return false
	}

	val := reflect.ValueOf(tasks) // 获取值的反射信息
	switch val.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice: // 对于数组、通道、映射和切片类型
		return val.Len() > 0 // 判断长度是否大于0
	default:
		// 未知类型，让调用者执行
		return true
	}
}

// shallQuit 检查是否应该退出。
func (pe *PeriodicalExecutor) shallQuit(last time.Duration) (stop bool) {
	if timex.Since(last) <= pe.interval*idleRound {
		return
	}

	pe.lock.Lock()
	if atomic.LoadInt32(&pe.inflight) == 0 { // 如果没有正在执行的任务
		pe.guarded = false
		stop = true
	}
	pe.lock.Unlock()

	return
}
