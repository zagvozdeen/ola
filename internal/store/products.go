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
		return nil, err
	}
	defer rows.Close()
	products := make([]models.Product, 0)
	for rows.Next() {
		var product models.Product
		err = rows.Scan(&product.ID, &product.UUID, &product.Name, &product.Description, &product.PriceFrom, &product.PriceTo, &product.Type, &product.FileID, &product.FileContent, &product.UserID, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *Store) GetProductByUUID(ctx context.Context, uuid uuid.UUID) (*models.Product, error) {
	product := &models.Product{}
	fileContent := ""
	err := s.pool.QueryRow(ctx, "SELECT p.id, p.uuid::text, p.name, p.description, p.price_from, p.price_to, p.type, p.file_id, f.content, p.user_id, p.created_at, p.updated_at FROM products p JOIN files f ON f.id = p.file_id WHERE p.uuid = $1", uuid.String()).Scan(&product.ID, &product.UUID, &product.Name, &product.Description, &product.PriceFrom, &product.PriceTo, &product.Type, &product.FileID, &fileContent, &product.UserID, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		return nil, err
	}
	product.FileContent = &fileContent
	return product, nil
}

func (s *Store) CreateProduct(ctx context.Context, product *models.Product) error {
	_, err := s.pool.Exec(ctx, "INSERT INTO products (id, uuid, name, description, price_from, price_to, type, file_id, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)", product.ID, product.UUID, product.Name, product.Description, product.PriceFrom, product.PriceTo, product.Type, product.FileID, product.UserID, product.CreatedAt, product.UpdatedAt)
	return err
}

func (s *Store) UpdateProduct(ctx context.Context, product *models.Product) error {
	_, err := s.pool.Exec(ctx, "UPDATE products SET name = $1, description = $2, price_from = $3, price_to = $4, type = $5, file_id = $6, user_id = $7, updated_at = $8 WHERE id = $9", product.Name, product.Description, product.PriceFrom, product.PriceTo, product.Type, product.FileID, product.UserID, product.UpdatedAt, product.ID)
	return err
}
