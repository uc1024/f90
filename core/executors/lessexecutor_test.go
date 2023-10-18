package executors

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/uc1024/f90/core/timex"
)

func TestLessExecutor_DoOrDiscard(t *testing.T) {
	executor := NewLessExecutor(time.Minute)
	// 执行
	assert.True(t, executor.DoOrDiscard(func() {}))
	// 丢弃因为在间隔内
	assert.False(t, executor.DoOrDiscard(func() {}))
	// 更新执行时间
	executor.lastTime.Set(timex.Now() - time.Minute - time.Second*30)
	// 成功执行
	assert.True(t, executor.DoOrDiscard(func() {}))
	// 执行被丢弃
	assert.False(t, executor.DoOrDiscard(func() {}))
}

func BenchmarkLessExecutor(b *testing.B) {
	exec := NewLessExecutor(time.Millisecond)
	for i := 0; i < b.N; i++ {
		exec.DoOrDiscard(func() {
		})
	}
}
