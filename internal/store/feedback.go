package store

import (
	"context"

	"github.com/zagvozdeen/ola/internal/store/models"
)

func (s *Store) GetAllFeedback(ctx context.Context) ([]models.Feedback, error) {
	rows, err := s.pool.Query(ctx, "SELECT id, uuid, name, phone, content, user_id, created_at, updated_at FROM feedback ORDER BY created_at DESC")
	if err != nil {
		return nil, wrapDBError(err)
	}
	defer rows.Close()

	feedbacks := make([]models.Feedback, 0)
	for rows.Next() {
		var feedback models.Feedback
		err = rows.Scan(&feedback.ID, &feedback.UUID, &feedback.Name, &feedback.Phone, &feedback.Content, &feedback.UserID, &feedback.CreatedAt, &feedback.UpdatedAt)
		if err != nil {
			return nil, wrapDBError(err)
		}
		feedbacks = append(feedbacks, feedback)
	}
	if err = rows.Err(); err != nil {
		return nil, wrapDBError(err)
	}

	return feedbacks, nil
}

func (s *Store) CreateFeedback(ctx context.Context, feedback *models.Feedback) error {
	err := s.pool.QueryRow(
		ctx,
		"INSERT INTO feedback (uuid, name, phone, content, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		feedback.UUID, feedback.Name, feedback.Phone, feedback.Content, feedback.UserID, feedback.CreatedAt, feedback.UpdatedAt,
	).Scan(&feedback.ID)
	return wrapDBError(err)
}
