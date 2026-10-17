package repository

import (
	"context"
	"errors"
	"fmt"
	"uuid"

	"github.com/ItakawaM/grpc-product-service/internal/productservice"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// postgres implements productservice.Repository with PostgreSQL.
type postgres struct {
	db *pgxpool.Pool
}

// NewPostgres initializes a Product Repository with PostgreSQL.
func NewPostgres(db *pgxpool.Pool) productservice.Repository {
	return &postgres{
		db: db,
	}
}

// getProductQuery retrieves a productservice.Product entity by its UUID.
const getProductQuery = `
	SELECT id, name, description, price_cents, created_at, updated_at, deleted_at, version
	FROM products 
	WHERE id = $1 AND deleted_at IS NULL;
`

// Get implements productservice.Repository.
func (p *postgres) Get(ctx context.Context, id uuid.UUID) (*productservice.Product, error) {
	var product productservice.Product
	if err := p.db.QueryRow(ctx, getProductQuery, id).
		Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.PriceCents,
			&product.CreatedAt,
			&product.UpdatedAt,
			&product.DeletedAt,
			&product.Version,
		); err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, productservice.ErrNotFound

		default:
			return nil, fmt.Errorf("get product %s: %w", id, err)
		}
	}

	return &product, nil
}

// createProductQuery adds a new productservice.Product entity.
const createProductQuery = `
	INSERT INTO products (name, description, price_cents) 
	VALUES ($1, $2, $3)
	RETURNING id, created_at, updated_at, deleted_at, version;
`

// Create implements productservice.Repository.
func (p *postgres) Create(ctx context.Context, product *productservice.Product) error {
	if err := p.db.QueryRow(ctx, createProductQuery, product.Name, product.Description, product.PriceCents).
		Scan(
			&product.ID,
			&product.CreatedAt,
			&product.UpdatedAt,
			&product.DeletedAt,
			&product.Version,
		); err != nil {
		return fmt.Errorf("create product: %w", err)
	}

	return nil
}
