package repository

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"time"

	"context"

	"crypto/rand"

	"github.com/FACorreiaa/glasses-management-platform/app/models"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

const (
	RedisPrefix = "user_session:"
	RandSize    = 32
	MaxAge      = time.Hour * 24 * 60
)

type Token = string

type AccountRepository struct {
	pgpool    *pgxpool.Pool
	validator *validator.Validate
	sessions  *sessions.CookieStore
}

func NewAccountRepository(db *pgxpool.Pool,
	validator *validator.Validate,
	sessions *sessions.CookieStore,
) *AccountRepository {
	return &AccountRepository{
		pgpool:    db,
		validator: validator,
		sessions:  sessions,
	}
}

// Logout deletes the user token from the Redis store.
func (a *AccountRepository) Logout(ctx context.Context, token Token) error {
	// Delete the token from PostgreSQL
	cmdTag, err := a.pgpool.Exec(
		ctx,
		`
        DELETE FROM user_sessions WHERE token = $1
        `,
		token,
	)
	if err != nil {
		return errors.New("error deleting token")
	}

	if cmdTag.RowsAffected() == 0 {
		// Token not found, consider it already logged out
		return nil
	}

	return nil
}

func (a *AccountRepository) Login(ctx context.Context, form models.LoginForm) (*Token, error) {
	if err := a.validator.Struct(form); err != nil {
		return nil, err
	}

	rows, _ := a.pgpool.Query(
		ctx,
		`
        select
            user_id,
            username,
            email,
            password_hash,
            role,
            created_at,
            updated_at
        from "user" where email = $1 limit 1
        `,
		form.Email,
	)
	user, err := pgx.CollectOneRow[models.UserSession](rows, pgx.RowToStructByPos[models.UserSession])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("invalid email or password")
		}

		slog.Error(" querying user", "err", err)
		return nil, errors.New("internal server error")
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(form.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	tokenBytes := make([]byte, RandSize)
	if _, err = rand.Read(tokenBytes); err != nil {
		slog.Error(" generating token", "err", err)
		return nil, errors.New("internal server error")
	}

	token := fmt.Sprintf("%x", tokenBytes)
	log.Printf("Generated token: %s", token)

	_, err = a.pgpool.Exec(
		ctx,
		`
        INSERT INTO user_sessions (token, user_id) VALUES ($1, $2)
        `,
		token, user.ID,
	)
	if err != nil {
		log.Println(" inserting token into PostgreSQL:", err)
		return nil, errors.New("internal server error")
	}

	log.Println("Token successfully inserted into PostgreSQL")
	return &token, nil
}

func (m *MiddlewareRepository) UserFromSessionToken(ctx context.Context, token Token) (*models.UserSession, error) {
	log.Println("Retrieving user ID from PostgreSQL for token:", token)

	// Retrieve user ID from the user_sessions table
	var userID string
	var createdAt time.Time
	err := m.Pgpool.QueryRow(
		ctx,
		`
        SELECT user_id, created_at FROM user_sessions WHERE token = $1
        `,
		token,
	).Scan(&userID, &createdAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("auth session expired")
		}

		log.Println(" querying user ID with token from PostgreSQL:", err)
		return nil, errors.New("internal server error")
	}

	// Check if the session has expired
	if time.Since(createdAt) > MaxAge {
		return nil, errors.New("auth session expired")
	}

	// Retrieve user details from PostgreSQL
	rows, err := m.Pgpool.Query(
		ctx,
		`
        select
            user_id,
            username,
            email,
            password_hash,
            role,
            created_at,
            updated_at
        from "user" where user_id = $1 limit 1
        `,
		userID,
	)
	if err != nil {
		log.Println(" querying user from PostgreSQL:", err)
		return nil, errors.New("internal server error")
	}

	userWithToken, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[models.UserSession])
	if err != nil {
		return nil, errors.New("internal server error")
	}

	return &userWithToken, nil
}
