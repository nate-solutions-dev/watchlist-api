package middleware

import (
	"github.com/Friel909/watchlist-api/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestID injects a per-request UUID into gin and context.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.NewString()
		c.Set("request_id", requestID)
		c.Writer.Header().Set("X-Request-Id", requestID)

		ctx := logger.WithRequestID(c.Request.Context(), requestID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
