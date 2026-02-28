package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/zagvozdeen/ola/internal/api/core"
	"github.com/zagvozdeen/ola/internal/store/enums"
	"github.com/zagvozdeen/ola/internal/store/models"
)

type upsertProductRequest struct {
	Name          string            `json:"name" mold:"trim" validate:"required,max=255"`
	Description   string            `json:"description" mold:"trim" validate:"required,max=3000"`
	PriceFrom     int               `json:"price_from" validate:"required,gte=0"`
	PriceTo       *int              `json:"price_to" validate:"omitempty,gte=0"`
	Type          enums.ProductType `json:"type"`
	FileContent   string            `json:"file_content" validate:"required"`
	CategoryUUIDs []uuid.UUID       `json:"category_uuids"`
}

func (s *Service) getProducts(r *http.Request, user *models.User) core.Response {
	products, err := s.store.GetAllProducts(r.Context())
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to get products: %w", err))
	}
	if err = s.attachProductCategories(r.Context(), products); err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to get product categories: %w", err))
	}
	return core.JSON(http.StatusOK, products)
}

func (s *Service) getProduct(r *http.Request, user *models.User) core.Response {
	res := allowForModeratorOrAdmin(user)
	if res != nil {
		return res
	}

	uid, err := uuid.Parse(r.PathValue("uuid"))
	if err != nil {
		return core.Err(http.StatusBadRequest, fmt.Errorf("invalid product uuid: %w", err))
	}

	product, err := s.store.GetProductByUUID(r.Context(), uid)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return core.Err(http.StatusNotFound, fmt.Errorf("product not found"))
		}
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to get product: %w", err))
	}

	categoriesByProductID, err := s.store.GetCategoriesByProductIDs(r.Context(), []int{product.ID})
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to get product categories: %w", err))
	}
	product.Categories = categoriesByProductID[product.ID]
	if product.Categories == nil {
		product.Categories = []models.Category{}
	}

	return core.JSON(http.StatusOK, product)
}

func (s *Service) createProduct(r *http.Request, user *models.User) core.Response {
	res := allowForModeratorOrAdmin(user)
	if res != nil {
		return res
	}

	req, res := core.Validate[upsertProductRequest](r, s.conform, s.validate)
	if res != nil {
		return res
	}
	if req.PriceTo != nil && *req.PriceTo < req.PriceFrom {
		return core.Err(http.StatusBadRequest, fmt.Errorf("price_to must be greater than or equal to price_from"))
	}

	//pType, err := enums.NewProductType(req.Type)
	//if err != nil {
	//	return core.Err(http.StatusBadRequest, fmt.Errorf("invalid product type: %w", err))
	//}

	uid, err := uuid.NewV7()
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to generate uuid v7: %w", err))
	}

	categories, err := s.resolveProductCategories(r.Context(), req.CategoryUUIDs)
	if err != nil {
		return core.Err(http.StatusBadRequest, err)
	}

	ctx, err := s.store.Begin(r.Context())
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to begin transaction: %w", err))
	}
	defer s.store.Rollback(ctx)

	now := time.Now()
	product := &models.Product{
		UUID:        uid,
		Name:        req.Name,
		Description: req.Description,
		PriceFrom:   req.PriceFrom,
		PriceTo:     req.PriceTo,
		Type:        req.Type,
		FileContent: req.FileContent,
		UserID:      user.ID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err = s.store.CreateProduct(ctx, product); err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to create product: %w", err))
	}

	categoryIDs := make([]int, 0, len(categories))
	for _, category := range categories {
		categoryIDs = append(categoryIDs, category.ID)
	}
	if err = s.store.ReplaceProductCategories(ctx, product.ID, categoryIDs); err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to save product categories: %w", err))
	}

	s.store.Commit(ctx)

	created, err := s.store.GetProductByUUID(r.Context(), uid)
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to load created product: %w", err))
	}
	created.Categories = categories
	if created.Categories == nil {
		created.Categories = []models.Category{}
	}

	return core.JSON(http.StatusCreated, created)
}

func (s *Service) updateProduct(r *http.Request, user *models.User) core.Response {
	res := allowForModeratorOrAdmin(user)
	if res != nil {
		return res
	}

	uid, err := uuid.Parse(r.PathValue("uuid"))
	if err != nil {
		return core.Err(http.StatusBadRequest, fmt.Errorf("invalid product uuid: %w", err))
	}

	req, res := core.Validate[upsertProductRequest](r, s.conform, s.validate)
	if res != nil {
		return res
	}
	if req.PriceTo != nil && *req.PriceTo < req.PriceFrom {
		return core.Err(http.StatusBadRequest, fmt.Errorf("price_to must be greater than or equal to price_from"))
	}

	//pType, err := enums.NewProductType(req.Type)
	//if err != nil {
	//	return core.Err(http.StatusBadRequest, fmt.Errorf("invalid product type: %w", err))
	//}

	categories, err := s.resolveProductCategories(r.Context(), req.CategoryUUIDs)
	if err != nil {
		return core.Err(http.StatusBadRequest, err)
	}

	product, err := s.store.GetProductByUUID(r.Context(), uid)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return core.Err(http.StatusNotFound, fmt.Errorf("product not found"))
		}
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to get product: %w", err))
	}

	product.Name = req.Name
	product.Description = req.Description
	product.PriceFrom = req.PriceFrom
	product.PriceTo = req.PriceTo
	product.Type = req.Type
	product.FileContent = req.FileContent
	product.UserID = user.ID
	product.UpdatedAt = time.Now()

	ctx, err := s.store.Begin(r.Context())
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to begin transaction: %w", err))
	}
	defer s.store.Rollback(ctx)

	if err = s.store.UpdateProduct(ctx, product); err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to update product: %w", err))
	}

	categoryIDs := make([]int, 0, len(categories))
	for _, category := range categories {
		categoryIDs = append(categoryIDs, category.ID)
	}
	if err = s.store.ReplaceProductCategories(ctx, product.ID, categoryIDs); err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to save product categories: %w", err))
	}

	s.store.Commit(ctx)

	updated, err := s.store.GetProductByUUID(r.Context(), uid)
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to load updated product: %w", err))
	}
	updated.Categories = categories
	if updated.Categories == nil {
		updated.Categories = []models.Category{}
	}
	return core.JSON(http.StatusOK, updated)
}

func (s *Service) deleteProduct(r *http.Request, user *models.User) core.Response {
	res := allowForModeratorOrAdmin(user)
	if res != nil {
		return res
	}

	uid, err := uuid.Parse(r.PathValue("uuid"))
	if err != nil {
		return core.Err(http.StatusBadRequest, fmt.Errorf("invalid product uuid: %w", err))
	}

	if err = s.store.DeleteProductByUUID(r.Context(), uid); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return core.Err(http.StatusNotFound, fmt.Errorf("product not found"))
		}
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to delete product: %w", err))
	}
	return core.JSON(http.StatusNoContent, nil)
}

func (s *Service) attachProductCategories(ctx context.Context, products []models.Product) error {
	if len(products) == 0 {
		return nil
	}

	productIDs := make([]int, 0, len(products))
	for _, product := range products {
		productIDs = append(productIDs, product.ID)
	}

	categoriesByProductID, err := s.store.GetCategoriesByProductIDs(ctx, productIDs)
	if err != nil {
		return err
	}

	for i := range products {
		products[i].Categories = categoriesByProductID[products[i].ID]
		if products[i].Categories == nil {
			products[i].Categories = []models.Category{}
		}
	}

	return nil
}

func (s *Service) resolveProductCategories(ctx context.Context, categoryUUIDs []uuid.UUID) ([]models.Category, error) {
	if len(categoryUUIDs) == 0 {
		return []models.Category{}, nil
	}

	categories := make([]models.Category, 0, len(categoryUUIDs))
	seen := make(map[uuid.UUID]struct{}, len(categoryUUIDs))
	for _, categoryUUID := range categoryUUIDs {
		if _, ok := seen[categoryUUID]; ok {
			continue
		}
		seen[categoryUUID] = struct{}{}

		category, err := s.store.GetCategoryByUUID(ctx, categoryUUID)
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				return nil, fmt.Errorf("category not found: %s", categoryUUID.String())
			}
			return nil, fmt.Errorf("failed to get category: %w", err)
		}
		categories = append(categories, *category)
	}

	return categories, nil
}
