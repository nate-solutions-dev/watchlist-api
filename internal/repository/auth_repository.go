package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Friel909/watchlist-api/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user model.User, actorUsername string) error
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
}

type authRepository struct {
	pool *pgxpool.Pool
}

func NewAuthRepository(pool *pgxpool.Pool) AuthRepository {
	return &authRepository{pool: pool}
}

func (r *authRepository) CreateUser(ctx context.Context, user model.User, actorUsername string) error {
	query := `
		INSERT INTO USER_DATA (
			USER_DATA_ID, USER_NAME, PASSWORD, EMAIL, REGION, PREFERRED_GENRES,
			AVATAR_URL, BIO, USR_CRT, USR_UPD, DTM_CRT, DTM_UPD
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := r.pool.Exec(ctx, query,
		user.UserDataID,
		user.UserName,
		user.Password,
		user.Email,
		user.Region,
		user.PreferredGenres,
		user.AvatarURL,
		user.Bio,
		actorUsername,
		actorUsername,
		user.DtmCrt,
		user.DtmUpd,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return mapUniqueViolationToError(pgErr)
		}
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *authRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT USER_DATA_ID, USER_NAME, PASSWORD, EMAIL, REGION, PREFERRED_GENRES,
		       AVATAR_URL, BIO, USR_CRT, USR_UPD, DTM_CRT, DTM_UPD
		FROM USER_DATA
		WHERE EMAIL = $1
		LIMIT 1
	`

	var user model.User
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.UserDataID,
		&user.UserName,
		&user.Password,
		&user.Email,
		&user.Region,
		&user.PreferredGenres,
		&user.AvatarURL,
		&user.Bio,
		&user.UsrCrt,
		&user.UsrUpd,
		&user.DtmCrt,
		&user.DtmUpd,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	return &user, nil
}

func mapUniqueViolationToError(pgErr *pgconn.PgError) error {
	constraint := strings.ToLower(pgErr.ConstraintName)
	detail := strings.ToLower(pgErr.Detail)

	if strings.Contains(constraint, "user_name") || strings.Contains(detail, "(user_name)") {
		return errors.New("username already taken")
	}
	if strings.Contains(constraint, "email") || strings.Contains(detail, "(email)") {
		return errors.New("email already registered")
	}

	return errors.New("duplicate user data")
}
