package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Friel909/watchlist-api/internal/logger"
	"github.com/Friel909/watchlist-api/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WatchListRepository interface {
	Create(ctx context.Context, item model.WatchListItem, actorUsername string) (*model.WatchListItem, error)
	GetAllByUserID(ctx context.Context, userDataID string) ([]model.WatchListItem, error)
	GetByID(ctx context.Context, watchListID, userDataID string) (*model.WatchListItem, error)
	Update(ctx context.Context, watchListID, userDataID, actorUsername string, req model.WatchListItem) (*model.WatchListItem, error)
	Delete(ctx context.Context, watchListID, userDataID string) error
}

type watchListRepository struct {
	pool *pgxpool.Pool
}

// NewWatchListRepository creates a watchlist repository backed by pgx pool.
func NewWatchListRepository(pool *pgxpool.Pool) WatchListRepository {
	return &watchListRepository{pool: pool}
}

// Create inserts a new row into WATCH_LIST and returns the inserted record.
func (r *watchListRepository) Create(ctx context.Context, item model.WatchListItem, actorUsername string) (*model.WatchListItem, error) {
	logger.Info(ctx, "WatchListRepository.Create", "executing INSERT on WATCH_LIST")
	query := `
		INSERT INTO WATCH_LIST (
			WATCH_LIST_ID, USER_DATA_ID, TMDB_ID, MEDIA_TYPE, MOVIE_TITLE,
			POSTER_URL, GENRES, STATUS, RATING, REVIEW, CURRENT_SEASON,
			CURRENT_EPISODE, WATCHED_AT, USR_CRT, USR_UPD, DTM_CRT, DTM_UPD
		)
		VALUES ($1, $2, $3, $4, $5,
				$6, $7, $8, $9, $10, $11,
				$12, $13, $14, $15, $16, $17)
		RETURNING WATCH_LIST_ID, USER_DATA_ID, TMDB_ID, MEDIA_TYPE, MOVIE_TITLE,
		          POSTER_URL, GENRES, STATUS, RATING, REVIEW, CURRENT_SEASON,
		          CURRENT_EPISODE, WATCHED_AT, USR_CRT, USR_UPD, DTM_CRT, DTM_UPD
	`

	var created model.WatchListItem
	err := r.pool.QueryRow(ctx, query,
		item.WatchListID,
		item.UserDataID,
		item.TMDBID,
		item.MediaType,
		item.MovieTitle,
		item.PosterURL,
		item.Genres,
		item.Status,
		item.Rating,
		item.Review,
		item.CurrentSeason,
		item.CurrentEpisode,
		item.WatchedAt,
		actorUsername,
		actorUsername,
		item.DtmCrt,
		item.DtmUpd,
	).Scan(
		&created.WatchListID,
		&created.UserDataID,
		&created.TMDBID,
		&created.MediaType,
		&created.MovieTitle,
		&created.PosterURL,
		&created.Genres,
		&created.Status,
		&created.Rating,
		&created.Review,
		&created.CurrentSeason,
		&created.CurrentEpisode,
		&created.WatchedAt,
		&created.UsrCrt,
		&created.UsrUpd,
		&created.DtmCrt,
		&created.DtmUpd,
	)
	if err != nil {
		logger.Error(ctx, "WatchListRepository.Create", "insert failed", "error", err.Error())
		return nil, fmt.Errorf("create watchlist item: %w", err)
	}

	return &created, nil
}

// GetAllByUserID selects all WATCH_LIST rows for one user in descending creation time.
func (r *watchListRepository) GetAllByUserID(ctx context.Context, userDataID string) ([]model.WatchListItem, error) {
	logger.Info(ctx, "WatchListRepository.GetAllByUserID", "executing SELECT on WATCH_LIST")
	query := `
		SELECT WATCH_LIST_ID, USER_DATA_ID, TMDB_ID, MEDIA_TYPE, MOVIE_TITLE,
		       POSTER_URL, GENRES, STATUS, RATING, REVIEW, CURRENT_SEASON,
		       CURRENT_EPISODE, WATCHED_AT, USR_CRT, USR_UPD, DTM_CRT, DTM_UPD
		FROM WATCH_LIST
		WHERE USER_DATA_ID = $1
		ORDER BY DTM_CRT DESC
	`

	rows, err := r.pool.Query(ctx, query, userDataID)
	if err != nil {
		logger.Error(ctx, "WatchListRepository.GetAllByUserID", "select failed", "error", err.Error())
		return nil, fmt.Errorf("get all watchlist by user: %w", err)
	}
	defer rows.Close()

	items := make([]model.WatchListItem, 0)
	for rows.Next() {
		var item model.WatchListItem
		err := rows.Scan(
			&item.WatchListID,
			&item.UserDataID,
			&item.TMDBID,
			&item.MediaType,
			&item.MovieTitle,
			&item.PosterURL,
			&item.Genres,
			&item.Status,
			&item.Rating,
			&item.Review,
			&item.CurrentSeason,
			&item.CurrentEpisode,
			&item.WatchedAt,
			&item.UsrCrt,
			&item.UsrUpd,
			&item.DtmCrt,
			&item.DtmUpd,
		)
		if err != nil {
			logger.Error(ctx, "WatchListRepository.GetAllByUserID", "scan failed", "error", err.Error())
			return nil, fmt.Errorf("scan watchlist row: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		logger.Error(ctx, "WatchListRepository.GetAllByUserID", "row iteration failed", "error", err.Error())
		return nil, fmt.Errorf("iterate watchlist rows: %w", err)
	}

	return items, nil
}

// GetByID selects one WATCH_LIST row by watchlist id and user id.
func (r *watchListRepository) GetByID(ctx context.Context, watchListID, userDataID string) (*model.WatchListItem, error) {
	logger.Info(ctx, "WatchListRepository.GetByID", "executing SELECT on WATCH_LIST")
	query := `
		SELECT WATCH_LIST_ID, USER_DATA_ID, TMDB_ID, MEDIA_TYPE, MOVIE_TITLE,
		       POSTER_URL, GENRES, STATUS, RATING, REVIEW, CURRENT_SEASON,
		       CURRENT_EPISODE, WATCHED_AT, USR_CRT, USR_UPD, DTM_CRT, DTM_UPD
		FROM WATCH_LIST
		WHERE WATCH_LIST_ID = $1 AND USER_DATA_ID = $2
		LIMIT 1
	`

	var item model.WatchListItem
	err := r.pool.QueryRow(ctx, query, watchListID, userDataID).Scan(
		&item.WatchListID,
		&item.UserDataID,
		&item.TMDBID,
		&item.MediaType,
		&item.MovieTitle,
		&item.PosterURL,
		&item.Genres,
		&item.Status,
		&item.Rating,
		&item.Review,
		&item.CurrentSeason,
		&item.CurrentEpisode,
		&item.WatchedAt,
		&item.UsrCrt,
		&item.UsrUpd,
		&item.DtmCrt,
		&item.DtmUpd,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		logger.Error(ctx, "WatchListRepository.GetByID", "select failed", "error", err.Error())
		return nil, fmt.Errorf("get watchlist by id: %w", err)
	}

	return &item, nil
}

// Update updates one WATCH_LIST row filtered by watchlist id and user id.
func (r *watchListRepository) Update(ctx context.Context, watchListID, userDataID, actorUsername string, req model.WatchListItem) (*model.WatchListItem, error) {
	logger.Info(ctx, "WatchListRepository.Update", "executing UPDATE on WATCH_LIST")
	setClauses := make([]string, 0)
	args := make([]any, 0)
	argPos := 1

	if req.Status != "" {
		setClauses = append(setClauses, fmt.Sprintf("STATUS = $%d", argPos))
		args = append(args, req.Status)
		argPos++
	}
	if req.Rating != nil {
		setClauses = append(setClauses, fmt.Sprintf("RATING = $%d", argPos))
		args = append(args, req.Rating)
		argPos++
	}
	if req.Review != nil {
		setClauses = append(setClauses, fmt.Sprintf("REVIEW = $%d", argPos))
		args = append(args, req.Review)
		argPos++
	}
	if req.CurrentSeason != nil {
		setClauses = append(setClauses, fmt.Sprintf("CURRENT_SEASON = $%d", argPos))
		args = append(args, req.CurrentSeason)
		argPos++
	}
	if req.CurrentEpisode != nil {
		setClauses = append(setClauses, fmt.Sprintf("CURRENT_EPISODE = $%d", argPos))
		args = append(args, req.CurrentEpisode)
		argPos++
	}

	setClauses = append(setClauses, fmt.Sprintf("USR_UPD = $%d", argPos))
	args = append(args, actorUsername)
	argPos++

	setClauses = append(setClauses, fmt.Sprintf("DTM_UPD = $%d", argPos))
	now := time.Now().UTC()
	args = append(args, now)
	argPos++

	args = append(args, watchListID, userDataID)
	watchListIDPos := argPos
	userDataIDPos := argPos + 1

	query := fmt.Sprintf(`
		UPDATE WATCH_LIST
		SET %s
		WHERE WATCH_LIST_ID = $%d AND USER_DATA_ID = $%d
		RETURNING WATCH_LIST_ID, USER_DATA_ID, TMDB_ID, MEDIA_TYPE, MOVIE_TITLE,
		          POSTER_URL, GENRES, STATUS, RATING, REVIEW, CURRENT_SEASON,
		          CURRENT_EPISODE, WATCHED_AT, USR_UPD, DTM_CRT, DTM_UPD
	`, strings.Join(setClauses, ", "), watchListIDPos, userDataIDPos)

	var updated model.WatchListItem
	err := r.pool.QueryRow(ctx, query, args...).Scan(
		&updated.WatchListID,
		&updated.UserDataID,
		&updated.TMDBID,
		&updated.MediaType,
		&updated.MovieTitle,
		&updated.PosterURL,
		&updated.Genres,
		&updated.Status,
		&updated.Rating,
		&updated.Review,
		&updated.CurrentSeason,
		&updated.CurrentEpisode,
		&updated.WatchedAt,
		&updated.UsrUpd,
		&updated.DtmCrt,
		&updated.DtmUpd,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		logger.Error(ctx, "WatchListRepository.Update", "update failed", "error", err.Error())
		return nil, fmt.Errorf("update watchlist item: %w", err)
	}

	return &updated, nil
}

// Delete removes one WATCH_LIST row filtered by watchlist id and user id.
func (r *watchListRepository) Delete(ctx context.Context, watchListID, userDataID string) error {
	logger.Info(ctx, "WatchListRepository.Delete", "executing DELETE on WATCH_LIST")
	query := `
		DELETE FROM WATCH_LIST
		WHERE WATCH_LIST_ID = $1 AND USER_DATA_ID = $2
	`

	result, err := r.pool.Exec(ctx, query, watchListID, userDataID)
	if err != nil {
		logger.Error(ctx, "WatchListRepository.Delete", "delete failed", "error", err.Error())
		return fmt.Errorf("delete watchlist item: %w", err)
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
