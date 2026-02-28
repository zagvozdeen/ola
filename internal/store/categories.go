package store

import (
	"context"
	"fmt"
	"strings"

	"github.com/zagvozdeen/ola/internal/store/models"
)

// GetAllCategories
func (s *Store) GetAllCategories(ctx context.Context) ([]models.Category, error) {
	rows, err := s.querier(ctx).Query(ctx, "SELECT id, slug, name, created_at, updated_at FROM categories ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	categories := make([]models.Category, 0)
	for rows.Next() {
		var category models.Category
		err = rows.Scan(&category.ID, &category.Slug, &category.Name, &category.CreatedAt, &category.UpdatedAt)
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

func (s *Store) GetCategoryByName(ctx context.Context, name string) (*models.Category, error) {
	category := &models.Category{}
	err := s.querier(ctx).QueryRow(
		ctx,
		"SELECT id, slug, name, created_at, updated_at FROM categories WHERE name = $1",
		name,
	).Scan(&category.ID, &category.Slug, &category.Name, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		return nil, wrapDBError(err)
	}
	return category, nil
}

func (s *Store) GetCategoryBySlug(ctx context.Context, categorySlug string) (*models.Category, error) {
	category := &models.Category{}
	err := s.querier(ctx).QueryRow(
		ctx,
		"SELECT id, slug, name, created_at, updated_at FROM categories WHERE slug = $1",
		categorySlug,
	).Scan(&category.ID, &category.Slug, &category.Name, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		return nil, wrapDBError(err)
	}

	return category, nil
}

func (s *Store) CreateCategory(ctx context.Context, category *models.Category) error {
	err := s.querier(ctx).QueryRow(
		ctx,
		"INSERT INTO categories (slug, name, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id",
		category.Slug, category.Name, category.CreatedAt, category.UpdatedAt,
	).Scan(&category.ID)
	return wrapDBError(err)
}

func (s *Store) UpdateCategory(ctx context.Context, category *models.Category) error {
	tag, err := s.querier(ctx).Exec(
		ctx,
		"UPDATE categories SET slug = $1, name = $2, updated_at = $3 WHERE id = $4",
		category.Slug, category.Name, category.UpdatedAt, category.ID,
	)
	if err != nil {
		return wrapDBError(err)
	}
	if tag.RowsAffected() == 0 {
		return models.ErrNotFound
	}
	return nil
}

func (s *Store) DeleteCategoryBySlug(ctx context.Context, categorySlug string) error {
	_, err := s.querier(ctx).Exec(
		ctx,
		"DELETE FROM category_product WHERE category_id = (SELECT id FROM categories WHERE slug = $1)",
		categorySlug,
	)
	if err != nil {
		return wrapDBError(err)
	}

	tag, err := s.querier(ctx).Exec(ctx, "DELETE FROM categories WHERE slug = $1", categorySlug)
	if err != nil {
		return wrapDBError(err)
	}
	if tag.RowsAffected() == 0 {
		return models.ErrNotFound
	}
	return nil
}

func (s *Store) GetCategoriesByProductIDs(ctx context.Context, productIDs []int) (map[int][]models.Category, error) {
	categoriesByProductID := make(map[int][]models.Category)
	if len(productIDs) == 0 {
		return categoriesByProductID, nil
	}

	placeholders := make([]string, 0, len(productIDs))
	args := make([]any, 0, len(productIDs))
	for _, productID := range productIDs {
		args = append(args, productID)
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(args)))
	}

	rows, err := s.querier(ctx).Query(
		ctx,
		"SELECT cp.product_id, c.id, c.slug, c.name, c.created_at, c.updated_at FROM category_product cp JOIN categories c ON c.id = cp.category_id WHERE cp.product_id IN ("+strings.Join(placeholders, ", ")+") ORDER BY c.name",
		args...,
	)
	if err != nil {
		return nil, wrapDBError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			productID int
			category  models.Category
		)
		err = rows.Scan(&productID, &category.ID, &category.Slug, &category.Name, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			return nil, wrapDBError(err)
		}
		categoriesByProductID[productID] = append(categoriesByProductID[productID], category)
	}

	err = rows.Err()
	if err != nil {
		return nil, wrapDBError(err)
	}

	return categoriesByProductID, nil
}

func (s *Store) ReplaceProductCategories(ctx context.Context, productID int, categoryIDs []int) error {
	_, err := s.querier(ctx).Exec(ctx, "DELETE FROM category_product WHERE product_id = $1", productID)
	if err != nil {
		return wrapDBError(err)
	}

	seen := make(map[int]struct{}, len(categoryIDs))
	for _, categoryID := range categoryIDs {
		if _, ok := seen[categoryID]; ok {
			continue
		}
		seen[categoryID] = struct{}{}

		_, err = s.querier(ctx).Exec(
			ctx,
			"INSERT INTO category_product (category_id, product_id) VALUES ($1, $2)",
			categoryID,
			productID,
		)
		if err != nil {
			return wrapDBError(err)
		}
	}

	return nil
}
