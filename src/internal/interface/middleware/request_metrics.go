package middleware

import (
	"time"
	"wetalk-academy/package/logger"

	"github.com/gin-gonic/gin"
)

func RequestMetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		logger.ObserveRequest(time.Since(start), c.Writer.Status())
	}
}
