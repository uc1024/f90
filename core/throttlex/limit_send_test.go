package throttlex

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/uc1024/f90/core/throttlex/script"
)

func TestSmsLimiter(t *testing.T) {
	// 创建一个 MiniRedis 服务器
	mr, err := miniredis.Run()
	assert.NoError(t, err)
	defer mr.Close()

	// 创建 RedisMock 客户端
	// client,_ := redismock.NewClientMock()
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	// * 加载脚本
	script_loaded := script.GetLimitSendScript(context.Background(), client)
	_ = script_loaded

	key := "sms:13800138000"
	ex, err := NewLimitSender(client, func(slo *SendLimitOptions) {
		slo.Count = 1
		slo.Period = time.Second
		slo.KeyPrefix = "limit_send"
	})
	for i := 0; i < 12; i++ {
		time.Sleep(1 * time.Second)
		a, err := ex.SendLimit(key)
		assert.NoError(t, err)
		t.Logf("%+v", a)
	}
}
