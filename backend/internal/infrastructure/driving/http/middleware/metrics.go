package middleware

import (
	"github.com/gin-gonic/gin"
)

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Placeholder for metrics collection (e.g., Prometheus)
		// In a real implementation, we would record request duration,
		// status codes, and request counts here.

		c.Next()
	}
}
