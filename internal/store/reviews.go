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
		return nil, err
	}
	defer rows.Close()
	reviews := make([]models.Review, 0)
	for rows.Next() {
		var review models.Review
		err = rows.Scan(&review.ID, &review.UUID, &review.Name, &review.Content, &review.FileID, &review.FileContent, &review.UserID, &review.PublishedAt, &review.CreatedAt, &review.UpdatedAt)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return reviews, nil
}

func (s *Store) GetReviewByUUID(ctx context.Context, uuid uuid.UUID) (*models.Review, error) {
	review := &models.Review{}
	fileContent := ""
	err := s.pool.QueryRow(ctx, "SELECT r.id, r.uuid::text, r.name, r.content, r.file_id, f.content, r.user_id, r.published_at, r.created_at, r.updated_at FROM reviews r JOIN files f ON f.id = r.file_id WHERE r.uuid = $1", uuid.String()).Scan(&review.ID, &review.UUID, &review.Name, &review.Content, &review.FileID, &fileContent, &review.UserID, &review.PublishedAt, &review.CreatedAt, &review.UpdatedAt)
	if err != nil {
		return nil, err
	}
	review.FileContent = &fileContent
	return review, nil
}

func (s *Store) CreateReview(ctx context.Context, review *models.Review) error {
	_, err := s.pool.Exec(ctx, "INSERT INTO reviews (id, uuid, name, content, file_id, user_id, published_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)", review.ID, review.UUID, review.Name, review.Content, review.FileID, review.UserID, review.PublishedAt, review.CreatedAt, review.UpdatedAt)
	return err
}

func (s *Store) UpdateReview(ctx context.Context, review *models.Review) error {
	_, err := s.pool.Exec(ctx, "UPDATE reviews SET name = $1, content = $2, file_id = $3, user_id = $4, published_at = $5, updated_at = $6 WHERE id = $7", review.Name, review.Content, review.FileID, review.UserID, review.PublishedAt, review.UpdatedAt, review.ID)
	return err
}
