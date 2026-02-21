package store

import (
	"context"

	"github.com/google/uuid"
	"github.com/zagvozdeen/ola/internal/store/models"
)

func (s *Store) GetAllReviews(ctx context.Context) ([]models.Review, error) {
	panic("implement")
}

func (s *Store) GetReviewByUUID(ctx context.Context, uuid uuid.UUID) (*models.Review, error) {
	panic("implement")
}

func (s *Store) CreateReview(ctx context.Context, service *models.Review) error {
	panic("implement")
}

func (s *Store) UpdateReview(ctx context.Context, service *models.Review) error {
	panic("implement")
}
