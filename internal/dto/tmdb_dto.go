package dto

type TMDBGenre struct {
	Name string `json:"name"`
}

type TMDBDetailResponse struct {
	Title      string      `json:"title"`
	Name       string      `json:"name"`
	PosterPath string      `json:"poster_path"`
	Genres     []TMDBGenre `json:"genres"`
}

type TMDBRequestTokenResponse struct {
	RequestToken string `json:"request_token"`
	Success      bool   `json:"success"`
}

type TMDBSessionResponse struct {
	SessionID string `json:"session_id"`
	Success   bool   `json:"success"`
}

type TMDBErrorResponse struct {
	StatusCode    int    `json:"status_code"`
	StatusMessage string `json:"status_message"`
	Success       bool   `json:"success"`
}

type TMDBListItem struct {
	ID           int     `json:"id"`
	MediaType    string  `json:"media_type"`
	Title        string  `json:"title"`
	Name         string  `json:"name"`
	PosterPath   string  `json:"poster_path"`
	ReleaseDate  string  `json:"release_date"`
	FirstAirDate string  `json:"first_air_date"`
	Overview     string  `json:"overview"`
	VoteAverage  float64 `json:"vote_average"`
}

type TMDBListPayload struct {
	Page       int            `json:"page"`
	TotalPages int            `json:"total_pages"`
	Results    []TMDBListItem `json:"results"`
}

type TitleResult struct {
	TMDBID      int     `json:"tmdb_id"`
	MediaType   string  `json:"media_type"`
	Title       string  `json:"title"`
	PosterURL   string  `json:"poster_url"`
	Year        string  `json:"year"`
	Overview    string  `json:"overview"`
	VoteAverage float64 `json:"vote_average"`
}

type TMDBListResponse struct {
	Page       int           `json:"page"`
	TotalPages int           `json:"total_pages"`
	Results    []TitleResult `json:"results"`
}
