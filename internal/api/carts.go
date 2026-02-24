package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/zagvozdeen/ola/internal/api/core"
	"github.com/zagvozdeen/ola/internal/store/models"
)

type upsertCartItemRequest struct {
	ProductID int `json:"product_id" validate:"required,gt=0"`
	Qty       int `json:"qty" validate:"required,gt=0"`
}

func (s *Service) getCart(r *http.Request, user *models.User) core.Response {
	items, err := s.store.GetUserCartItems(r.Context(), user.ID)
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to get cart: %w", err))
	}

	return core.JSON(http.StatusOK, items)
}

func (s *Service) upsertCartItem(r *http.Request, user *models.User) core.Response {
	req, res := core.Validate[upsertCartItemRequest](r, s.conform, s.validate)
	if res != nil {
		return res
	}

	err := s.store.UpsertUserCartItem(r.Context(), user.ID, req.ProductID, req.Qty)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return core.Err(http.StatusNotFound, fmt.Errorf("product not found"))
		}

		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to upsert cart item: %w", err))
	}

	return core.JSON(http.StatusNoContent, nil)
}

func (s *Service) deleteCartItem(r *http.Request, user *models.User) core.Response {
	productUUID, err := uuid.Parse(r.PathValue("product_uuid"))
	if err != nil {
		return core.Err(http.StatusBadRequest, fmt.Errorf("invalid product uuid"))
	}

	err = s.store.DeleteUserCartItem(r.Context(), user.ID, productUUID)
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to delete cart item: %w", err))
	}

	return core.JSON(http.StatusNoContent, nil)
}
