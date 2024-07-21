package repository

import "github.com/jackc/pgx/v5/pgxpool"

type CustomerRepository struct {
	pgpool *pgxpool.Pool
}

func NewCustomerRepository(db *pgxpool.Pool) *CustomerRepository {
	return &CustomerRepository{pgpool: db}
}
