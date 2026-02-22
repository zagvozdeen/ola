package api

import (
	"encoding/json/v2"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/zagvozdeen/ola/internal/api/core"
	"github.com/zagvozdeen/ola/internal/store/models"
)

type createOrderRequest struct {
	Name    string `json:"name" mold:"trim" validate:"required,max=255"`
	Phone   string `json:"phone" mold:"trim" validate:"required,max=255,ru_phone"`
	Content string `json:"content" mold:"trim" validate:"required,max=3000"`
}

func (s *Service) getOrders(r *http.Request, user *models.User) core.Response {
	res := allowForModeratorOrAdmin(user)
	if res != nil {
		return res
	}

	orders, err := s.store.GetAllOrders(r.Context())
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to get orders: %w", err))
	}
	return core.JSON(http.StatusOK, orders)
}

func (s *Service) createOrder(r *http.Request, user *models.User) core.Response {
	req := &createOrderRequest{}
	err := json.UnmarshalRead(r.Body, req)
	if err != nil {
		return core.Err(http.StatusBadRequest, err)
	}
	err = s.conform.Struct(r.Context(), req)
	if err != nil {
		return core.Err(http.StatusBadRequest, err)
	}
	err = s.validate.StructCtx(r.Context(), req)
	if err != nil {
		return core.Err(http.StatusBadRequest, err)
	}

	uid, err := uuid.NewV7()
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to generate uuid v7: %w", err))
	}
	order := &models.Order{
		UUID:      uid.String(),
		Name:      req.Name,
		Phone:     req.Phone,
		Content:   req.Content,
		UserID:    &user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = s.store.CreateOrder(r.Context(), order)
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to create order: %w", err))
	}

	return core.JSON(http.StatusCreated, order)
}

type createGuestOrderRequest struct {
	Name    string `json:"name" mold:"trim" validate:"required,max=255"`
	Phone   string `json:"phone" mold:"trim" validate:"required,max=255,ru_phone"`
	Content string `json:"content" mold:"trim" validate:"required,max=3000"`
}

func (s *Service) createGuestOrder(r *http.Request) core.Response {
	req, res := core.Validate[createGuestOrderRequest](r, s.conform, s.validate)
	if res != nil {
		return res
	}

	uid, err := uuid.NewV7()
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to generate uuid v7: %w", err))
	}
	order := &models.Order{
		UUID:      uid.String(),
		Name:      req.Name,
		Phone:     req.Phone,
		Content:   req.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = s.store.CreateOrder(r.Context(), order)
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to create guest order: %w", err))
	}

	return core.JSON(http.StatusCreated, order)
}
