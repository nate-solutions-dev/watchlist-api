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
	BackdropPath string  `json:"backdrop_path"`
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
	BackdropURL string  `json:"backdrop_url"`
	Year        string  `json:"year"`
	Overview    string  `json:"overview"`
	VoteAverage float64 `json:"vote_average"`
}

type TMDBListResponse struct {
	Page       int           `json:"page"`
	TotalPages int           `json:"total_pages"`
	Results    []TitleResult `json:"results"`
}

// ─── Title detail DTOs ────────────────────────────────────────────────────────

type TMDBCastMember struct {
	Name        string `json:"name"`
	Character   string `json:"character"`
	ProfilePath string `json:"profile_path"`
	Order       int    `json:"order"`
}

type TMDBCrewMember struct {
	Name string `json:"name"`
	Job  string `json:"job"`
}

type TMDBCredits struct {
	Cast []TMDBCastMember `json:"cast"`
	Crew []TMDBCrewMember `json:"crew"`
}

type TMDBCreatedBy struct {
	Name string `json:"name"`
}

type TMDBFullDetailResponse struct {
	ID               int            `json:"id"`
	Title            string         `json:"title"`
	Name             string         `json:"name"`
	Overview         string         `json:"overview"`
	PosterPath       string         `json:"poster_path"`
	BackdropPath     string         `json:"backdrop_path"`
	ReleaseDate      string         `json:"release_date"`
	FirstAirDate     string         `json:"first_air_date"`
	Runtime          int            `json:"runtime"`
	EpisodeRunTime   []int          `json:"episode_run_time"`
	VoteAverage      float64        `json:"vote_average"`
	Genres           []TMDBGenre    `json:"genres"`
	Status           string         `json:"status"`
	NumberOfSeasons  int            `json:"number_of_seasons"`
	NumberOfEpisodes int            `json:"number_of_episodes"`
	CreatedBy        []TMDBCreatedBy `json:"created_by"`
	Credits          TMDBCredits    `json:"credits"`
}

type CastMember struct {
	Name       string `json:"name"`
	Character  string `json:"character"`
	ProfileURL string `json:"profile_url"`
}

type TitleDetailResult struct {
	TMDBID      int          `json:"tmdb_id"`
	MediaType   string       `json:"media_type"`
	Title       string       `json:"title"`
	Overview    string       `json:"overview"`
	PosterURL   string       `json:"poster_url"`
	BackdropURL string       `json:"backdrop_url"`
	Year        string       `json:"year"`
	Runtime     int          `json:"runtime"`
	Seasons     int          `json:"seasons"`
	Episodes    int          `json:"episodes"`
	VoteAverage float64      `json:"vote_average"`
	Genres      []string     `json:"genres"`
	Status      string       `json:"status"`
	Cast        []CastMember `json:"cast"`
	Directors   []string     `json:"directors"`
	Creators    []string     `json:"creators"`
}
