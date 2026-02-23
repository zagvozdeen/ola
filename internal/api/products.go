package api

import (
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
	Name        string `json:"name" mold:"trim" validate:"required,max=255"`
	Description string `json:"description" mold:"trim" validate:"required,max=3000"`
	PriceFrom   int    `json:"price_from" validate:"required,gte=0"`
	PriceTo     *int   `json:"price_to" validate:"omitempty,gte=0"`
	Type        string `json:"type" mold:"trim,lcase" validate:"required,oneof=product service"`
	FileID      int    `json:"file_id" validate:"required,gt=0"`
}

func (s *Service) getProducts(r *http.Request, user *models.User) core.Response {
	products, err := s.store.GetAllProducts(r.Context())
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to get products: %w", err))
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

	pType, err := enums.NewProductType(req.Type)
	if err != nil {
		return core.Err(http.StatusBadRequest, fmt.Errorf("invalid product type: %w", err))
	}

	uid, err := uuid.NewV7()
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to generate uuid v7: %w", err))
	}

	now := time.Now()
	product := &models.Product{
		UUID:        uid,
		Name:        req.Name,
		Description: req.Description,
		PriceFrom:   req.PriceFrom,
		PriceTo:     req.PriceTo,
		Type:        pType,
		FileID:      req.FileID,
		UserID:      user.ID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err = s.store.CreateProduct(r.Context(), product); err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to create product: %w", err))
	}

	created, err := s.store.GetProductByUUID(r.Context(), uid)
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to load created product: %w", err))
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

	pType, err := enums.NewProductType(req.Type)
	if err != nil {
		return core.Err(http.StatusBadRequest, fmt.Errorf("invalid product type: %w", err))
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
	product.Type = pType
	product.FileID = req.FileID
	product.UserID = user.ID
	product.UpdatedAt = time.Now()

	if err = s.store.UpdateProduct(r.Context(), product); err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to update product: %w", err))
	}

	updated, err := s.store.GetProductByUUID(r.Context(), uid)
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to load updated product: %w", err))
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
