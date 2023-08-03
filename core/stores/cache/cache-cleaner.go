package cache

import (
	"fmt"

	"strings"
	"time"

	"github.com/uc1024/f90/core/collection"
	"github.com/uc1024/f90/core/slogx"
	"github.com/uc1024/f90/core/stringx"
	"github.com/uc1024/f90/core/threadingx"
)

const cleanWorkers = 5
const numSlots = 300

var (
	// * time wheel
	timingWheel *collection.TimingWheel
	// * task runner  number of workers
	taskRunner = threadingx.NewTaskRunner(cleanWorkers)
)

/*
缓存清理任务
*/
type cacheCleaner struct {
	delay time.Duration
	task  func() error
	keys  []string
}

func init() {
	var err error
	timingWheel, err = collection.NewTimingWheel(time.Second, numSlots, clean)
	slogx.Default.MustSucc(nil, err)
}

func AddCleanTask(task func() error, keys ...string) {
	timingWheel.SetTimer(stringx.Randn(8), cacheCleaner{
		delay: time.Second,
		task:  task,
		keys:  keys,
	}, time.Second)
}

/*
registered function to be called by timing wheel
*/
func clean(key, value any) {
	taskRunner.Schedule(func() {

		dt := value.(cacheCleaner)
		err := dt.task()
		if err == nil {
			return
		}

		next, ok := nextDelay(dt.delay)
		if ok {
			// * postpone task reschedule
			dt.delay = next
			timingWheel.SetTimer(key, dt, next)
		} else {
			// * log error and give up
			msg := fmt.Sprintf("retried but failed to clear cache with keys: %q, error: %v",
				strings.Join(dt.keys, ","), err)
			slogx.Default.Error(nil, msg)
		}
	})
}

/*
add delay to next retry
*/
func nextDelay(delay time.Duration) (time.Duration, bool) {
	switch delay {
	case time.Second:
		return time.Second * 5, true
	case time.Second * 5:
		return time.Minute, true
	case time.Minute:
		return time.Minute * 5, true
	case time.Minute * 5:
		return time.Hour, true
	default:
		return 0, false
	}
}
