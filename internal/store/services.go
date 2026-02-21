package store

import (
	"context"

	"github.com/google/uuid"
	"github.com/zagvozdeen/ola/internal/store/models"
)

func (s *Store) GetAllServices(ctx context.Context) ([]models.Service, error) {
	panic("implement")
}

func (s *Store) GetServiceByUUID(ctx context.Context, uuid uuid.UUID) (*models.Service, error) {
	panic("implement")
}

func (s *Store) CreateService(ctx context.Context, service *models.Service) error {
	panic("implement")
}

func (s *Store) UpdateService(ctx context.Context, service *models.Service) error {
	panic("implement")
}
