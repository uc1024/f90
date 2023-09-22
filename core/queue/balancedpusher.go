package queue

import (
	"context"
	"errors"
	"sync/atomic"

	"github.com/uc1024/f90/core/slogx"
)

// ErrNoAvailablePusher 表示没有可用的 pusher。
var ErrNoAvailablePusher = errors.New("no available pusher")

// BalancedPusher 用于使用轮询算法将消息推送到多个 pushe
type BalancedPusher struct {
	name    string   // pusher 的名称
	pushers []Pusher // pusher 列表
	index   uint64   // 当前 pusher 的下标
}

func NewBalancedPusher(pushers []Pusher) Pusher {
	return &BalancedPusher{
		name:    generateName(pushers),
		pushers: pushers,
	}
}

// Name 返回 pusher 的名称。
func (pusher *BalancedPusher) Name() string {
	return pusher.name
}

// Push 将消息推送到其中一个 pusher。
func (pusher *BalancedPusher) Push(message string) error {
	size := len(pusher.pushers)
	for i := 0; i < size; i++ {
		// 获取下一个 pusher 的下标
		index := atomic.AddUint64(&pusher.index, 1) % uint64(size)
		target := pusher.pushers[index]

		// 将消息推送到当前 pusher
		if err := target.Push(message); err != nil {
			slogx.Default.Error(context.Background(), err.Error())
		} else {
			return nil
		}
	}

	// 所有 pusher 都无法推送消息，返回错误
	return ErrNoAvailablePusher
}
