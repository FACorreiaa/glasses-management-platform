package repository

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

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
		// This Scan expects exactly 8 columns from the query:
		if err := rows.Scan(
			&a.UserName,              // 1
			&a.UserEmail,             // 2
			&a.LeftPrescription.Sph,  // 3
			&a.RightPrescription.Sph, // 4
			&a.Reference,             // 5
			&a.IsInStock,             // 6
			&a.UpdatedAt,             // 7
			&a.CreatedAt,             // 8
		); err != nil {
			slog.Error(" scanning glasses details", "err", err) // Adjusted log message
			return nil, errors.New("internal server error")
		}
		al = append(al, a)
	}

	if err := rows.Err(); err != nil {
		slog.Error(" fetching glasses details", "err", err) // Adjusted log message
		return nil, errors.New("internal server error")
	}

	slog.Info("Glasses details fetched", "glasses", al) // Optional: Log fetched data
	return al, nil
}

func (r *GlassesRepository) GetGlasses(ctx context.Context, page, pageSize int,
	orderBy, sortBy, reference string, leftEye, rightEye *float64) ([]models.Glasses, error) {
	dbCtx, dbSpan := tracer.Start(ctx, "db.GetGlasses")
	defer dbSpan.End()

	// Validate and default sortBy
	if sortBy != "ASC" && sortBy != "DESC" {
		sortBy = "ASC" // Default to ASC if sortBy is invalid or empty
	}

	dbSpan.SetAttributes(
		semconv.DBSystemKey.String("postgresql"),
		semconv.DBNameKey.String("glasses"),
		semconv.DBStatementKey.String("SELECT glasses_id, color, brand, right_sph, left_sph, type, reference, is_in_stock, features, COALESCE(updated_at, '1970-01-01 00:00:00') AS updated_at, created_at FROM glasses g WHERE Trim(Upper(g.reference)) ILIKE trim(upper('%' || $4 || '%')) AND ($5::float8 IS NULL OR g.left_sph = $5) AND ($6::float8 IS NULL OR g.right_sph = $6) ORDER BY CASE WHEN $1 = 'Brand' THEN g.brand WHEN $1 = 'Color' THEN g.color WHEN $1 = 'Reference' THEN g.reference WHEN $1 = 'Type' THEN g.type WHEN $1 = 'Features' THEN g.features ELSE g.brand END "+sortBy),
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
					WHEN UPPER($1) = 'BRAND' THEN g.brand
					WHEN UPPER($1) = 'COLOR' THEN g.color
					WHEN UPPER($1) = 'REFERENCE' THEN g.reference
					WHEN UPPER($1) = 'TYPE' THEN g.type
					WHEN UPPER($1) = 'FEATURES' THEN g.features
					-- WHEN UPPER($1) = 'LEFT_SPH' THEN g.left_sph
					-- WHEN UPPER($1) = 'RIGHT_SPH' THEN g.right_sph
					ELSE g.brand
				END ` + sortBy + `
              OFFSET $2 LIMIT $3`
	offset := (page - 1) * pageSize
	slog.Info("Glasses fetched", "offset", offset, "sortBy", sortBy)
	return r.fetchGlasses(dbCtx, query, orderBy, offset, pageSize, reference, leftEye, rightEye)
}

func (r *GlassesRepository) GetGlassesByID(ctx context.Context, glassesID uuid.UUID) (*models.Glasses, error) {
	query := `SELECT glasses_id, color, brand,
					 left_sph, left_cyl, left_axis, left_add,
					 right_sph, right_cyl, right_axis, right_add,
       				 reference, type, is_in_stock, features,
					 COALESCE(updated_at, '1970-01-01 00:00:00') AS updated_at, created_at
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

func (r *GlassesRepository) UpdateGlasses(ctx context.Context, g models.GlassesForm, restock bool, currentStockStatus bool) error {
	// Start a transaction for atomicity
	tx, err := r.pgpool.Begin(ctx)
	if err != nil {
		slog.Error("starting transaction for glasses update", "err", err)
		return errors.New("internal server error: could not start transaction")
	}
	// Ensure rollback happens if commit fails or panics occur
	defer func() {
		if err != nil { // Rollback if any error occurred during the function
			slog.Warn("rolling back transaction due to error", "err", err)
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				slog.Error("error rolling back transaction", "rollbackErr", rollbackErr)
			}
		}
	}()

	// Scenario 1: Restocking an item that was previously out of stock
	if restock && !currentStockStatus {
		slog.Info("Restocking glasses, deleting associated customer and shipping details", "glasses_id", g.GlassesID)

		shippingDeleteQuery := `DELETE FROM shipping_details WHERE glasses_id = $1`
		cmdTagShipping, errShip := tx.Exec(ctx, shippingDeleteQuery, g.GlassesID)
		if errShip != nil {
			err = fmt.Errorf("deleting shipping details: %w", errShip) // Capture error for defer rollback
			slog.Error("deleting shipping details during restock", "err", err, "glasses_id", g.GlassesID)
			return errors.New("internal server error: could not delete shipping details")
		}
		slog.Debug("Shipping details deletion result", "rows_affected", cmdTagShipping.RowsAffected(), "glasses_id", g.GlassesID)

		customerDeleteQuery := `DELETE FROM customer WHERE glasses_id = $1`
		cmdTagCustomer, errCust := tx.Exec(ctx, customerDeleteQuery, g.GlassesID)
		if errCust != nil {
			err = fmt.Errorf("deleting customer: %w", errCust) // Capture error for defer rollback
			slog.Error("deleting customer during restock", "err", err, "glasses_id", g.GlassesID)
			return errors.New("internal server error: could not delete customer")
		}
		if cmdTagCustomer.RowsAffected() == 0 {
			// This might be okay if glasses were out of stock but somehow had no customer link
			slog.Warn("No customer record found to delete during restock", "glasses_id", g.GlassesID)
		} else {
			slog.Info("Customer record deleted during restock", "glasses_id", g.GlassesID)
		}

		// 3. Update Glasses - Set IsInStock to TRUE and update other fields
		updateQuery := `
			UPDATE glasses
			SET color = $1, brand = $2,
			    reference = $3, type = $4, features = $5,
                left_sph = $6, left_cyl = $7, left_axis = $8, left_add = $9,
                right_sph = $10, right_cyl = $11, right_axis = $12, right_add = $13,
			    is_in_stock = TRUE, -- Set stock to true
                updated_at = NOW()
			WHERE glasses_id = $14` // Parameter count increased
		_, errExec := tx.Exec(ctx, updateQuery,
			g.Color, g.Brand, g.Reference, g.Type, g.Feature, // $1 - $5
			g.LeftSph, g.LeftCyl, g.LeftAxis, g.LeftAdd, // $6 - $9
			g.RightSph, g.RightCyl, g.RightAxis, g.RightAdd, // $10 - $13
			g.GlassesID) // $14
		if errExec != nil {
			err = fmt.Errorf("updating glasses during restock: %w", errExec) // Capture error
			slog.Error("updating glasses during restock", "err", err, "glasses_id", g.GlassesID)
			return errors.New("internal server error: could not update glasses")
		}
		slog.Info("Glasses updated and restocked", "glasses_id", g.GlassesID)

	} else {
		// Scenario 2: Standard update (or restock checkbox wasn't checked, or item was already in stock)
		slog.Info("Performing standard update for glasses", "glasses_id", g.GlassesID, "restock_flag", restock, "current_stock", currentStockStatus)
		// Original update query (does NOT touch is_in_stock)
		updateQuery := `
			UPDATE glasses
			SET color = $1, brand = $2,
			    reference = $3, type = $4, features = $5,
			    left_sph = $6, left_cyl = $7, left_axis = $8, left_add = $9,
                right_sph = $10, right_cyl = $11, right_axis = $12, right_add = $13,
			    updated_at = NOW()
			WHERE glasses_id = $14` // Parameter count increased
		_, errExec := tx.Exec(ctx, updateQuery,
			g.Color, g.Brand, g.Reference, g.Type, g.Feature, // $1 - $5
			g.LeftSph, g.LeftCyl, g.LeftAxis, g.LeftAdd, // $6 - $9
			g.RightSph, g.RightCyl, g.RightAxis, g.RightAdd, // $10 - $13
			g.GlassesID) // $14
		if errExec != nil {
			err = fmt.Errorf("updating glasses: %w", errExec) // Capture error
			slog.Error("updating glasses", "err", err, "glasses_id", g.GlassesID)
			return errors.New("internal server error")
		}
		slog.Info("Updated glasses (standard)", "glasses_id", g.GlassesID)
	}

	if err = tx.Commit(ctx); err != nil {
		slog.Error("committing transaction for glasses update", "err", err, "glasses_id", g.GlassesID)
		return errors.New("internal server error: could not commit changes")
	}

	return nil
}

func (r *GlassesRepository) InsertGlasses(ctx context.Context, g models.GlassesForm) error {
	query := `
		INSERT INTO glasses (reference, brand,
							right_sph, right_cyl, right_axis, right_add,
							left_sph, left_cyl, left_axis, left_add,
							color, type, features,
		                     is_in_stock, created_at, updated_at, user_id)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
		RETURNING glasses_id
	`
	err := r.pgpool.QueryRow(ctx, query,
		g.Reference, g.Brand,
		g.RightSph, g.RightCyl, g.RightAxis, g.RightAdd,
		g.LeftSph, g.LeftCyl, g.LeftAxis, g.LeftAdd,
		g.Color, g.Type, g.Feature,
		true,
		time.Now(),
		time.Now(),
		g.UserID,
	).Scan(&g.GlassesID)

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
					 left_sph, left_cyl, left_axis, left_add,
					 right_sph, right_cyl, right_axis, right_add,
       				 reference, type, is_in_stock, features,
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
