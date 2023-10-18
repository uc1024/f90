package executors


import (
	"time"

	"github.com/uc1024/f90/core/syncx"
	"github.com/uc1024/f90/core/timex"
)


// 确保在给定的时间间隔内只执行一次任务，防止任务过于频繁地执行。
// 在间隔内的任务会被丢弃，只有在间隔外的任务才会被执行。
// A LessExecutor is an executor to limit execution once within given time interval.
type LessExecutor struct {
	threshold time.Duration
	lastTime  *syncx.AtomicDuration
}

// NewLessExecutor returns a LessExecutor with given threshold as time interval.
func NewLessExecutor(threshold time.Duration) *LessExecutor {
	return &LessExecutor{
		threshold: threshold,
		lastTime:  syncx.NewAtomicDuration(),
	}
}

// DoOrDiscard executes or discards the task depends on if
// another task was executed within the time interval.
func (le *LessExecutor) DoOrDiscard(execute func()) bool {
	now := timex.Now()
	lastTime := le.lastTime.Load()
	if lastTime == 0 || lastTime+le.threshold < now {
		le.lastTime.Set(now)
		execute()
		return true
	}

	return false
}
