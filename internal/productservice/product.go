package productservice

import (
	"context"
	"errors"
	"time"
	"uuid"

	pb "github.com/ItakawaM/grpc-product-service/gen/go/productservice/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	// ErrNotFound is returned when a product with the provided UUID does not exist.
	ErrNotFound = errors.New("product not found")
)

// Product is a single business entity stored in the database.
type Product struct {
	ID          uuid.UUID
	Name        string
	Description string
	PriceCents  uint32
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
	Version     uint32
}

// toProto maps Product to its appropriate proto counterpart.
func (p *Product) toProto() *pb.Product {
	out := &pb.Product{
		Id:          p.ID.String(),
		Name:        p.Name,
		Description: p.Description,
		PriceCents:  p.PriceCents,
		CreatedAt:   timestamppb.New(p.CreatedAt),
		UpdatedAt:   timestamppb.New(p.UpdatedAt),
		Version:     p.Version,
	}

	if p.DeletedAt != nil {
		out.DeletedAt = timestamppb.New(*p.DeletedAt)
	}

	return out
}

// Repository defines the storage operations available for Product.
type Repository interface {
	// Get retrieves a Product by its UUID.
	// Returns ErrNotFound if no product with that UUID exists.
	Get(ctx context.Context, id uuid.UUID) (*Product, error)

	// Create adds a new Product entity to the database; populates its
	// ID, CreatedAt, UpdatedAt, DeletedAt and Version fields on success.
	Create(ctx context.Context, product *Product) error
}
