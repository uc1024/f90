package ginx

import (
	"testing"
	"time"

	"github.com/uc1024/f90/core/proc"
	"github.com/uc1024/f90/core/syncx"
)

// Path: ginx/engine_test.go
func TestNewServer(t *testing.T) {
	srv := NewServer()
	go func() {
		srv.Run()
	}()
	cond := syncx.NewCond()
	cond.WaitWithTimeout(5 * time.Second)
	proc.Shutdown()
}
