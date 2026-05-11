package controller

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Friel909/watchlist-api/internal/logger"
	"github.com/Friel909/watchlist-api/internal/service"
	"github.com/gin-gonic/gin"
)

type DiscoverController struct {
	tmdbService service.TMDBService
}

// NewDiscoverController creates a discover controller backed by TMDB service.
func NewDiscoverController(tmdbService service.TMDBService) *DiscoverController {
	return &DiscoverController{tmdbService: tmdbService}
}

// SearchTitles searches movie or show titles based on query parameters.
func (dc *DiscoverController) SearchTitles(c *gin.Context) {
	ctx := c.Request.Context()

	query := strings.TrimSpace(c.Query("q"))
	mediaType := strings.ToLower(strings.TrimSpace(c.DefaultQuery("type", "movie")))
	page := parsePage(c.DefaultQuery("page", "1"))

	logger.Info(ctx, "DiscoverController.SearchTitles", "request received", "method", c.Request.Method, "path", c.FullPath(), "type", mediaType, "page", page)

	if query == "" {
		logger.Warn(ctx, "DiscoverController.SearchTitles", "validation failed", "reason", "empty query")
		c.JSON(http.StatusBadRequest, gin.H{"data": nil, "error": "invalid search query"})
		return
	}

	if !isValidDiscoverType(mediaType) {
		logger.Warn(ctx, "DiscoverController.SearchTitles", "validation failed", "reason", "invalid type", "type", mediaType)
		c.JSON(http.StatusBadRequest, gin.H{"data": nil, "error": "invalid media type"})
		return
	}

	resp, err := dc.tmdbService.SearchTitles(ctx, query, mediaType, page)
	if err != nil {
		logger.Error(ctx, "DiscoverController.SearchTitles", "search failed", "error", err.Error(), "type", mediaType, "page", page)
		c.JSON(http.StatusInternalServerError, gin.H{"data": nil, "error": "search failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp, "error": nil})
}

// GetTrending returns weekly trending titles for a media type.
func (dc *DiscoverController) GetTrending(c *gin.Context) {
	ctx := c.Request.Context()

	mediaType := strings.ToLower(strings.TrimSpace(c.DefaultQuery("type", "movie")))
	page := parsePage(c.DefaultQuery("page", "1"))

	logger.Info(ctx, "DiscoverController.GetTrending", "request received", "method", c.Request.Method, "path", c.FullPath(), "type", mediaType, "page", page)

	if !isValidDiscoverType(mediaType) {
		logger.Warn(ctx, "DiscoverController.GetTrending", "validation failed", "reason", "invalid type", "type", mediaType)
		c.JSON(http.StatusBadRequest, gin.H{"data": nil, "error": "invalid media type"})
		return
	}

	resp, err := dc.tmdbService.GetTrending(ctx, mediaType, page)
	if err != nil {
		logger.Error(ctx, "DiscoverController.GetTrending", "fetch trending failed", "error", err.Error(), "type", mediaType, "page", page)
		c.JSON(http.StatusInternalServerError, gin.H{"data": nil, "error": "search failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp, "error": nil})
}

// isValidDiscoverType checks whether media type is supported by discover handlers.
func isValidDiscoverType(mediaType string) bool {
	return mediaType == "movie" || mediaType == "show"
}

// parsePage parses and normalizes page query param with default fallback.
func parsePage(pageRaw string) int {
	page, err := strconv.Atoi(pageRaw)
	if err != nil || page < 1 {
		return 1
	}
	return page
}
