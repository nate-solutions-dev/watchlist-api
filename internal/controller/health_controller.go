package controller

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthController exposes dependency-free liveness endpoints.
type HealthController struct{}

// NewHealthController constructs a HealthController.
func NewHealthController() *HealthController {
	return &HealthController{}
}

// Health reports that the HTTP server is up (no DB or external checks).
func (hc *HealthController) Health(c *gin.Context) {
	slog.Info("health check")

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"status":    "ok",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"service":   "watchlist-api",
		},
		"error": nil,
	})
}
