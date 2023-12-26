package ginx

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/uc1024/f90/core/proc"
	"github.com/uc1024/f90/core/syncx"
	"go.opentelemetry.io/otel/trace"
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

func TestNewServerTest(t *testing.T) {
	srv := NewServer()
	srv.Engine.PATCH("/test/:id", func(c *gin.Context) {
		span := trace.SpanFromContext(c.Request.Context())
		// 获取TraceID
		traceID := span.SpanContext().TraceID()
		// 将TraceID转换为字符串并打印
		traceIDStr := traceID.String()
		t.Logf("TraceID: %s\n", traceIDStr)
		user := &struct {
			Id   int    `uri:"id" `
			Name string `form:"name"`
			Age  int    `json:"age"`
		}{}
		err := c.BindUri(user)
		assert.NoError(t, err)
		err = c.BindQuery(user)
		assert.NoError(t, err)
		err = c.Bind(user)
		assert.NoError(t, err)
		t.Log(user)
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	go func() {
		srv.Run()
	}()
	time.Sleep(1 * time.Second)
	// http call /test/:id
	httpUrl, err := url.Parse("http://localhost:10122/test/10086?name=coderx")
	assert.NoError(t, err)
	method := "PATCH"
	payload := strings.NewReader(`{"age": 18}`)
	client := &http.Client{}
	t.Log(httpUrl.String())
	req, err := http.NewRequest(method, httpUrl.String(), payload)
	assert.NoError(t, err)
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	assert.NoError(t, err)
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	t.Log(string(body))
	cond := syncx.NewCond()
	cond.WaitWithTimeout(5 * time.Second)
	proc.Shutdown()
}
