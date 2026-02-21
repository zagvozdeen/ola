package store

import (
	"context"

	"github.com/zagvozdeen/ola/internal/store/models"
)

func (s *Store) CreateFeedback(ctx context.Context, feedback *models.Feedback) error {
	_, err := s.pool.Exec(ctx, "INSERT INTO feedback (id, uuid, name, phone, content, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", feedback.ID, feedback.UUID, feedback.Name, feedback.Phone, feedback.Content, feedback.UserID, feedback.CreatedAt, feedback.UpdatedAt)
	return err
}
