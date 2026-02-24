package store

import (
	"context"

	"github.com/google/uuid"
	"github.com/zagvozdeen/ola/internal/store/enums"
	"github.com/zagvozdeen/ola/internal/store/models"
)

func (s *Store) GetAllUsers(ctx context.Context) ([]models.User, error) {
	rows, err := s.pool.Query(ctx, "SELECT id, tid, uuid, first_name, last_name, username, email, phone, password, role, created_at, updated_at FROM users ORDER BY created_at DESC")
	if err != nil {
		return nil, wrapDBError(err)
	}
	defer rows.Close()

	users := make([]models.User, 0)
	for rows.Next() {
		user := models.User{}
		err = rows.Scan(&user.ID, &user.TID, &user.UUID, &user.FirstName, &user.LastName, &user.Username, &user.Email, &user.Phone, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)
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
	err := s.pool.QueryRow(ctx, "SELECT id, tid, uuid, first_name, last_name, username, email, phone, password, role, created_at, updated_at FROM users WHERE id = $1", id).Scan(&user.ID, &user.TID, &user.UUID, &user.FirstName, &user.LastName, &user.Username, &user.Email, &user.Phone, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, wrapDBError(err)
	}
	return user, nil
}

func (s *Store) GetUserByTID(ctx context.Context, tid int64) (*models.User, error) {
	user := &models.User{}
	err := s.pool.QueryRow(ctx, "SELECT id, tid, uuid, first_name, last_name, username, email, phone, password, role, created_at, updated_at FROM users WHERE tid = $1", tid).Scan(&user.ID, &user.TID, &user.UUID, &user.FirstName, &user.LastName, &user.Username, &user.Email, &user.Phone, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, wrapDBError(err)
	}
	return user, nil
}

func (s *Store) GetUserByUUID(ctx context.Context, userUUID uuid.UUID) (*models.User, error) {
	user := &models.User{}
	err := s.pool.QueryRow(
		ctx,
		"SELECT id, tid, uuid, first_name, last_name, username, email, phone, password, role, created_at, updated_at FROM users WHERE uuid = $1",
		userUUID,
	).Scan(&user.ID, &user.TID, &user.UUID, &user.FirstName, &user.LastName, &user.Username, &user.Email, &user.Phone, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, wrapDBError(err)
	}
	return user, nil
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	err := s.pool.QueryRow(
		ctx,
		"SELECT id, tid, uuid, first_name, last_name, username, email, phone, password, role, created_at, updated_at FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.TID, &user.UUID, &user.FirstName, &user.LastName, &user.Username, &user.Email, &user.Phone, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, wrapDBError(err)
	}
	return user, nil
}

func (s *Store) CreateUser(ctx context.Context, user *models.User) error {
	err := s.pool.QueryRow(
		ctx,
		"INSERT INTO users (tid, uuid, first_name, last_name, username, email, phone, password, role, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id",
		user.TID, user.UUID, user.FirstName, user.LastName, user.Username, user.Email, user.Phone, user.Password, user.Role, user.CreatedAt, user.UpdatedAt,
	).Scan(&user.ID)
	return wrapDBError(err)
}

func (s *Store) UpdateUserPhone(ctx context.Context, user *models.User) error {
	_, err := s.pool.Exec(
		ctx,
		"UPDATE users SET phone = $1, updated_at = $2 WHERE id = $3",
		user.Phone, user.UpdatedAt, user.ID,
	)
	return wrapDBError(err)
}

func (s *Store) UpdateUserRole(ctx context.Context, userID int, role enums.UserRole) error {
	tag, err := s.pool.Exec(
		ctx,
		"UPDATE users SET role = $1, updated_at = NOW() WHERE id = $2",
		role, userID,
	)
	if err != nil {
		return wrapDBError(err)
	}
	if tag.RowsAffected() == 0 {
		return models.ErrNotFound
	}
	return nil
}

func (s *Store) CountUsersByRole(ctx context.Context, role enums.UserRole) (int, error) {
	var count int
	err := s.pool.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE role = $1", role).Scan(&count)
	if err != nil {
		return 0, wrapDBError(err)
	}
	return count, nil
}
