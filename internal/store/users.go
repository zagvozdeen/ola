package store

import (
	"context"

	"github.com/zagvozdeen/ola/internal/store/models"
)

func (s *Store) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	user := &models.User{}
	err := s.pool.QueryRow(ctx, "SELECT id, tid, uuid::text, first_name, last_name, username, email, password, role, created_at, updated_at FROM users WHERE id = $1", id).Scan(&user.ID, &user.TID, &user.UUID, &user.FirstName, &user.LastName, &user.Username, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Store) GetUserByTID(ctx context.Context, tid int64) (*models.User, error) {
	user := &models.User{}
	err := s.pool.QueryRow(ctx, "SELECT id, tid, uuid::text, first_name, last_name, username, email, password, role, created_at, updated_at FROM users WHERE tid = $1", tid).Scan(&user.ID, &user.TID, &user.UUID, &user.FirstName, &user.LastName, &user.Username, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Store) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user := &models.User{}
	err := s.pool.QueryRow(ctx, "SELECT id, tid, uuid::text, first_name, last_name, username, email, password, role, created_at, updated_at FROM users WHERE username = $1", username).Scan(&user.ID, &user.TID, &user.UUID, &user.FirstName, &user.LastName, &user.Username, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Store) CreateUser(ctx context.Context, user *models.User) error {
	_, err := s.pool.Exec(ctx, "INSERT INTO users (id, tid, uuid, first_name, last_name, username, email, password, role, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)", user.ID, user.TID, user.UUID, user.FirstName, user.LastName, user.Username, user.Email, user.Password, user.Role, user.CreatedAt, user.UpdatedAt)
	return err
}
