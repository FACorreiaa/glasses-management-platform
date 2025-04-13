package repository

import (
	"errors"
	"log/slog"

	"github.com/FACorreiaa/glasses-management-platform/app/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"

	"context"
)

var tracer = otel.Tracer("github.com/FACorreiaa/glasses-management-platform/services")

type GlassesRepository struct {
	pgpool *pgxpool.Pool
}

func NewGlassesRepository(db *pgxpool.Pool) *GlassesRepository {
	return &GlassesRepository{pgpool: db}
}

func (r *GlassesRepository) fetchGlasses(ctx context.Context, query string, args ...interface{}) ([]models.Glasses, error) {
	var al []models.Glasses

	rows, err := r.pgpool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var a models.Glasses
		if err := rows.Scan(
			&a.GlassesID, &a.Color, &a.Brand,
			&a.LeftPrescription.Sph, &a.LeftPrescription.Cyl,
			&a.LeftPrescription.Axis, &a.LeftPrescription.Add,
			&a.RightPrescription.Sph, &a.RightPrescription.Cyl,
			&a.RightPrescription.Axis, &a.RightPrescription.Add,
			&a.Reference, &a.Type, &a.IsInStock,
			&a.Feature, &a.UpdatedAt, &a.CreatedAt,
		); err != nil {
			slog.Error(" scanning glasses", "err", err)
			return nil, errors.New("internal server error")
		}
		al = append(al, a)
	}

	if err := rows.Err(); err != nil {
		slog.Error(" fetching glasses", "err", err)
		return nil, errors.New("internal server error")
	}

	slog.Info("Glasses fetched", "glasses", al)
	return al, nil
}

func (r *GlassesRepository) fetchGlassesDetails(ctx context.Context, query string, args ...interface{}) ([]models.Glasses, error) {
	var al []models.Glasses

	rows, err := r.pgpool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var a models.Glasses
		if err := rows.Scan(
			&a.UserName, &a.UserEmail,
			&a.LeftPrescription.Sph, &a.LeftPrescription.Cyl,
			&a.LeftPrescription.Axis, &a.LeftPrescription.Add,
			&a.RightPrescription.Sph, &a.RightPrescription.Cyl,
			&a.RightPrescription.Axis, &a.RightPrescription.Add,
			&a.Reference, &a.IsInStock, &a.UpdatedAt, &a.CreatedAt,
		); err != nil {
			slog.Error(" scanning glasses", "err", err)
			return nil, errors.New("internal server error")
		}
		al = append(al, a)
	}

	if err := rows.Err(); err != nil {
		slog.Error(" fetching glasses", "err", err)
		return nil, errors.New("internal server error")
	}

	slog.Info("Glasses fetched", "glasses", al)
	return al, nil
}

func (r *GlassesRepository) GetGlasses(ctx context.Context, page, pageSize int,
	orderBy, sortBy, reference string, leftEye, rightEye *float64) ([]models.Glasses, error) {
	dbCtx, dbSpan := tracer.Start(ctx, "db.GetGlasses")
	defer dbSpan.End()

	dbSpan.SetAttributes(
		semconv.DBSystemKey.String("postgresql"),
		semconv.DBNameKey.String("glasses"),
		semconv.DBStatementKey.String("SELECT glasses_id, color, brand, right_sph, left_sph, type, reference, is_in_stock, features, COALESCE(updated_at, '1970-01-01 00:00:00') AS updated_at, created_at FROM glasses g WHERE Trim(Upper(g.reference)) ILIKE trim(upper('%' || $4 || '%')) AND ($5::float8 IS NULL OR g.left_sph = $5) AND ($6::float8 IS NULL OR g.right_sph = $6) ORDER BY CASE WHEN $1 = 'Brand' THEN g.brand WHEN $1 = 'Color' THEN g.color WHEN $1 = 'Reference' THEN g.reference WHEN $1 = 'Type' THEN g.type WHEN $1 = 'Features' THEN g.features ELSE g.brand END OFFSET $2 LIMIT $3"),
		semconv.DBOperationKey.String("query"),
	)
	query := `SELECT glasses_id, color, brand, 
					 left_sph, left_cyl, left_axis, left_add,
					 right_sph, right_cyl, right_axis, right_add,
       				 reference, type, is_in_stock, features, 
					 COALESCE(updated_at, '1970-01-01 00:00:00') AS updated_at, created_at
			 	FROM glasses g
			 	WHERE Trim(Upper(g.reference)) ILIKE trim(upper('%' || $4 || '%'))
			 	AND ($5::float8 IS NULL OR g.left_sph = $5)
			 	AND ($6::float8 IS NULL OR g.right_sph = $6)
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
	return r.fetchGlasses(dbCtx, query, orderBy, offset, pageSize, reference, leftEye, rightEye)
}

func (r *GlassesRepository) GetGlassesByID(ctx context.Context, glassesID uuid.UUID) (*models.Glasses, error) {
	query := `SELECT glasses_id, color, brand, 
	 				right_sph, left_sph, type,
       				reference, is_in_stock, features, updated_at, created_at
				FROM glasses
				WHERE glasses_id = $1`
	var a models.Glasses

	if err := r.pgpool.QueryRow(ctx, query, glassesID).Scan(
		&a.GlassesID, &a.Color, &a.Brand,
		&a.LeftPrescription.Sph, &a.LeftPrescription.Cyl,
		&a.LeftPrescription.Axis, &a.LeftPrescription.Add,
		&a.RightPrescription.Sph, &a.RightPrescription.Cyl,
		&a.RightPrescription.Axis, &a.RightPrescription.Add,
		&a.Reference, &a.Type, &a.IsInStock,
		&a.Feature, &a.UpdatedAt, &a.CreatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Error("No rows", "err", err)
			return nil, errors.New("internal server error")
		}
		slog.Error(" scanning glasses", "err", err)
		return nil, errors.New("internal server error")
	}
	slog.Info("Glasses fetched", "glasses", a)
	return &a, nil
}

func (r *GlassesRepository) DeleteGlasses(ctx context.Context, glassesID uuid.UUID) error {
	query := `DELETE FROM glasses WHERE glasses_id = $1`
	_, err := r.pgpool.Exec(ctx, query, glassesID)
	if err != nil {
		slog.Error(" deleting glasses", "err", err)
		return errors.New("internal server error")
	}
	slog.Info("Deleted glasses", "glasses_id", glassesID)
	return err
}

func (r *GlassesRepository) UpdateGlasses(ctx context.Context, g models.GlassesForm) error {
	query := `
		UPDATE glasses
		SET color = $1, brand = $2, right_sph = $3, 
			left_sph = $4, reference = $5, type =$6,
		    features = $7, updated_at = NOW()
		WHERE glasses_id = $8
	`
	_, err := r.pgpool.Exec(ctx, query, g.Color, g.Brand,
		g.RightSph, g.LeftSph, g.Reference, g.Type, g.Feature, g.GlassesID)
	if err != nil {
		slog.Error(" updating glasses", "err", err)
		return errors.New("internal server error")
	}
	slog.Info("Updated glasses", "glasses_id", g.GlassesID)
	return err
}

func (r *GlassesRepository) InsertGlasses(ctx context.Context, g models.GlassesForm) error {
	query := `
		INSERT INTO glasses (reference, brand, right_sph, left_sph, 
							color, type, features,
		                     is_in_stock, created_at, updated_at, user_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, true, NOW(), NOW(), $8)
		RETURNING glasses_id
	`
	err := r.pgpool.QueryRow(ctx, query, g.Reference, g.Brand, g.RightSph, g.LeftSph,
		g.Color, g.Type, g.Feature, g.UserID).Scan(&g.GlassesID)

	if err != nil {
		slog.Error(" inserting glasses", "err", err)
		return errors.New("internal server error")
	}
	slog.Info("Inserted glasses", "glasses_id", g.GlassesID)
	return err
}

func (r *GlassesRepository) GetSum(ctx context.Context) (int, error) {
	var count int
	row := r.pgpool.QueryRow(ctx, `SELECT Count(DISTINCT g.glasses_id) FROM glasses g`)
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	slog.Info("Glasses count", "count", count)
	return count, nil
}

func (r *GlassesRepository) GetGlassesByType(ctx context.Context,
	page, pageSize int, orderBy, sortBy, glassesType string) ([]models.Glasses, error) {
	query := `SELECT glasses_id, color, brand, 
					 left_sph, left_cyl, left_axis, left_add,
					 right_sph, right_cyl, right_axis, right_add,
       				 reference, type, is_in_stock, features, 
					 COALESCE(updated_at, '1970-01-01 00:00:00') AS updated_at, created_at
			 	FROM glasses g
			 	where type = $5
			 	ORDER BY
			    CASE
			        WHEN $1 = 'Brand' AND $2 = 'ASC' THEN g.brand
			        WHEN $1 = 'Brand' AND $2 = 'DESC' THEN g.brand
			    END,
			    g.created_at
			    OFFSET $3 LIMIT $4`
	offset := (page - 1) * pageSize

	slog.Info("Fetching glasses", "page", page, "pageSize", pageSize, "offset", offset)
	return r.fetchGlasses(ctx, query, orderBy, sortBy, offset, pageSize, glassesType)
}

func (r *GlassesRepository) GetGlassesByStock(ctx context.Context,
	page, pageSize int, orderBy, sortBy string, isInStock bool) ([]models.Glasses, error) {
	query := `SELECT glasses_id, color, brand, 
					 right_sph, left_sph, type,
                     reference, is_in_stock, features,  
					 COALESCE(updated_at, '1970-01-01 00:00:00') AS updated_at, created_at
                 FROM glasses g
                 WHERE is_in_stock = $5
                 ORDER BY
                 CASE
                     WHEN $1 = 'Brand' AND $2 = 'ASC' THEN g.brand
                     WHEN $1 = 'Brand' AND $2 = 'DESC' THEN g.brand
                 END,
                 g.created_at
                 OFFSET $3 LIMIT $4`
	offset := (page - 1) * pageSize

	slog.Info("Fetching glasses", "page", page, "pageSize", pageSize, "offset", offset)
	return r.fetchGlasses(ctx, query, orderBy, sortBy, offset, pageSize, isInStock)
}

func (r *GlassesRepository) GetSumByType(ctx context.Context, glassesType string) (int, error) {
	var count int
	row := r.pgpool.QueryRow(ctx, `
			SELECT Count(DISTINCT g.glasses_id)
			FROM glasses g
			WHERE g.type = $1`, glassesType)
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	slog.Info("Glasses count", "count", count)
	return count, nil
}

func (r *GlassesRepository) GetSumByStock(ctx context.Context, isInStock bool) (int, error) {
	var count int
	row := r.pgpool.QueryRow(ctx, `
			SELECT Count(DISTINCT g.glasses_id)
			FROM glasses g
			WHERE g.is_in_stock = $1`, isInStock)
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	slog.Info("Glasses count", "count", count)
	return count, nil
}

func (r *GlassesRepository) GetGlassesReference(ctx context.Context, id uuid.UUID) (string, error) {
	query := `SELECT reference FROM glasses WHERE glasses_id = $1`
	var reference string
	if err := r.pgpool.QueryRow(ctx, query, id).Scan(&reference); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Error("No rows", "err", err)
			return "", errors.New("internal server error")
		}
		slog.Error("scanning glasses", "err", err)
		return "", errors.New("internal server error")
	}

	slog.Info("Glasses fetched", "reference", reference)
	return reference, nil
}
