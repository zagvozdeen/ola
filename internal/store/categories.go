package store

import (
	"context"

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
