package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/Friel909/watchlist-api/config"
	"github.com/Friel909/watchlist-api/internal/dto"
	"github.com/Friel909/watchlist-api/internal/model"
	"github.com/Friel909/watchlist-api/internal/repository"
	"github.com/Friel909/watchlist-api/internal/validator"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) error
	Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error)
}

type authService struct {
	authRepo repository.AuthRepository
	cfg      *config.Config
	tmdb     TMDBService
}

func NewAuthService(authRepo repository.AuthRepository, cfg *config.Config, tmdb TMDBService) AuthService {
	return &authService{authRepo: authRepo, cfg: cfg, tmdb: tmdb}
}

func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) error {
	existing, err := s.authRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("email already registered")
	}

	if err := validatePassword(req.Password); err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	normalizedRegion := validator.NormalizeRegion(req.Region)
	now := time.Now().UTC()
	user := model.User{
		UserDataID:      uuid.NewString(),
		UserName:        req.Username,
		Password:        string(hashedPassword),
		Email:           req.Email,
		Region:          normalizedRegion,
		PreferredGenres: []string{},
		AvatarURL:       "",
		Bio:             "",
		UsrCrt:          req.Username,
		UsrUpd:          req.Username,
		DtmCrt:          now,
		DtmUpd:          now,
	}

	return s.authRepo.CreateUser(ctx, user, req.Username)
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.authRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	tmdbSessionID, err := s.tmdb.GetSessionToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("get tmdb session token: %w", err)
	}

	expiresAt := time.Now().UTC().Add(24 * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"USER_DATA_ID":    user.UserDataID,
		"USER_NAME":       user.UserName,
		"TMDB_SESSION_ID": tmdbSessionID,
		"exp":             expiresAt.Unix(),
	})

	signedToken, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return nil, fmt.Errorf("sign jwt token: %w", err)
	}

	return &dto.AuthResponse{
		Token:         signedToken,
		UserDataID:    user.UserDataID,
		Username:      user.UserName,
		Email:         user.Email,
		TMDBSessionID: tmdbSessionID,
	}, nil
}

func validatePassword(password string) error {
	if utf8.RuneCountInString(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	var hasUpper bool
	var hasNumber bool
	var hasSpecial bool

	for _, r := range password {
		if unicode.IsUpper(r) {
			hasUpper = true
		}
		if unicode.IsDigit(r) {
			hasNumber = true
		}
		if unicode.IsPunct(r) || unicode.IsSymbol(r) {
			hasSpecial = true
		}
	}

	if !hasUpper {
		return errors.New("password must contain at least 1 uppercase letter")
	}
	if !hasNumber {
		return errors.New("password must contain at least 1 number")
	}
	if !hasSpecial {
		return errors.New("password must contain at least 1 special character")
	}

	// Block all-whitespace edge cases that can still pass punctuation checks in some unicode sets.
	if strings.TrimSpace(password) == "" {
		return errors.New("password cannot be blank")
	}

	return nil
}
