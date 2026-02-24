package store

import (
	"context"

	"github.com/google/uuid"
	"github.com/zagvozdeen/ola/internal/store/models"
)

// GetAllCategories
func (s *Store) GetAllCategories(ctx context.Context) ([]models.Category, error) {
	rows, err := s.pool.Query(ctx, "SELECT id, uuid, name, created_at, updated_at FROM categories ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	categories := make([]models.Category, 0)
	for rows.Next() {
		var category models.Category
		err = rows.Scan(&category.ID, &category.UUID, &category.Name, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (s *Store) GetCategoryByUUID(ctx context.Context, categoryUUID uuid.UUID) (*models.Category, error) {
	category := &models.Category{}
	err := s.pool.QueryRow(
		ctx,
		"SELECT id, uuid, name, created_at, updated_at FROM categories WHERE uuid = $1",
		categoryUUID,
	).Scan(&category.ID, &category.UUID, &category.Name, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		return nil, wrapDBError(err)
	}

	return category, nil
}

func (s *Store) CreateCategory(ctx context.Context, category *models.Category) error {
	err := s.pool.QueryRow(
		ctx,
		"INSERT INTO categories (uuid, name, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id",
		category.UUID, category.Name, category.CreatedAt, category.UpdatedAt,
	).Scan(&category.ID)
	return wrapDBError(err)
}

func (s *Store) UpdateCategory(ctx context.Context, category *models.Category) error {
	tag, err := s.pool.Exec(
		ctx,
		"UPDATE categories SET name = $1, updated_at = $2 WHERE id = $3",
		category.Name, category.UpdatedAt, category.ID,
	)
	if err != nil {
		return wrapDBError(err)
	}
	if tag.RowsAffected() == 0 {
		return models.ErrNotFound
	}
	return nil
}

func (s *Store) DeleteCategoryByUUID(ctx context.Context, categoryUUID uuid.UUID) error {
	tag, err := s.pool.Exec(ctx, "DELETE FROM categories WHERE uuid = $1", categoryUUID)
	if err != nil {
		return wrapDBError(err)
	}
	if tag.RowsAffected() == 0 {
		return models.ErrNotFound
	}
	return nil
}
