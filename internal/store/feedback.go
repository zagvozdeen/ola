package store

import (
	"context"

	"github.com/google/uuid"
	"github.com/zagvozdeen/ola/internal/store/models"
)

func (s *Store) GetAllFeedback(ctx context.Context) ([]models.Feedback, error) {
	rows, err := s.querier(ctx).Query(ctx, "SELECT id, uuid, status, source, type, name, phone, content, user_id, created_at, updated_at FROM feedback ORDER BY created_at DESC")
	if err != nil {
		return nil, wrapDBError(err)
	}
	defer rows.Close()

	feedbacks := make([]models.Feedback, 0)
	for rows.Next() {
		var feedback models.Feedback
		err = rows.Scan(&feedback.ID, &feedback.UUID, &feedback.Status, &feedback.Source, &feedback.Type, &feedback.Name, &feedback.Phone, &feedback.Content, &feedback.UserID, &feedback.CreatedAt, &feedback.UpdatedAt)
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

func (s *Store) GetFeedbackByID(ctx context.Context, id int) (*models.Feedback, error) {
	feedback := &models.Feedback{}
	err := s.querier(ctx).QueryRow(
		ctx,
		"SELECT id, uuid, status, source, type, name, phone, content, user_id, created_at, updated_at FROM feedback WHERE id = $1",
		id,
	).Scan(
		&feedback.ID, &feedback.UUID, &feedback.Status, &feedback.Source, &feedback.Type, &feedback.Name, &feedback.Phone, &feedback.Content, &feedback.UserID, &feedback.CreatedAt, &feedback.UpdatedAt,
	)
	return feedback, wrapDBError(err)
}

func (s *Store) GetFeedbackByUUID(ctx context.Context, feedbackUUID uuid.UUID) (*models.Feedback, error) {
	feedback := &models.Feedback{}
	err := s.querier(ctx).QueryRow(
		ctx,
		"SELECT id, uuid, status, source, type, name, phone, content, user_id, created_at, updated_at FROM feedback WHERE uuid = $1",
		feedbackUUID,
	).Scan(
		&feedback.ID, &feedback.UUID, &feedback.Status, &feedback.Source, &feedback.Type, &feedback.Name, &feedback.Phone, &feedback.Content, &feedback.UserID, &feedback.CreatedAt, &feedback.UpdatedAt,
	)
	if err != nil {
		return nil, wrapDBError(err)
	}

	return feedback, nil
}

func (s *Store) CreateFeedback(ctx context.Context, feedback *models.Feedback) error {
	err := s.querier(ctx).QueryRow(
		ctx,
		"INSERT INTO feedback (uuid, status, source, type, name, phone, content, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id",
		feedback.UUID, feedback.Status, feedback.Source, feedback.Type, feedback.Name, feedback.Phone, feedback.Content, feedback.UserID, feedback.CreatedAt, feedback.UpdatedAt,
	).Scan(&feedback.ID)
	return wrapDBError(err)
}

func (s *Store) UpdateFeedbackStatus(ctx context.Context, feedback *models.Feedback) error {
	_, err := s.querier(ctx).Exec(
		ctx,
		"UPDATE feedback SET status = $2, updated_at = $3 WHERE id = $1",
		feedback.ID, feedback.Status, feedback.UpdatedAt,
	)
	return wrapDBError(err)
}
