package store

import (
	"context"

	"github.com/zagvozdeen/ola/internal/store/models"
)

func (s *Store) GetAllOrders(ctx context.Context) ([]models.Order, error) {
	rows, err := s.pool.Query(ctx, "SELECT id, uuid::text, name, phone, content, user_id, created_at, updated_at FROM orders ORDER BY created_at DESC")
	if err != nil {
		return nil, wrapDBError(err)
	}
	defer rows.Close()

	orders := make([]models.Order, 0)
	for rows.Next() {
		order := models.Order{}
		err = rows.Scan(&order.ID, &order.UUID, &order.Name, &order.Phone, &order.Content, &order.UserID, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, wrapDBError(err)
		}
		orders = append(orders, order)
	}
	if err = rows.Err(); err != nil {
		return nil, wrapDBError(err)
	}

	return orders, nil
}

func (s *Store) CreateOrder(ctx context.Context, order *models.Order) error {
	err := s.pool.QueryRow(
		ctx,
		"INSERT INTO orders (uuid, name, phone, content, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		order.UUID, order.Name, order.Phone, order.Content, order.UserID, order.CreatedAt, order.UpdatedAt,
	).Scan(&order.ID)
	return wrapDBError(err)
}
