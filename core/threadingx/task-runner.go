package threadingx

/*
	go携程管理器
*/
import (
	"context"
	"fmt"

	"github.com/uc1024/f90/core/rescue"
	"github.com/uc1024/f90/core/slogx"
)

type (
	TaskRunnerErrorCall func(interface{})

	// TaskRunner 类用于实现 Go 协程管理器
	TaskRunner struct {
		limitChan chan struct{}       // * 用于限制并发数量的 channel
		errorCall TaskRunnerErrorCall // * 错误回调函数
	}
)

/*
NewTaskRunner 方法创建并返回 TaskRunner 实例
concurrent 参数用于指定并发数量
*/
func NewTaskRunner(concurrent int) *TaskRunner {
	tr := TaskRunner{}
	// * 初始化并发数量限制 channel
	tr.limitChan = make(chan struct{}, concurrent)
	// * 初始化错误回调函数为默认函数
	tr.errorCall = tr.onError
	return &tr
}

// * Schedule 方法用于添加一个任务到协程队列中
func (tr *TaskRunner) Schedule(task func()) {

	// * 往 channel 中添加一个值，如果 channel 已满则阻塞等待
	tr.limitChan <- struct{}{}

	go func() {
		defer rescue.CatchError(
			tr.done,      // * 默认执行
			tr.errorCall, // * 错误时执行
		)
		task()
	}()
}

// * done 方法从 channel 中移除一个值
func (tr *TaskRunner) done() {
	<-tr.limitChan
}

// * onError 方法是默认的错误回调函数，将错误信息打印到控制台
func (tr *TaskRunner) onError(e interface{}) {
	slogx.Default.Error(context.Background(),
		fmt.Sprintf("TaskRunner error: %v", e))
}

// * SetErrorCall 方法用于设置错误回调函数
func (tr *TaskRunner) SetErrorCall(t TaskRunnerErrorCall) {
	tr.errorCall = t
}
