package store

import (
	"context"

	"github.com/google/uuid"
	"github.com/zagvozdeen/ola/internal/store/models"
)

// GetAllProducts
func (s *Store) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	rows, err := s.pool.Query(ctx, "SELECT p.id, p.uuid, p.name, p.description, p.price_from, p.price_to, p.type, p.file_id, f.content, p.user_id, p.created_at, p.updated_at FROM products p JOIN files f ON f.id = p.file_id ORDER BY p.created_at DESC")
	if err != nil {
		return nil, wrapDBError(err)
	}
	defer rows.Close()
	products := make([]models.Product, 0)
	for rows.Next() {
		var product models.Product
		err = rows.Scan(&product.ID, &product.UUID, &product.Name, &product.Description, &product.PriceFrom, &product.PriceTo, &product.Type, &product.FileID, &product.FileContent, &product.UserID, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return nil, wrapDBError(err)
		}
		products = append(products, product)
	}
	err = rows.Err()
	if err != nil {
		return nil, wrapDBError(err)
	}
	return products, nil
}

func (s *Store) GetProductByID(ctx context.Context, id int) (*models.Product, error) {
	product := &models.Product{}
	err := s.pool.QueryRow(
		ctx,
		"SELECT id, uuid, name, description, price_from, price_to, type, file_id, user_id, created_at, updated_at FROM products p WHERE id = $1",
		id,
	).Scan(&product.ID, &product.UUID, &product.Name, &product.Description, &product.PriceFrom, &product.PriceTo, &product.Type, &product.FileID, &product.UserID, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		return nil, wrapDBError(err)
	}
	return product, nil
}

func (s *Store) GetProductByUUID(ctx context.Context, uuid uuid.UUID) (*models.Product, error) {
	product := &models.Product{}
	err := s.pool.QueryRow(
		ctx,
		"SELECT p.id, p.uuid, p.name, p.description, p.price_from, p.price_to, p.type, p.file_id, f.content, p.user_id, p.created_at, p.updated_at FROM products p JOIN files f ON f.id = p.file_id WHERE p.uuid = $1",
		uuid,
	).Scan(&product.ID, &product.UUID, &product.Name, &product.Description, &product.PriceFrom, &product.PriceTo, &product.Type, &product.FileID, &product.FileContent, &product.UserID, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		return nil, wrapDBError(err)
	}
	return product, nil
}

func (s *Store) CreateProduct(ctx context.Context, product *models.Product) error {
	err := s.pool.QueryRow(
		ctx,
		"INSERT INTO products (uuid, name, description, price_from, price_to, type, file_id, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id",
		product.UUID, product.Name, product.Description, product.PriceFrom, product.PriceTo, product.Type, product.FileID, product.UserID, product.CreatedAt, product.UpdatedAt,
	).Scan(&product.ID)
	return wrapDBError(err)
}

func (s *Store) UpdateProduct(ctx context.Context, product *models.Product) error {
	_, err := s.pool.Exec(ctx, "UPDATE products SET name = $1, description = $2, price_from = $3, price_to = $4, type = $5, file_id = $6, user_id = $7, updated_at = $8 WHERE id = $9", product.Name, product.Description, product.PriceFrom, product.PriceTo, product.Type, product.FileID, product.UserID, product.UpdatedAt, product.ID)
	return wrapDBError(err)
}

func (s *Store) DeleteProductByUUID(ctx context.Context, uuid uuid.UUID) error {
	tag, err := s.pool.Exec(ctx, "DELETE FROM products WHERE uuid = $1", uuid)
	if err != nil {
		return wrapDBError(err)
	}
	if tag.RowsAffected() == 0 {
		return models.ErrNotFound
	}
	return nil
}
