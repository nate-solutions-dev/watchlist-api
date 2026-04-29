package dto

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20,alphanum"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Region   string `json:"region" binding:"omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token         string `json:"token"`
	UserDataID    string `json:"user_data_id"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	TMDBSessionID string `json:"tmdb_session_id"`
}

type MeResponse struct {
	UserDataID string `json:"user_data_id"`
	Username   string `json:"username"`
}
