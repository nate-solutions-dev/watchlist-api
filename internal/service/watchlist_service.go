package service

import (
	"context"
	"time"

	"github.com/Friel909/watchlist-api/internal/dto"
	"github.com/Friel909/watchlist-api/internal/logger"
	"github.com/Friel909/watchlist-api/internal/model"
	"github.com/Friel909/watchlist-api/internal/repository"
	"github.com/google/uuid"
)

type WatchListService interface {
	Create(ctx context.Context, userDataID, callerUsername string, req dto.CreateWatchListRequest) (*model.WatchListItem, error)
	GetAllByUserID(ctx context.Context, userDataID string) ([]model.WatchListItem, error)
	Update(ctx context.Context, watchListID, userDataID, callerUsername string, req dto.UpdateWatchListRequest) (*model.WatchListItem, error)
	Delete(ctx context.Context, watchListID, userDataID string) error
}

type watchListService struct {
	repo repository.WatchListRepository
	tmdb TMDBService
}

// NewWatchListService creates a watchlist service with repository and TMDB dependencies.
func NewWatchListService(repo repository.WatchListRepository, tmdb TMDBService) WatchListService {
	return &watchListService{repo: repo, tmdb: tmdb}
}

// Create validates metadata source and persists a new watchlist item.
func (s *watchListService) Create(ctx context.Context, userDataID, callerUsername string, req dto.CreateWatchListRequest) (*model.WatchListItem, error) {
	logger.Info(ctx, "WatchListService.Create", "entry", "tmdb_id", req.TMDBID, "media_type", req.MediaType, "caller_id", userDataID)
	logger.Debug(ctx, "WatchListService.Create", "calling external TMDB metadata")
	title, posterURL, genres, err := s.tmdb.FetchMetadata(ctx, req.TMDBID, req.MediaType)
	if err != nil {
		logger.Error(ctx, "WatchListService.Create", "tmdb metadata fetch failed", "error", err.Error())
		return nil, err
	}

	now := time.Now().UTC()
	item := model.WatchListItem{
		WatchListID: uuid.NewString(),
		UserDataID:  userDataID,
		TMDBID:      req.TMDBID,
		MediaType:   req.MediaType,
		MovieTitle:  title,
		PosterURL:   posterURL,
		Genres:      genres,
		Status:      req.Status,
		DtmCrt:      now,
		DtmUpd:      now,
	}

	logger.Debug(ctx, "WatchListService.Create", "calling repo create")
	created, err := s.repo.Create(ctx, item, callerUsername)
	if err != nil {
		logger.Error(ctx, "WatchListService.Create", "repo create failed", "error", err.Error())
		return nil, err
	}
	logger.Info(ctx, "WatchListService.Create", "success", "watch_list_id", created.WatchListID)
	return created, nil
}

// GetAllByUserID fetches all watchlist items for one user id.
func (s *watchListService) GetAllByUserID(ctx context.Context, userDataID string) ([]model.WatchListItem, error) {
	return s.repo.GetAllByUserID(ctx, userDataID)
}

// Update applies partial field changes to one watchlist item.
func (s *watchListService) Update(ctx context.Context, watchListID, userDataID, callerUsername string, req dto.UpdateWatchListRequest) (*model.WatchListItem, error) {
	payload := model.WatchListItem{
		Rating:         req.Rating,
		Review:         req.Review,
		CurrentSeason:  req.CurrentSeason,
		CurrentEpisode: req.CurrentEpisode,
	}
	if req.Status != nil {
		payload.Status = *req.Status
	}

	return s.repo.Update(ctx, watchListID, userDataID, callerUsername, payload)
}

// Delete removes one watchlist item for the given user id.
func (s *watchListService) Delete(ctx context.Context, watchListID, userDataID string) error {
	return s.repo.Delete(ctx, watchListID, userDataID)
}
