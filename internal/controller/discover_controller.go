package controller

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Friel909/watchlist-api/internal/dto"
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
		c.JSON(http.StatusBadRequest, dto.Response{Message: "invalid search query", Response: http.StatusBadRequest})
		return
	}

	if !isValidDiscoverType(mediaType) {
		logger.Warn(ctx, "DiscoverController.SearchTitles", "validation failed", "reason", "invalid type", "type", mediaType)
		c.JSON(http.StatusBadRequest, dto.Response{Message: "invalid media type", Response: http.StatusBadRequest})
		return
	}

	resp, err := dc.tmdbService.SearchTitles(ctx, query, mediaType, page)
	if err != nil {
		logger.Error(ctx, "DiscoverController.SearchTitles", "search failed", "error", err.Error(), "type", mediaType, "page", page)
		c.JSON(http.StatusInternalServerError, dto.Response{Message: "search failed", Response: http.StatusInternalServerError})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Message: "success", Response: http.StatusOK, Result: resp})
}

// GetTrending returns weekly trending titles for a media type.
func (dc *DiscoverController) GetTrending(c *gin.Context) {
	ctx := c.Request.Context()

	mediaType := strings.ToLower(strings.TrimSpace(c.DefaultQuery("type", "movie")))
	page := parsePage(c.DefaultQuery("page", "1"))

	logger.Info(ctx, "DiscoverController.GetTrending", "request received", "method", c.Request.Method, "path", c.FullPath(), "type", mediaType, "page", page)

	if !isValidDiscoverType(mediaType) {
		logger.Warn(ctx, "DiscoverController.GetTrending", "validation failed", "reason", "invalid type", "type", mediaType)
		c.JSON(http.StatusBadRequest, dto.Response{Message: "invalid media type", Response: http.StatusBadRequest})
		return
	}

	resp, err := dc.tmdbService.GetTrending(ctx, mediaType, page)
	if err != nil {
		logger.Error(ctx, "DiscoverController.GetTrending", "fetch trending failed", "error", err.Error(), "type", mediaType, "page", page)
		c.JSON(http.StatusInternalServerError, dto.Response{Message: "fetch trending failed", Response: http.StatusInternalServerError})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Message: "success", Response: http.StatusOK, Result: resp})
}

// GetPopular returns currently popular titles for a media type.
func (dc *DiscoverController) GetPopular(c *gin.Context) {
	ctx := c.Request.Context()

	mediaType := strings.ToLower(strings.TrimSpace(c.DefaultQuery("type", "movie")))
	page := parsePage(c.DefaultQuery("page", "1"))

	logger.Info(ctx, "DiscoverController.GetPopular", "request received", "method", c.Request.Method, "path", c.FullPath(), "type", mediaType, "page", page)

	if !isValidDiscoverType(mediaType) {
		logger.Warn(ctx, "DiscoverController.GetPopular", "validation failed", "reason", "invalid type", "type", mediaType)
		c.JSON(http.StatusBadRequest, dto.Response{Message: "invalid media type", Response: http.StatusBadRequest})
		return
	}

	resp, err := dc.tmdbService.GetPopular(ctx, mediaType, page)
	if err != nil {
		logger.Error(ctx, "DiscoverController.GetPopular", "fetch popular failed", "error", err.Error(), "type", mediaType, "page", page)
		c.JSON(http.StatusInternalServerError, dto.Response{Message: "fetch popular failed", Response: http.StatusInternalServerError})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Message: "success", Response: http.StatusOK, Result: resp})
}

// GetTitleDetail returns full detail for a single title by TMDB ID and media type.
func (dc *DiscoverController) GetTitleDetail(c *gin.Context) {
	ctx := c.Request.Context()

	tmdbID, err := strconv.Atoi(c.Param("tmdb_id"))
	if err != nil || tmdbID < 1 {
		logger.Warn(ctx, "DiscoverController.GetTitleDetail", "validation failed", "reason", "invalid tmdb_id")
		c.JSON(http.StatusBadRequest, dto.Response{Message: "invalid tmdb_id", Response: http.StatusBadRequest})
		return
	}

	mediaType := strings.ToLower(strings.TrimSpace(c.DefaultQuery("type", "movie")))
	if !isValidDiscoverType(mediaType) {
		logger.Warn(ctx, "DiscoverController.GetTitleDetail", "validation failed", "reason", "invalid type", "type", mediaType)
		c.JSON(http.StatusBadRequest, dto.Response{Message: "invalid media type", Response: http.StatusBadRequest})
		return
	}

	logger.Info(ctx, "DiscoverController.GetTitleDetail", "request received", "tmdb_id", tmdbID, "type", mediaType)

	result, err := dc.tmdbService.GetTitleDetail(ctx, tmdbID, mediaType)
	if err != nil {
		logger.Error(ctx, "DiscoverController.GetTitleDetail", "fetch failed", "error", err.Error(), "tmdb_id", tmdbID)
		c.JSON(http.StatusInternalServerError, dto.Response{Message: "fetch title detail failed", Response: http.StatusInternalServerError})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Message: "success", Response: http.StatusOK, Result: result})
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
