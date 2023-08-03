package threadingx

import (
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTaskRunner(t *testing.T) {
	tr := NewTaskRunner(2)

	// * 创建等待组，用于等待所有任务完成
	var wg sync.WaitGroup
	wg.Add(3)

	// * 添加 3 个任务到协程队列中
	for i := 0; i < 3; i++ {
		tr.Schedule(func() {
			defer wg.Done()
			t.Logf("Task %d is running", i)
		})
	}

	// * 等待所有任务完成
	wg.Wait()
}

func TestNewTaskRunner(t *testing.T) {
	tr := NewTaskRunner(10)
	if tr == nil {
		t.Error("NewTaskRunner returned nil")
	}
}

func TestTaskRunner_Schedule(t *testing.T) {
	tr := NewTaskRunner(1)

	var wg sync.WaitGroup
	wg.Add(2)

	// Test that the first task runs without error
	tr.Schedule(func() {
		defer wg.Done()
	})

	// Test that the second task runs without error
	tr.Schedule(func() {
		defer wg.Done()
	})

	wg.Wait()
}

func TestTaskRunner_SetErrorCall(t *testing.T) {
	tr := NewTaskRunner(1)

	// Set a custom error callback
	tr.SetErrorCall(func(e interface{}) {
		assert.Equal(t, "test error", e.(error).Error())
	})

	// Schedule a task that will panic
	tr.Schedule(func() {
		panic(errors.New("test error"))
	})
}
