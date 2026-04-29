package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Friel909/watchlist-api/config"
	"github.com/Friel909/watchlist-api/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTMiddleware validates bearer tokens and enriches request context with caller data.
func JWTMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warn(ctx, "JWTMiddleware.Handle", "authentication failed", "reason", "missing authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{"data": nil, "error": "missing authorization header"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			logger.Warn(ctx, "JWTMiddleware.Handle", "authentication failed", "reason", "invalid authorization format")
			c.JSON(http.StatusUnauthorized, gin.H{"data": nil, "error": "invalid authorization format"})
			c.Abort()
			return
		}
		maskedToken := parts[1]
		if len(maskedToken) > 8 {
			maskedToken = maskedToken[:8] + "..."
		}

		token, err := jwt.Parse(parts[1], func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid token signing method")
			}
			return []byte(cfg.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			logger.Warn(ctx, "JWTMiddleware.Handle", "authentication failed", "reason", "invalid token", "token_prefix", maskedToken)
			c.JSON(http.StatusUnauthorized, gin.H{"data": nil, "error": "invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			logger.Warn(ctx, "JWTMiddleware.Handle", "authentication failed", "reason", "invalid token claims")
			c.JSON(http.StatusUnauthorized, gin.H{"data": nil, "error": "invalid token claims"})
			c.Abort()
			return
		}

		callerID, ok := claims["USER_DATA_ID"].(string)
		if !ok || callerID == "" {
			logger.Warn(ctx, "JWTMiddleware.Handle", "authentication failed", "reason", "missing USER_DATA_ID claim")
			c.JSON(http.StatusUnauthorized, gin.H{"data": nil, "error": "missing USER_DATA_ID claim"})
			c.Abort()
			return
		}
		callerUsername, ok := claims["USER_NAME"].(string)
		if !ok || callerUsername == "" {
			logger.Warn(ctx, "JWTMiddleware.Handle", "authentication failed", "reason", "missing USER_NAME claim")
			c.JSON(http.StatusUnauthorized, gin.H{"data": nil, "error": "missing USER_NAME claim"})
			c.Abort()
			return
		}

		c.Set("caller_id", callerID)
		c.Set("caller_username", callerUsername)
		ctx = logger.WithCallerID(c.Request.Context(), callerID)
		c.Request = c.Request.WithContext(ctx)
		if tmdbSessionID, ok := claims["TMDB_SESSION_ID"].(string); ok && tmdbSessionID != "" {
			c.Set("tmdb_session_id", tmdbSessionID)
		}

		logger.Info(c.Request.Context(), "JWTMiddleware.Handle", "authentication success", "caller_id", callerID)
		c.Next()
	}
}
