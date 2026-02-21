package store

import (
	"context"

	"github.com/zagvozdeen/ola/internal/store/models"
)

func (s *Store) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	panic("implement")
}

func (s *Store) GetUserByTID(ctx context.Context, tid int64) (*models.User, error) {
	panic("implement")
}

func (s *Store) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	panic("implement")
}

func (s *Store) CreateUser(ctx context.Context, user *models.User) error {
	panic("implement")
}
