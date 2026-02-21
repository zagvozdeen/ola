package store

import (
	"context"

	"github.com/google/uuid"
	"github.com/zagvozdeen/ola/internal/store/models"
)

func (s *Store) GetAllServices(ctx context.Context) ([]models.Service, error) {
	rows, err := s.pool.Query(ctx, "SELECT s.id, s.uuid::text, s.name, s.description, s.price_from, s.price_to, s.file_id, f.content, s.user_id, s.created_at, s.updated_at FROM services s JOIN files f ON f.id = s.file_id ORDER BY s.created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	services := make([]models.Service, 0)
	for rows.Next() {
		service := models.Service{}
		fileContent := ""
		err = rows.Scan(&service.ID, &service.UUID, &service.Name, &service.Description, &service.PriceFrom, &service.PriceTo, &service.FileID, &fileContent, &service.UserID, &service.CreatedAt, &service.UpdatedAt)
		if err != nil {
			return nil, err
		}
		service.FileContent = &fileContent
		services = append(services, service)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return services, nil
}

func (s *Store) GetServiceByUUID(ctx context.Context, uuid uuid.UUID) (*models.Service, error) {
	service := &models.Service{}
	fileContent := ""
	err := s.pool.QueryRow(ctx, "SELECT s.id, s.uuid::text, s.name, s.description, s.price_from, s.price_to, s.file_id, f.content, s.user_id, s.created_at, s.updated_at FROM services s JOIN files f ON f.id = s.file_id WHERE s.uuid = $1", uuid.String()).Scan(&service.ID, &service.UUID, &service.Name, &service.Description, &service.PriceFrom, &service.PriceTo, &service.FileID, &fileContent, &service.UserID, &service.CreatedAt, &service.UpdatedAt)
	if err != nil {
		return nil, err
	}
	service.FileContent = &fileContent
	return service, nil
}

func (s *Store) CreateService(ctx context.Context, service *models.Service) error {
	_, err := s.pool.Exec(ctx, "INSERT INTO services (id, uuid, name, description, price_from, price_to, file_id, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", service.ID, service.UUID, service.Name, service.Description, service.PriceFrom, service.PriceTo, service.FileID, service.UserID, service.CreatedAt, service.UpdatedAt)
	return err
}

func (s *Store) UpdateService(ctx context.Context, service *models.Service) error {
	_, err := s.pool.Exec(ctx, "UPDATE services SET name = $1, description = $2, price_from = $3, price_to = $4, file_id = $5, user_id = $6, updated_at = $7 WHERE id = $8", service.Name, service.Description, service.PriceFrom, service.PriceTo, service.FileID, service.UserID, service.UpdatedAt, service.ID)
	return err
}
