package trace

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

func SetHeaderTrace(c *gin.Context) {
	c.Header("X-Trace-Id", getTrace(c.Request.Context()))
	c.Next()
}

func getTrace(ctx context.Context) string {
	if sc := trace.SpanContextFromContext(ctx); sc.HasTraceID() {
		return sc.TraceID().String()
	}
	return ""
}
