package syncx

import (
	"time"

)

// 阻塞到等待通知
type Cond struct {
	signal chan struct{}
}

func NewCond() *Cond {
	return &Cond{
		signal: make(chan struct{}),
	}
}

// false : timer return
// true  : <-c.signal
func (c *Cond) WaitWithTimeout(timeout time.Duration) (time.Duration, bool) {
	timer := time.NewTicker(timeout)
	defer timer.Stop()

	begin := time.Now()

	select {
	case <-c.signal:
		end := time.Now().Sub(begin)
		return timeout - end, true
	case <-timer.C:
		return 0, false
	}
}

// wait signal
func (c *Cond) Wait() {
	<-c.signal
}

// send signal
func (c *Cond) Signal() {
	select {
	case c.signal <- struct{}{}:
	default:
	}
}
