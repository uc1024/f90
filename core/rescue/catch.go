package rescue

import (
	"context"
	"fmt"

	"github.com/uc1024/f90/core/slogx"
)

// 捕获
func Catch(fns ...func()) {
	for _, fn := range fns {
		fn()
	}

	if err := recover(); err != nil {
		slogx.Default.Error(context.Background(), fmt.Sprintf("catch: %v", err))
	}
}

// * 错误捕获
func CatchError(fn func(), err_fns ...func(interface{})) {
	fn()
	if err := recover(); err != nil {
		slogx.Default.Error(context.Background(), fmt.Sprintf("catch: %v", err))
		for _, fn := range err_fns {
			fn(err)
		}
	}
}
