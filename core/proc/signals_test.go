//go:build linux || darwin

package proc

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uc1024/f90/core/slogx"
)

func TestDone(t *testing.T) {
	buf := []byte{}
	b := bytes.NewBuffer(buf)
	slogx.Default.SetWriter(b)
	select {
	case <-Done():
		assert.Fail(t, "should run")
	default:
		slogx.Default.Error(context.Background(), "def")
	}
	t.Log(b.String())
	assert.NotNil(t, Done())
}
