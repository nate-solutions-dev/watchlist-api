package controller

import (
	"errors"
	"net/http"

	"github.com/Friel909/watchlist-api/internal/dto"
	"github.com/Friel909/watchlist-api/internal/logger"
	"github.com/Friel909/watchlist-api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type WatchListController struct {
	watchListService service.WatchListService
}

// NewWatchListController creates a watchlist controller with injected service.
func NewWatchListController(watchListService service.WatchListService) *WatchListController {
	return &WatchListController{watchListService: watchListService}
}

// @Summary      Get watchlist items
// @Description  Return all watchlist entries for the authenticated user
// @Tags         watchlist
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Success      200            {object}  dto.Response{result=[]dto.WatchListResponse}
// @Failure      401            {object}  dto.Response
// @Failure      500            {object}  dto.Response
// @Router       /private/watchlist [get]
// GetAll returns all watchlist items for the authenticated user.
func (wc *WatchListController) GetAll(c *gin.Context) {
	callerID := c.GetString("caller_id")
	items, err := wc.watchListService.GetAllByUserID(c.Request.Context(), callerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Message:  err.Error(),
			Response: http.StatusInternalServerError,
			Result:   nil,
		})
		return
	}
	c.JSON(http.StatusOK, dto.Response{
		Message:  "success",
		Response: http.StatusOK,
		Result:   items,
	})
}

// @Summary      Create watchlist item
// @Description  Create a new watchlist entry for the authenticated user
// @Tags         watchlist
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                     true  "Bearer {token}"
// @Param        body           body      dto.CreateWatchListRequest true  "Request body"
// @Success      201            {object}  dto.Response{result=dto.WatchListResponse}
// @Failure      400            {object}  dto.Response
// @Failure      401            {object}  dto.Response
// @Router       /private/watchlist [post]
// Create creates a watchlist item using authenticated user context.
func (wc *WatchListController) Create(c *gin.Context) {
	callerID := c.GetString("caller_id")
	callerUsername := c.GetString("caller_username")
	ctx := c.Request.Context()

	var req dto.CreateWatchListRequest
	logger.Info(ctx, "WatchListController.Create", "request received", "method", c.Request.Method, "path", c.FullPath(), "caller_id", callerID)
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn(ctx, "WatchListController.Create", "request validation failed", "error", err.Error())
		c.JSON(http.StatusBadRequest, dto.Response{
			Message:  "invalid request body",
			Response: http.StatusBadRequest,
			Result:   nil,
		})
		return
	}
	logger.Debug(ctx, "WatchListController.Create", "calling service", "tmdb_id", req.TMDBID, "media_type", req.MediaType, "status", req.Status)

	item, err := wc.watchListService.Create(ctx, callerID, callerUsername, req)
	if err != nil {
		logger.Error(ctx, "WatchListController.Create", "service failed", "error", err.Error())
		c.JSON(http.StatusBadRequest, dto.Response{
			Message:  "failed to create watchlist item",
			Response: http.StatusBadRequest,
			Result:   nil,
		})
		return
	}
	logger.Info(ctx, "WatchListController.Create", "request completed", "watch_list_id", item.WatchListID)

	c.JSON(http.StatusCreated, dto.Response{
		Message:  "watchlist item created",
		Response: http.StatusCreated,
		Result:   item,
	})
}

// @Summary      Update watchlist item
// @Description  Update a watchlist entry owned by the authenticated user
// @Tags         watchlist
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                     true  "Bearer {token}"
// @Param        id             path      string                     true  "Watchlist ID"
// @Param        body           body      dto.UpdateWatchListRequest true  "Request body"
// @Success      200            {object}  dto.Response{result=dto.WatchListResponse}
// @Failure      400            {object}  dto.Response
// @Failure      401            {object}  dto.Response
// @Failure      404            {object}  dto.Response
// @Router       /private/watchlist/{id} [patch]
// Update updates one watchlist item belonging to the authenticated user.
func (wc *WatchListController) Update(c *gin.Context) {
	callerID := c.GetString("caller_id")
	callerUsername := c.GetString("caller_username")
	watchListID := c.Param("id")

	var req dto.UpdateWatchListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Message:  err.Error(),
			Response: http.StatusBadRequest,
			Result:   nil,
		})
		return
	}

	updated, err := wc.watchListService.Update(c.Request.Context(), watchListID, callerID, callerUsername, req)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, dto.Response{
				Message:  "watchlist item not found",
				Response: http.StatusNotFound,
				Result:   nil,
			})
			return
		}
		c.JSON(http.StatusBadRequest, dto.Response{
			Message:  err.Error(),
			Response: http.StatusBadRequest,
			Result:   nil,
		})
		return
	}
	if updated == nil {
		c.JSON(http.StatusNotFound, dto.Response{
			Message:  "watchlist item not found",
			Response: http.StatusNotFound,
			Result:   nil,
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Message:  "watchlist item updated",
		Response: http.StatusOK,
		Result:   updated,
	})
}

// @Summary      Delete watchlist item
// @Description  Delete a watchlist entry owned by the authenticated user
// @Tags         watchlist
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer {token}"
// @Param        id             path      string  true  "Watchlist ID"
// @Success      200            {object}  dto.Response{result=map[string]string}
// @Failure      400            {object}  dto.Response
// @Failure      401            {object}  dto.Response
// @Failure      404            {object}  dto.Response
// @Router       /private/watchlist/{id} [delete]
// Delete removes one watchlist item belonging to the authenticated user.
func (wc *WatchListController) Delete(c *gin.Context) {
	callerID := c.GetString("caller_id")
	watchListID := c.Param("id")

	err := wc.watchListService.Delete(c.Request.Context(), watchListID, callerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, dto.Response{
				Message:  "watchlist item not found",
				Response: http.StatusNotFound,
				Result:   nil,
			})
			return
		}
		c.JSON(http.StatusBadRequest, dto.Response{
			Message:  err.Error(),
			Response: http.StatusBadRequest,
			Result:   nil,
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Message:  "watchlist item deleted",
		Response: http.StatusOK,
		Result:   gin.H{"message": "deleted"},
	})
}
