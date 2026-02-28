package store

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/zagvozdeen/ola/internal/store/enums"
	"github.com/zagvozdeen/ola/internal/store/models"
)

// GetAllProducts
func (s *Store) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	rows, err := s.querier(ctx).Query(ctx, "SELECT id, uuid, name, description, price_from, price_to, type, file_content, user_id, created_at, updated_at FROM products ORDER BY created_at DESC")
	if err != nil {
		return nil, wrapDBError(err)
	}
	defer rows.Close()
	products := make([]models.Product, 0)
	for rows.Next() {
		var product models.Product
		err = rows.Scan(&product.ID, &product.UUID, &product.Name, &product.Description, &product.PriceFrom, &product.PriceTo, &product.Type, &product.FileContent, &product.UserID, &product.CreatedAt, &product.UpdatedAt)
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

// GetMainProducts
func (s *Store) GetMainProducts(ctx context.Context) ([]models.Product, error) {
	rows, err := s.querier(ctx).Query(
		ctx,
		"SELECT name, description, price_from, price_to, file_content FROM products WHERE is_main = TRUE AND type = $1 ORDER BY created_at DESC LIMIT 4",
		enums.ProductTypeProduct,
	)
	if err != nil {
		return nil, wrapDBError(err)
	}
	defer rows.Close()
	products := make([]models.Product, 0)
	for rows.Next() {
		var product models.Product
		err = rows.Scan(&product.Name, &product.Description, &product.PriceFrom, &product.PriceTo, &product.FileContent)
		if err != nil {
			return nil, wrapDBError(err)
		}
		products = append(products, product)
	}
	err = rows.Err()
	return products, wrapDBError(err)
}

// GetMainServices
func (s *Store) GetMainServices(ctx context.Context) ([]models.Product, error) {
	rows, err := s.querier(ctx).Query(
		ctx,
		"SELECT name, description, price_from, price_to, file_content FROM products WHERE is_main = TRUE AND type = $1 ORDER BY created_at DESC LIMIT 4",
		enums.ProductTypeService,
	)
	if err != nil {
		return nil, wrapDBError(err)
	}
	defer rows.Close()
	products := make([]models.Product, 0)
	for rows.Next() {
		var product models.Product
		err = rows.Scan(&product.Name, &product.Description, &product.PriceFrom, &product.PriceTo, &product.FileContent)
		if err != nil {
			return nil, wrapDBError(err)
		}
		products = append(products, product)
	}
	err = rows.Err()
	return products, wrapDBError(err)
}

// GetCatalogProducts
func (s *Store) GetCatalogProducts(ctx context.Context, categorySlugs []string, productType *enums.ProductType) ([]models.Product, error) {
	query := strings.Builder{}
	query.WriteString("SELECT p.id, p.uuid, p.name, p.description, p.price_from, p.price_to, p.type, p.file_content, p.user_id, p.created_at, p.updated_at FROM products p")

	args := make([]any, 0, len(categorySlugs)+1)
	conditions := make([]string, 0, 2)

	if productType != nil {
		args = append(args, *productType)
		conditions = append(conditions, fmt.Sprintf("p.type = $%d", len(args)))
	}

	if len(categorySlugs) > 0 {
		placeholders := make([]string, 0, len(categorySlugs))
		for _, categorySlug := range categorySlugs {
			args = append(args, categorySlug)
			placeholders = append(placeholders, fmt.Sprintf("$%d", len(args)))
		}

		conditions = append(
			conditions,
			"EXISTS (SELECT 1 FROM category_product cp JOIN categories c ON c.id = cp.category_id WHERE cp.product_id = p.id AND c.slug IN ("+strings.Join(placeholders, ", ")+"))",
		)
	}

	if len(conditions) > 0 {
		query.WriteString(" WHERE ")
		query.WriteString(strings.Join(conditions, " AND "))
	}

	query.WriteString(" ORDER BY p.created_at DESC")

	rows, err := s.querier(ctx).Query(ctx, query.String(), args...)
	if err != nil {
		return nil, wrapDBError(err)
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var product models.Product
		err = rows.Scan(&product.ID, &product.UUID, &product.Name, &product.Description, &product.PriceFrom, &product.PriceTo, &product.Type, &product.FileContent, &product.UserID, &product.CreatedAt, &product.UpdatedAt)
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
	err := s.querier(ctx).QueryRow(
		ctx,
		"SELECT id, uuid, name, description, price_from, price_to, type, file_content, user_id, created_at, updated_at FROM products WHERE id = $1",
		id,
	).Scan(&product.ID, &product.UUID, &product.Name, &product.Description, &product.PriceFrom, &product.PriceTo, &product.Type, &product.FileContent, &product.UserID, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		return nil, wrapDBError(err)
	}
	return product, nil
}

func (s *Store) GetProductByUUID(ctx context.Context, uuid uuid.UUID) (*models.Product, error) {
	product := &models.Product{}
	err := s.querier(ctx).QueryRow(
		ctx,
		"SELECT id, uuid, name, description, price_from, price_to, type, is_main, file_content, user_id, created_at, updated_at FROM products WHERE uuid = $1",
		uuid,
	).Scan(&product.ID, &product.UUID, &product.Name, &product.Description, &product.PriceFrom, &product.PriceTo, &product.Type, &product.IsMain, &product.FileContent, &product.UserID, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		return nil, wrapDBError(err)
	}
	return product, nil
}

func (s *Store) CreateProduct(ctx context.Context, product *models.Product) error {
	err := s.querier(ctx).QueryRow(
		ctx,
		"INSERT INTO products (uuid, name, description, price_from, price_to, type, is_main, file_content, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id",
		product.UUID, product.Name, product.Description, product.PriceFrom, product.PriceTo, product.Type, product.IsMain, product.FileContent, product.UserID, product.CreatedAt, product.UpdatedAt,
	).Scan(&product.ID)
	return wrapDBError(err)
}

func (s *Store) UpdateProduct(ctx context.Context, product *models.Product) error {
	_, err := s.querier(ctx).Exec(
		ctx,
		"UPDATE products SET name = $1, description = $2, price_from = $3, price_to = $4, type = $5, is_main = $6, file_content = $7, user_id = $8, updated_at = $9 WHERE id = $10",
		product.Name, product.Description, product.PriceFrom, product.PriceTo, product.Type, product.IsMain, product.FileContent, product.UserID, product.UpdatedAt, product.ID,
	)
	return wrapDBError(err)
}

func (s *Store) DeleteProductByUUID(ctx context.Context, uuid uuid.UUID) error {
	_, err := s.querier(ctx).Exec(
		ctx,
		"DELETE FROM category_product WHERE product_id = (SELECT id FROM products WHERE uuid = $1)",
		uuid,
	)
	if err != nil {
		return wrapDBError(err)
	}

	tag, err := s.querier(ctx).Exec(ctx, "DELETE FROM products WHERE uuid = $1", uuid)
	if err != nil {
		return wrapDBError(err)
	}
	if tag.RowsAffected() == 0 {
		return models.ErrNotFound
	}
	return nil
}
