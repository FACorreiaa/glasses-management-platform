package repository

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"context"

	"github.com/FACorreiaa/glasses-management-platform/app/models"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type AdminRepository struct {
	pgpool      *pgxpool.Pool
	redisClient *redis.Client
	validator   *validator.Validate
	sessions    *sessions.CookieStore
}

func NewAdminRepository(db *pgxpool.Pool,
	redisClient *redis.Client,
	validator *validator.Validate,
	sessions *sessions.CookieStore,
) *AdminRepository {
	return &AdminRepository{
		pgpool:      db,
		redisClient: redisClient,
		validator:   validator,
		sessions:    sessions,
	}
}

func (a *AdminRepository) fetchUsers(ctx context.Context, query string, args ...interface{}) ([]models.UserSession, error) {
	var us []models.UserSession

	rows, err := a.pgpool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u models.UserSession
		if err := rows.Scan(
			&u.ID,
			&u.Username,
			&u.Email,
			&u.Role,
			&u.UpdatedAt,
			&u.CreatedAt,
		); err != nil {
			slog.Error(" scanning users", "err", err)
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

func (a *AdminRepository) GetUsers(ctx context.Context, page, pageSize int, orderBy, sortBy, email string) ([]models.UserSession, error) {
	query := `
		SELECT
			user_id,
			username,
			email,
			role,
			updated_at,
			created_at
		FROM
			"user" u
		WHERE u.role = 'employee'
		AND Trim(Upper(u.email)) ILIKE trim(upper('%' || $4 || '%'))
		ORDER BY
		CASE
					WHEN $1 = 'Username' THEN u.username
					WHEN $1 = 'Email' THEN u.email
					WHEN $1 = 'Role' THEN u.role
					ELSE u.email
				END ` + sortBy + `
			    OFFSET $2 LIMIT $3`
	offset := (page - 1) * pageSize
	slog.Info("Fetching users", "page", page, "pageSize", pageSize, "offset", offset)
	return a.fetchUsers(ctx, query, orderBy, offset, pageSize, email)
}

func (a *AdminRepository) GetUsersByID(ctx context.Context, userID uuid.UUID) (*models.UserSession, error) {
	query := `SELECT
					user_id,
					username,
					email,
					role,
					updated_at,
					created_at
				FROM
					"user" u
				WHERE u.role = 'employee' AND user_id = $1`
	var u models.UserSession

	if err := a.pgpool.QueryRow(ctx, query, userID).Scan(
		&u.ID, &u.Username, &u.Email, &u.Role, &u.UpdatedAt, &u.CreatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		slog.Error(" getting user", "err", err)
		return nil, errors.New("internal server error")
	}

	slog.Info("Found user", "user_id", userID)
	return &u, nil
}

func (a *AdminRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM "user" u WHERE u.role = 'employee' AND u.user_id = $1`
	_, err := a.pgpool.Exec(ctx, query, userID)
	slog.Info("Deleted user", "user_id", userID)
	return err
}

func (a *AdminRepository) UpdateUser(ctx context.Context, form models.UpdateUserForm) error {
	// Check if the username already exists for a different user
	err := a.pgpool.QueryRow(ctx, `SELECT user_id
										FROM "user"
										WHERE username = $1
										AND user_id != $2`, form.Username, form.UserID).Scan(&form.UserID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		slog.Error("checking existing username", "err", err)
		return errors.New("internal server error")
	}

	// Check if the email already exists for a different user
	err = a.pgpool.QueryRow(ctx, `SELECT user_id FROM "user"
               							WHERE email = $1
               							AND user_id != $2`, form.Email, form.UserID).Scan(&form.UserID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		slog.Error("checking existing email", "err", err)
		return errors.New("internal server error")
	}

	// Hash the new password if provided
	var passwordHash []byte
	var setPasswordHash bool
	if form.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
		if err != nil {
			slog.Error("hashing password", "err", err)
			return errors.New("internal server error")
		}
		passwordHash = hashedPassword
		setPasswordHash = true
	}

	// Construct the base query
	query := `
        UPDATE "user"
        SET username = $1, email = $2, updated_at = NOW()
        WHERE user_id = $3
    `

	// Prepare arguments
	args := []interface{}{form.Username, form.Email, form.UserID}

	// Conditionally append password_hash to the query and args
	if setPasswordHash {
		query = `
            UPDATE "user"
            SET username = $1, email = $2, updated_at = NOW(), password_hash = $4
            WHERE user_id = $3
        `
		args = append(args, passwordHash)
	}

	// Log the query and arguments for debugging
	slog.Info("Executing query", "query", query, "args", args)

	// Execute the update query
	_, err = a.pgpool.Exec(ctx, query, args...)
	if err != nil {
		slog.Error("updating user", "err", err)
		return errors.New("internal server error")
	}

	slog.Info("Updated user", "user_id", form.UserID)
	return nil
}

func (a *AdminRepository) InsertUser(ctx context.Context, form models.RegisterFormValues) (*Token, error) {
	if err := a.validator.Struct(form); err != nil {
		slog.Warn("Validation error")
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error(" hashing password", "err", err)
		return nil, errors.New("internal server error")
	}

	var user models.UserSession
	var token Token

	if err = pgx.BeginFunc(ctx, a.pgpool, func(tx pgx.Tx) error {
		row, _ := tx.Query(
			ctx,
			`
			insert into "user" (username, email, password_hash, created_at, updated_at)
				values ($1, $2, $3, $4, $4)
			returning
				user_id,
				username,
				email,
				password_hash,
				role,
				created_at,
				updated_at
			`,
			form.Username,
			form.Email,
			passwordHash,
			time.Now(),
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
	}); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, errors.New("username or email already taken")
		}

		slog.Error(" creating account", "err", err)
		return nil, errors.New("internal server error")
	}

	slog.Info("Created account", "user_id", user.ID)
	return &token, nil
}

func (a *AdminRepository) GetUsersSum(ctx context.Context) (int, error) {
	var count int
	row := a.pgpool.QueryRow(ctx, `SELECT Count(DISTINCT u.user_id) FROM "user" u`)
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	slog.Info("User count", "count", count)
	return count, nil
}

func (a *AdminRepository) GetAdminID(ctx context.Context, userID uuid.UUID) (*models.UserSession, error) {
	query := `SELECT
					user_id,
					username,
					email,
					role,
					updated_at,
					created_at
				FROM
					"user" u
				WHERE u.role = 'admin' AND user_id = $1`
	var u models.UserSession

	if err := a.pgpool.QueryRow(ctx, query, userID).Scan(
		&u.ID, &u.Username, &u.Email, &u.Role, &u.UpdatedAt, &u.CreatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		slog.Error(" getting user", "err", err)
		return nil, errors.New("internal server error")
	}

	slog.Info("Found user", "user_id", userID)
	return &u, nil
}

func (r *GlassesRepository) GetGlassesDetails(ctx context.Context, page, pageSize int, orderBy, sortBy, reference string, leftEye, rightEye *float64) ([]models.Glasses, error) {
	query := `
		SELECT
			u.username, u.email, g.left_eye_strength, g.right_eye_strength,
			g.reference, g.is_in_stock,
			COALESCE(g.updated_at, '1970-01-01 00:00:00') AS updated_at, g.created_at
		FROM glasses g
		JOIN "user" u ON g.user_id = u.user_id
		WHERE
			Trim(Upper(g.reference)) ILIKE trim(upper('%' || $4 || '%'))
			AND ($5::float8 IS NULL OR g.left_eye_strength = $5)
			AND ($6::float8 IS NULL OR g.right_eye_strength = $6)
		ORDER BY
			CASE
				WHEN $1 = 'Brand' THEN g.brand
				WHEN $1 = 'Color' THEN g.color
				WHEN $1 = 'Reference' THEN g.reference
				WHEN $1 = 'Type' THEN g.type
				WHEN $1 = 'Features' THEN g.features
				ELSE g.brand
			END ` + sortBy + `
		OFFSET $2 LIMIT $3`
	offset := (page - 1) * pageSize
	slog.Info("Glasses fetched", "offset", offset)
	return r.fetchGlassesDetails(ctx, query, orderBy, offset, pageSize, reference, leftEye, rightEye)
}

func (a *AdminRepository) GetEmail(ctx context.Context, email string) error {
	var retrievedEmail string
	row := a.pgpool.QueryRow(ctx, `SELECT email FROM "user" WHERE email = $1`, email)
	if err := row.Scan(&retrievedEmail); err != nil {
		return err
	}
	return nil
}
