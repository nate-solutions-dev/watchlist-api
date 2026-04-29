package model

import "time"

type WatchListItem struct {
	WatchListID    string     `json:"watch_list_id" db:"WATCH_LIST_ID"`
	UserDataID     string     `json:"user_data_id" db:"USER_DATA_ID"`
	TMDBID         int        `json:"tmdb_id" db:"TMDB_ID"`
	MediaType      string     `json:"media_type" db:"MEDIA_TYPE"`
	MovieTitle     string     `json:"movie_title" db:"MOVIE_TITLE"`
	PosterURL      string     `json:"poster_url" db:"POSTER_URL"`
	Genres         []string   `json:"genres" db:"GENRES"`
	Status         string     `json:"status" db:"STATUS"`
	Rating         *int       `json:"rating" db:"RATING"`
	Review         *string    `json:"review" db:"REVIEW"`
	CurrentSeason  *int       `json:"current_season" db:"CURRENT_SEASON"`
	CurrentEpisode *int       `json:"current_episode" db:"CURRENT_EPISODE"`
	WatchedAt      *time.Time `json:"watched_at" db:"WATCHED_AT"`
	UsrCrt         string     `json:"usr_crt" db:"USR_CRT"`
	UsrUpd         string     `json:"usr_upd" db:"USR_UPD"`
	DtmCrt         time.Time  `json:"dtm_crt" db:"DTM_CRT"`
	DtmUpd         time.Time  `json:"dtm_upd" db:"DTM_UPD"`
}
