package store

import (
	"context"

	"github.com/zagvozdeen/ola/internal/store/models"
)

func (s *Store) CreateOrder(ctx context.Context, order *models.Order) error {
	_, err := s.pool.Exec(ctx, "INSERT INTO orders (id, uuid, name, phone, content, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", order.ID, order.UUID, order.Name, order.Phone, order.Content, order.UserID, order.CreatedAt, order.UpdatedAt)
	return err
}
