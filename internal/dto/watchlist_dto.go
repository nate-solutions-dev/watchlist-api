package dto

import "time"

type CreateWatchListRequest struct {
	TMDBID    int    `json:"tmdb_id" binding:"required"`
	MediaType string `json:"media_type" binding:"required"`
	Status    string `json:"status" binding:"required"`
}

type UpdateWatchListRequest struct {
	Status         *string `json:"status"`
	Rating         *int    `json:"rating"`
	Review         *string `json:"review"`
	CurrentSeason  *int    `json:"current_season"`
	CurrentEpisode *int    `json:"current_episode"`
}

type WatchListResponse struct {
	WatchListID    string     `json:"watch_list_id"`
	UserDataID     string     `json:"user_data_id"`
	TMDBID         int        `json:"tmdb_id"`
	MediaType      string     `json:"media_type"`
	MovieTitle     string     `json:"movie_title"`
	PosterURL      string     `json:"poster_url"`
	Genres         []string   `json:"genres"`
	Status         string     `json:"status"`
	Rating         *int       `json:"rating"`
	Review         *string    `json:"review"`
	CurrentSeason  *int       `json:"current_season"`
	CurrentEpisode *int       `json:"current_episode"`
	WatchedAt      *time.Time `json:"watched_at"`
	UsrCrt         string     `json:"usr_crt"`
	UsrUpd         string     `json:"usr_upd"`
	DtmCrt         time.Time  `json:"dtm_crt"`
	DtmUpd         time.Time  `json:"dtm_upd"`
}
