package repository

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"context"

	"github.com/FACorreiaa/glasses-management-platform/app/models"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

func (a *AccountRepository) fetchUsers(ctx context.Context, query string, args ...interface{}) ([]models.UserSession, error) {
	var us []models.UserSession

	rows, err := a.pgpool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u models.UserSession
		err := rows.Scan(
			&u.ID,
			&u.Username,
			&u.Email,
			&u.PasswordHash,
			&u.Role,
			&u.UpdatedAt,
			&u.CreatedAt,
		)

		if err != nil {
			slog.Error("Error scanning users", "err", err)
			return nil, errors.New("internal server error")
		}
		us = append(us, u)
	}

	if err := rows.Err(); err != nil {
		slog.Error("error fetching users", "err", err)
		return nil, errors.New("internal server error")
	}

	slog.Info("Users fetched", "users", us)
	return us, nil
}

func (a *AccountRepository) GetUsers(ctx context.Context, page, pageSize int, orderBy, sortBy, email string) ([]models.UserSession, error) {
	query := `
		SELECT
			user_id,
			username,
			email,
			password_hash,
			role,
			updated_at,
			created_at
		FROM
			"user" u
		WHERE u.role == 'employee'
		AND Trim(Upper(u.email)) ILIKE trim(upper('%' || $4 || '%'))
		ORDER BY
		CASE
					WHEN $1 = 'Username' THEN u.username
					ELSE u.email
				END ` + sortBy + `
			    OFFSET $2 LIMIT $3`
	offset := (page - 1) * pageSize
	slog.Info("Fetching users", "page", page, "pageSize", pageSize, "offset", offset)
	return a.fetchUsers(ctx, query, orderBy, offset, pageSize, email)
}

func (r *GlassesRepository) GetUsersByID(ctx context.Context, userID uuid.UUID) (*models.UserSession, error) {
	query := `SELECT
					user_id,
					username,
					email,
					password_hash,
					role,
					updated_at,
					created_at
				FROM
					"user" u
				WHERE u.role == 'employee' AND user_id = $1`
	var u models.UserSession

	err := r.pgpool.QueryRow(ctx, query, userID).Scan(
		&u.ID, &u.Username, &u.Email, &u.UpdatedAt, &u.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		slog.Error("Error getting user", "err", err)
		return nil, errors.New("internal server error")
	}

	slog.Info("Found user", "user_id", userID)
	return &u, nil
}

func (a *AccountRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM "user" u WHERE u.role == 'employee' AND u.user_id = $1`
	_, err := a.pgpool.Exec(ctx, query, userID)
	slog.Info("Deleted user", "user_id", userID)
	return err
}

func (a *AccountRepository) UpdateUser(ctx context.Context, form models.UpdateUserForm) error {
	// Validate the input form
	if err := a.validator.Struct(form); err != nil {
		slog.Warn("Validation error")
		return err
	}

	var passwordHash []byte
	var err error
	// Hash the new password if provided
	if form.Password != "" {
		passwordHash, err = bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
		if err != nil {
			slog.Error("Error hashing password", "err", err)
			return errors.New("internal server error")
		}
	}

	// Start the database transaction
	err = pgx.BeginFunc(ctx, a.pgpool, func(tx pgx.Tx) error {
		// Construct the query and arguments
		query := `
			UPDATE "user"
			SET username = $1, email = $2, bio = $3, image = $4, role = $5, updated_at = NOW()
		`
		args := []interface{}{form.Username, form.Email, form.Bio, form.Image, form.Role, form.UserID}

		// If the password was provided, update it too
		if form.Password != "" {
			query += `, password_hash = $6`
			args = append(args, passwordHash)
		}

		// Add the WHERE clause
		query += ` WHERE user_id = $7`
		args = append(args, form.UserID)

		// Execute the update query
		_, err := tx.Exec(ctx, query, args...)
		if err != nil {
			return errors.New("error updating user")
		}

		return nil
	})

	if err != nil {
		slog.Error("Error updating user", "err", err)
		return errors.New("internal server error")
	}

	slog.Info("Updated user", "user_id", form.UserID)
	return nil
}

func (a *AccountRepository) RegisterNewAccount(ctx context.Context, form models.RegisterForm) (*Token, error) {
	if err := a.validator.Struct(form); err != nil {
		slog.Warn("Validation error")
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Error hashing password", "err", err)
		return nil, errors.New("internal server error")
	}

	var user models.UserSession
	var token Token

	err = pgx.BeginFunc(ctx, a.pgpool, func(tx pgx.Tx) error {
		row, _ := tx.Query(
			ctx,
			`
			insert into "user" (username, email, password_hash)
				values ($1, $2, $3)
			returning
				user_id,
				username,
				email,
				password_hash,
				bio,
				role,
				image,
				created_at,
				updated_at
			`,
			form.Username,
			form.Email,
			passwordHash,
		)
		user, err = pgx.CollectOneRow(row, pgx.RowToStructByPos[models.UserSession])
		if err != nil {
			return errors.New("error inserting user")
		}

		tokenBytes := make([]byte, RandSize)
		if _, err = rand.Read(tokenBytes); err != nil {
			return errors.New("error generating token")
		}
		token = fmt.Sprintf("%x", tokenBytes)

		if err := a.redisClient.Set(ctx, token, user.ID, time.Hour*24*7).Err(); err != nil {
			return errors.New("error inserting token into Redis")
		}

		return nil
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, errors.New("username or email already taken")
		}

		slog.Error("Error creating account", "err", err)
		return nil, errors.New("internal server error")
	}

	slog.Info("Created account", "user_id", user.ID)
	return &token, nil
}
