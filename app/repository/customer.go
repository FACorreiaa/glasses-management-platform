package repository

import (
	"context"
	"errors"
	"fmt"
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
	query := `select c.customer_id, c.name, card_id_number, email, g.reference,
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
		if err := rows.Scan(&s.CustomerID, &s.Name, &s.CardID, &s.Email, &s.Reference,
			&s.LeftEye, &s.RightEye, &s.CreatedAt, &s.UpdatedAt); err != nil {
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
       			g.left_eye_strength, g.right_eye_strength, c.customer_id,
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
			&s.LeftEye, &s.RightEye, &s.CustomerID, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		sd = append(sd, s)
	}

	slog.Info("Fetching shippingDetails", "page", page, "pageSize", pageSize, "offset", offset)

	return sd, nil
}

func (r *GlassesRepository) DeleteCustomer(ctx context.Context, customerID uuid.UUID) error {
	query := `DELETE FROM customer WHERE customer_id = $1`
	_, err := r.pgpool.Exec(ctx, query, customerID)
	if err != nil {
		slog.Error(" deleting customer", "err", err)
		return errors.New("internal server error")
	}
	slog.Info("Deleted customer", "card_id", customerID)
	return err
}

func (r *CustomerRepository) UpdateShippingDetails(ctx context.Context, form models.ShippingDetailsForm, id uuid.UUID) error {
	tx, err := r.pgpool.Begin(ctx)
	if err != nil {
		slog.Error("starting transaction", "err", err)
		return errors.New("internal server error")
	}
	fmt.Println("ID ON QUERY", id.String())

	defer func(tx pgx.Tx, ctx context.Context) {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			slog.Error("rolling back transaction", "err", err)
		}
	}(tx, ctx)

	customerUpdateQuery := `
        UPDATE customer
        SET name = $1, card_id_number = $2, email = $3, updated_at = NOW()
        WHERE customer_id = $4
        RETURNING customer_id
    `
	// var customerID uuid.UUID
	if err := tx.QueryRow(ctx, customerUpdateQuery, form.Name, form.CardID, form.Email, id).Scan(&id); err != nil {
		slog.Error("updating customer details", "err", err)
		return errors.New("internal server error")
	}

	glassesUpdateQuery := `
        UPDATE glasses
        SET reference = $1, left_eye_strength = $2, right_eye_strength = $3, updated_at = NOW()
        WHERE glasses_id = (
            SELECT glasses_id
            FROM customer
            WHERE customer_id = $4
        )
        RETURNING glasses_id
    `
	var glassesID uuid.UUID
	if err := tx.QueryRow(ctx, glassesUpdateQuery, form.Reference, form.LeftEye, form.RightEye, id).Scan(&glassesID); err != nil {
		slog.Error("updating glasses details", "err", err)
		return errors.New("internal server error")
	}

	if err := tx.Commit(ctx); err != nil {
		slog.Error("committing transaction", "err", err)
		return errors.New("internal server error")
	}

	slog.Info("Shipping details updated", "customer_id", id, "glasses_id", glassesID)
	return nil
}

func (r *CustomerRepository) GetCustomerGlassesID(ctx context.Context, customerID uuid.UUID) (*models.ShippingDetails, error) {
	query := `SELECT c.customer_id, g.right_eye_strength, g.left_eye_strength
				FROM glasses g
				JOIN customer c ON c.glasses_id = g.glasses_id
				WHERE c.customer_id = $1`
	var a models.ShippingDetails

	err := r.pgpool.QueryRow(ctx, query, customerID).Scan(
		&a.CustomerID, &a.RightEye, &a.LeftEye,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Error("No rows", "err", err)
			return nil, errors.New("internal server error")
		}
		slog.Error(" scanning glasses", "err", err)
		return nil, errors.New("internal server error")
	}

	slog.Info("Customer fetched", "customer", a)
	return &a, nil
}

func (r *CustomerRepository) GetCardIDFromShipping(ctx context.Context, customerID uuid.UUID) (string, error) {
	query := `SELECT card_id_number FROM customer WHERE customer_id = $1`
	var cardID string
	if err := r.pgpool.QueryRow(ctx, query, customerID).Scan(&cardID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Error("No card_id_number found for customer", "customer_id", customerID)
			return "", nil
		}
		slog.Error("Error fetching card_id_number", "err", err)
		return "", err
	}
	slog.Info("Card id number fetched", "card_id", cardID)
	return cardID, nil
}

func (r *CustomerRepository) GetReferenceNumberFromShipping(ctx context.Context, customerID uuid.UUID) (string, error) {
	query := `SELECT reference FROM glasses
              join customer on customer.glasses_id = glasses.glasses_id
        	  WHERE customer_id = $1`
	var reference string
	if err := r.pgpool.QueryRow(ctx, query, customerID).Scan(&reference); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Error("No reference found for customer", "customer_id", customerID)
			return "", nil
		}
		slog.Error("Error fetching reference", "err", err)
		return "", err
	}
	slog.Info("Reference number fetched", "reference", reference)
	return reference, nil
}
