package middleware

import (
	"net/http"
	"runtime/debug"
	"zentinel/internal/infrastructure/driving/http/dto/response"
	"zentinel/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("stack", string(debug.Stack())),
				)

				c.AbortWithStatusJSON(http.StatusInternalServerError, response.NewErrorResponse(
					"INTERNAL_ERROR",
					"Internal Server Error",
					"An unexpected error occurred",
				))
			}
		}()
		c.Next()
	}
}
