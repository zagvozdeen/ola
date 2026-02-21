package store

import (
	"context"

	"github.com/google/uuid"
	"github.com/zagvozdeen/ola/internal/store/models"
)

func (s *Store) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	panic("implement")
}

func (s *Store) GetProductByUUID(ctx context.Context, uuid uuid.UUID) (*models.Product, error) {
	panic("implement")
}

func (s *Store) CreateProduct(ctx context.Context, service *models.Product) error {
	panic("implement")
}

func (s *Store) UpdateProduct(ctx context.Context, service *models.Product) error {
	panic("implement")
}
