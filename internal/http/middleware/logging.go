package middleware

import (
	"github.com/gin-gonic/gin"
	"go_shurtiner/pkg/logging"
	"net/http"
	"time"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		rawQuery := c.Request.URL.RawQuery

		// process request
		c.Next()

		logger := logging.FromContext(c.Request.Context())
		timestamp := time.Now()
		latency := timestamp.Sub(start)
		latencyValue := latency.String()
		clientIP := c.ClientIP()
		method := c.Request.Method
		status := c.Writer.Status()
		if rawQuery != "" {
			path = path + "?" + rawQuery
		}
		// append logger keys if not success or too slow latency.
		if status != http.StatusOK {
			logger = logger.With("status", status)
		}
		if latency > time.Second*3 {
			logger = logger.With("latency", latencyValue)
		}
		logger.Infof("[SHORT_API] %v | %3d | %s | %13v | %15s | %-7s %#v",
			timestamp.Format("2006/01/02 - 15:04:05"),
			status,
			latency,
			latencyValue,
			clientIP,
			method,
			path,
		)
	}
}
