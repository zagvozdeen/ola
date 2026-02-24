package store

import (
	"context"

	"github.com/google/uuid"
	"github.com/zagvozdeen/ola/internal/store/models"
)

// GetAllReviews
func (s *Store) GetAllReviews(ctx context.Context) ([]models.Review, error) {
	rows, err := s.pool.Query(ctx, "SELECT r.id, r.uuid, r.name, r.content, r.file_id, f.content, r.user_id, r.published_at, r.created_at, r.updated_at FROM reviews r JOIN files f ON f.id = r.file_id ORDER BY r.published_at DESC")
	if err != nil {
		return nil, wrapDBError(err)
	}
	defer rows.Close()
	reviews := make([]models.Review, 0)
	for rows.Next() {
		var review models.Review
		err = rows.Scan(&review.ID, &review.UUID, &review.Name, &review.Content, &review.FileID, &review.FileContent, &review.UserID, &review.PublishedAt, &review.CreatedAt, &review.UpdatedAt)
		if err != nil {
			return nil, wrapDBError(err)
		}
		reviews = append(reviews, review)
	}
	err = rows.Err()
	if err != nil {
		return nil, wrapDBError(err)
	}
	return reviews, nil
}

func (s *Store) GetReviewByID(ctx context.Context, id int) (*models.Review, error) {
	review := &models.Review{}
	err := s.pool.QueryRow(
		ctx,
		"SELECT id, uuid, name, content, file_id, user_id, published_at, created_at, updated_at FROM reviews r WHERE id = $1",
		id,
	).Scan(&review.ID, &review.UUID, &review.Name, &review.Content, &review.FileID, &review.UserID, &review.PublishedAt, &review.CreatedAt, &review.UpdatedAt)
	if err != nil {
		return nil, wrapDBError(err)
	}
	return review, nil
}

func (s *Store) GetReviewByUUID(ctx context.Context, uuid uuid.UUID) (*models.Review, error) {
	review := &models.Review{}
	fileContent := ""
	err := s.pool.QueryRow(ctx, "SELECT r.id, r.uuid, r.name, r.content, r.file_id, f.content, r.user_id, r.published_at, r.created_at, r.updated_at FROM reviews r JOIN files f ON f.id = r.file_id WHERE r.uuid = $1", uuid).Scan(&review.ID, &review.UUID, &review.Name, &review.Content, &review.FileID, &fileContent, &review.UserID, &review.PublishedAt, &review.CreatedAt, &review.UpdatedAt)
	if err != nil {
		return nil, wrapDBError(err)
	}
	review.FileContent = &fileContent
	return review, nil
}

func (s *Store) CreateReview(ctx context.Context, review *models.Review) error {
	err := s.pool.QueryRow(
		ctx,
		"INSERT INTO reviews (uuid, name, content, file_id, user_id, published_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id",
		review.UUID, review.Name, review.Content, review.FileID, review.UserID, review.PublishedAt, review.CreatedAt, review.UpdatedAt,
	).Scan(&review.ID)
	return wrapDBError(err)
}

func (s *Store) UpdateReview(ctx context.Context, review *models.Review) error {
	tag, err := s.pool.Exec(ctx, "UPDATE reviews SET name = $1, content = $2, file_id = $3, user_id = $4, published_at = $5, updated_at = $6 WHERE id = $7", review.Name, review.Content, review.FileID, review.UserID, review.PublishedAt, review.UpdatedAt, review.ID)
	if err != nil {
		return wrapDBError(err)
	}
	if tag.RowsAffected() == 0 {
		return models.ErrNotFound
	}
	return nil
}

func (s *Store) DeleteReviewByUUID(ctx context.Context, reviewUUID uuid.UUID) error {
	tag, err := s.pool.Exec(ctx, "DELETE FROM reviews WHERE uuid = $1", reviewUUID)
	if err != nil {
		return wrapDBError(err)
	}
	if tag.RowsAffected() == 0 {
		return models.ErrNotFound
	}
	return nil
}
