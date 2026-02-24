package store

import (
	"context"

	"github.com/zagvozdeen/ola/internal/store/models"
)

func (s *Store) GetAllUsers(ctx context.Context) ([]models.User, error) {
	rows, err := s.pool.Query(ctx, "SELECT id, tid, uuid, first_name, last_name, username, email, password, role, created_at, updated_at FROM users ORDER BY created_at DESC")
	if err != nil {
		return nil, wrapDBError(err)
	}
	defer rows.Close()

	users := make([]models.User, 0)
	for rows.Next() {
		user := models.User{}
		err = rows.Scan(&user.ID, &user.TID, &user.UUID, &user.FirstName, &user.LastName, &user.Username, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, wrapDBError(err)
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, wrapDBError(err)
	}

	return users, nil
}

func (s *Store) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	user := &models.User{}
	err := s.pool.QueryRow(ctx, "SELECT id, tid, uuid, first_name, last_name, username, email, password, role, created_at, updated_at FROM users WHERE id = $1", id).Scan(&user.ID, &user.TID, &user.UUID, &user.FirstName, &user.LastName, &user.Username, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, wrapDBError(err)
	}
	return user, nil
}

func (s *Store) GetUserByTID(ctx context.Context, tid int64) (*models.User, error) {
	user := &models.User{}
	err := s.pool.QueryRow(ctx, "SELECT id, tid, uuid, first_name, last_name, username, email, password, role, created_at, updated_at FROM users WHERE tid = $1", tid).Scan(&user.ID, &user.TID, &user.UUID, &user.FirstName, &user.LastName, &user.Username, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, wrapDBError(err)
	}
	return user, nil
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	err := s.pool.QueryRow(
		ctx,
		"SELECT id, tid, uuid, first_name, last_name, username, email, password, role, created_at, updated_at FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.TID, &user.UUID, &user.FirstName, &user.LastName, &user.Username, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, wrapDBError(err)
	}
	return user, nil
}

func (s *Store) CreateUser(ctx context.Context, user *models.User) error {
	err := s.pool.QueryRow(
		ctx,
		"INSERT INTO users (tid, uuid, first_name, last_name, username, email, password, role, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id",
		user.TID, user.UUID, user.FirstName, user.LastName, user.Username, user.Email, user.Password, user.Role, user.CreatedAt, user.UpdatedAt,
	).Scan(&user.ID)
	return wrapDBError(err)
}
