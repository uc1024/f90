package timex

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCurrentTimeOut(t *testing.T) {
	offset := time.Minute
	ts := time.Now()
	ts = ts.Add(-(offset + time.Second))
	err := IsCurrentTimeWithinInterval(ts.UnixMilli(), offset)
	assert.EqualError(t, err, "time out")
	u := time.Now().Unix()
	today := u / 86400 * 86400
	t.Log(today)
}
