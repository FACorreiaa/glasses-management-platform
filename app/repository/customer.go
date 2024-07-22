package repository

import (
	"context"
	"errors"
	"log/slog"

	"github.com/FACorreiaa/glasses-management-platform/app/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CustomerRepository struct {
	pgpool *pgxpool.Pool
}

func NewCustomerRepository(db *pgxpool.Pool) *CustomerRepository {
	return &CustomerRepository{pgpool: db}
}

func (r *CustomerRepository) InsertShippingDetails(ctx context.Context, glassesID, userID uuid.UUID,
	c models.CustomerShippingForm, s models.Shipping) error {
	tx, err := r.pgpool.Begin(ctx)
	if err != nil {
		slog.Error(" starting transaction", "err", err)
		return errors.New("internal server error")
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		if err := tx.Rollback(ctx); err != nil {
			slog.Error(" rolling back transaction", "err", err)
		}
	}(tx, ctx)

	customerQuery := `INSERT INTO customer (glasses_id, user_id, name, card_id_number, address, address_details, city,
                      country, continent, postal_code, phone_number, email, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW(), NOW()) RETURNING customer_id`

	var customerID uuid.UUID
	if err := tx.QueryRow(ctx, customerQuery, glassesID, userID, c.Name, c.CardID, c.Address, c.AddressDetails, c.City,
		c.Country, c.Continent, c.PostalCode, c.PhoneNumber, c.Email).Scan(&customerID); err != nil {
		slog.Error(" inserting customer shipping details", "err", err)
		return errors.New("internal server error")
	}

	updateStockQuery := `UPDATE glasses SET is_in_stock = false, updated_at = Now() WHERE glasses_id = $1 RETURNING glasses_id`
	var updatedGlassesID uuid.UUID
	if err := tx.QueryRow(ctx, updateStockQuery, glassesID).Scan(&updatedGlassesID); err != nil {
		slog.Error("updating glasses stock", "err", err)
		return errors.New("internal server error")
	}

	shippingQuery := `
		INSERT INTO shipping_details (glasses_id, customer_id, shipped_by, shipping_date, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW(), NOW())
		RETURNING shipping_id
	`

	var shippingID uuid.UUID
	err = tx.QueryRow(ctx, shippingQuery, glassesID, customerID, userID).Scan(&shippingID)
	if err != nil {
		slog.Error("inserting shipping details", "err", err)
		return errors.New("internal server error")
	}

	if err := tx.Commit(ctx); err != nil {
		slog.Error("committing transaction", "err", err)
		return errors.New("internal server error")
	}

	slog.Info("Shipping details inserted and glasses stock status updated", "shipping_id", shippingID, "glasses_id", updatedGlassesID)

	return nil
}

func (r *CustomerRepository) GetCardIDNumber(ctx context.Context, userID uuid.UUID) (string, error) {
	query := `SELECT card_id_number FROM customer WHERE user_id = $1`
	var cardID string
	_ = r.pgpool.QueryRow(ctx, query, userID).Scan(&cardID)
	slog.Info("Card id number fetched", "card_id", cardID)
	return cardID, nil
}

func (r *CustomerRepository) GetShippingDetails(ctx context.Context, page, pageSize int,
	orderBy, sortBy, reference string, leftEye, rightEye *float64) ([]models.ShippingDetails, error) {
	var sd []models.ShippingDetails
	query := `select c.name, card_id_number, email, g.reference,
       			g.left_eye_strength, g.right_eye_strength,
       			c.created_at, c.updated_at
				from customer c
				join glasses g on g.glasses_id = c.glasses_id
				WHERE Trim(Upper(g.reference)) ILIKE trim(upper('%' || $1 || '%'))
			 	AND ($2::float8 IS NULL OR g.left_eye_strength = $2)
			 	AND ($3::float8 IS NULL OR g.right_eye_strength = $3)
				ORDER BY
				CASE
					WHEN $4 = 'Brand' THEN g.brand
					WHEN $4 = 'Color' THEN g.color
					WHEN $4 = 'Reference' THEN g.reference
					WHEN $4 = 'Type' THEN g.type
					WHEN $4 = 'Features' THEN g.features

					ELSE g.brand
				END ` + sortBy + `
			    OFFSET $5 LIMIT $6`
	offset := (page - 1) * pageSize
	rows, err := r.pgpool.Query(ctx, query, reference, leftEye, rightEye, orderBy, offset, pageSize)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s models.ShippingDetails
		if err := rows.Scan(&s.Name, &s.CardID, &s.Email, &s.Reference,
			&s.LeftEyeStrength, &s.RightEyeStrength, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		sd = append(sd, s)
	}

	slog.Info("Fetching shippingDetails", "page", page, "pageSize", pageSize, "offset", offset)

	return sd, nil
}

func (r *CustomerRepository) GetShippingExpandedDetails(ctx context.Context, page, pageSize int,
	orderBy, sortBy, reference string, leftEye, rightEye *float64) ([]models.SettingsShippingDetails, error) {
	var sd []models.SettingsShippingDetails
	query := `select u.username, u.email as "collaborator_email", c.name, c.card_id_number, c.email, g.reference,
       			g.left_eye_strength, g.right_eye_strength,
       			c.created_at, c.updated_at
				from customer c
				join glasses g on g.glasses_id = c.glasses_id
				join "user" u on u.user_id = c.user_id
				WHERE Trim(Upper(g.reference)) ILIKE trim(upper('%' || $1 || '%'))
			 	AND ($2::float8 IS NULL OR g.left_eye_strength = $2)
			 	AND ($3::float8 IS NULL OR g.right_eye_strength = $3)
				ORDER BY
				CASE
					WHEN $4 = 'Brand' THEN g.brand
					WHEN $4 = 'Color' THEN g.color
					WHEN $4 = 'Reference' THEN g.reference
					WHEN $4 = 'Type' THEN g.type
					WHEN $4 = 'Features' THEN g.features

					ELSE g.brand
				END ` + sortBy + `
			    OFFSET $5 LIMIT $6`
	offset := (page - 1) * pageSize
	rows, err := r.pgpool.Query(ctx, query, reference, leftEye, rightEye, orderBy, offset, pageSize)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s models.SettingsShippingDetails
		if err := rows.Scan(&s.CollaboratorName, &s.CollaboratorEmail, &s.Name, &s.CardID, &s.Email, &s.Reference,
			&s.LeftEyeStrength, &s.RightEyeStrength, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		sd = append(sd, s)
	}

	slog.Info("Fetching shippingDetails", "page", page, "pageSize", pageSize, "offset", offset)

	return sd, nil
}

func (r *GlassesRepository) DeleteCustomer(ctx context.Context, customerID string) error {
	query := `DELETE FROM customer WHERE card_id_number = $1`
	_, err := r.pgpool.Exec(ctx, query, customerID)
	if err != nil {
		slog.Error(" deleting customer", "err", err)
		return errors.New("internal server error")
	}
	slog.Info("Deleted customer", "card_id", customerID)
	return err
}
