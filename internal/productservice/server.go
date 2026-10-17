package productservice

import (
	"context"
	"errors"
	"time"
	"uuid"

	pb "github.com/ItakawaM/grpc-product-service/gen/go/productservice/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// server implements ProductServiceServer.
type server struct {
	products Repository
	pb.UnimplementedProductServiceServer
}

// NewServer initializes new ProductServiceServer with the provided repository.
func NewServer(products Repository) pb.ProductServiceServer {
	return &server{
		products: products,
	}
}

// CreateProduct creates a new Product entity in the database with the provided
// name, description and price in cents.
//
// It errors on database failure or context cancellation.
func (s *server) CreateProduct(ctx context.Context, in *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	product := &Product{
		Name:        in.GetName(),
		Description: in.GetDescription(),
		PriceCents:  in.GetPriceCents(),
	}

	if err := s.products.Create(ctx, product); err != nil {
		return nil, s.toStatusError(ctx, err)
	}

	return &pb.CreateProductResponse{
		Product: product.toProto(),
	}, nil
}

// GetProduct retrieves a Product by its UUID.
//
// It returns codes.InvalidArgument when the provided UUID is not valid
// and codes.NotFound when no Product exists.
func (s *server) GetProduct(ctx context.Context, in *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	id, err := uuid.Parse(in.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product id")
	}

	product, err := s.products.Get(ctx, id)
	if err != nil {
		return nil, s.toStatusError(ctx, err)
	}

	return &pb.GetProductResponse{
		Product: product.toProto(),
	}, nil
}

// TODO: REPLACE WITH AN INTERCEPTOR LATER
func (s *server) toStatusError(ctx context.Context, err error) error {
	switch {
	case errors.Is(err, ErrNotFound):
		return status.Error(codes.NotFound, "product not found")

	case errors.Is(err, context.DeadlineExceeded):
		return status.Error(codes.DeadlineExceeded, "request timed out")

	case errors.Is(err, context.Canceled):
		return status.Error(codes.Canceled, "request canceled")

	default:
		// TODO: ADD LOGGING LATER
		return status.Error(codes.Internal, "internal error")
	}
}
