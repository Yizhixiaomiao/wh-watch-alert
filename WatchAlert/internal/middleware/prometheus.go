package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
	"watchAlert/internal/metrics"
)

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()
		status := c.Writer.Status()

		metrics.RecordHTTPRequest(c.Request.Method, c.FullPath(), fmt.Sprintf("%d", status), duration)
	}
}
