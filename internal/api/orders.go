package api

import (
	"encoding/json/v2"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/zagvozdeen/ola/internal/api/core"
	"github.com/zagvozdeen/ola/internal/store/enums"
	"github.com/zagvozdeen/ola/internal/store/models"
)

type createOrderRequest struct {
	Name    string `json:"name" mold:"trim" validate:"required,max=255"`
	Phone   string `json:"phone" mold:"trim" validate:"required,max=255,ru_phone"`
	Content string `json:"content" mold:"trim" validate:"required,max=3000"`
}

const orderSourceHeader = "X-App-Source"

func sourceFromAuthHeader(authorization string) enums.OrderSource {
	if strings.HasPrefix(authorization, "tma ") {
		return enums.OrderSourceTMA
	}

	return enums.OrderSourceSPA
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

func (s *Service) getOrder(r *http.Request, user *models.User) core.Response {
	res := allowForModeratorOrAdmin(user)
	if res != nil {
		return res
	}

	uid, err := uuid.Parse(r.PathValue("uuid"))
	if err != nil {
		return core.Err(http.StatusBadRequest, fmt.Errorf("invalid order uuid: %w", err))
	}

	order, err := s.store.GetOrderByUUID(r.Context(), uid)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return core.Err(http.StatusNotFound, fmt.Errorf("order not found"))
		}

		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to get order: %w", err))
	}

	return core.JSON(http.StatusOK, order)
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

	err = s.store.UpdateUserPhone(r.Context(), user.ID, req.Phone)
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to update user phone: %w", err))
	}
	user.Phone = &req.Phone

	uid, err := uuid.NewV7()
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to generate uuid v7: %w", err))
	}

	order := &models.Order{
		UUID:      uid,
		Status:    enums.RequestStatusCreated,
		Source:    sourceFromAuthHeader(r.Header.Get("Authorization")),
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

func (s *Service) createOrderFromCart(r *http.Request, user *models.User) core.Response {
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

	err = s.store.UpdateUserPhone(r.Context(), user.ID, req.Phone)
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to update user phone: %w", err))
	}
	user.Phone = &req.Phone

	order, err := s.store.CreateOrderFromUserCart(
		r.Context(),
		user.ID,
		sourceFromAuthHeader(r.Header.Get("Authorization")),
		req.Name,
		req.Phone,
		req.Content,
	)
	if err != nil {
		if errors.Is(err, models.ErrCartEmpty) {
			return core.Err(http.StatusBadRequest, fmt.Errorf("cart is empty"))
		}

		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to create order from cart: %w", err))
	}

	return core.JSON(http.StatusCreated, order)
}

type createGuestOrderRequest struct {
	Name    string `json:"name" mold:"trim" validate:"required,max=255"`
	Phone   string `json:"phone" mold:"trim" validate:"required,max=255,ru_phone"`
	Content string `json:"content" mold:"trim" validate:"required,max=3000"`
	Consent bool   `json:"consent" validate:"required"`
	Source  string `json:"source" mold:"trim,lcase" validate:"omitempty,oneof=landing spa tma"`
}

func (s *Service) createGuestOrder(r *http.Request) core.Response {
	req, res := core.Validate[createGuestOrderRequest](r, s.conform, s.validate)
	if res != nil {
		return res
	}

	source := enums.OrderSourceLanding

	if req.Source != "" {
		parsedSource, err := enums.NewOrderSource(req.Source)
		if err != nil {
			return core.Err(http.StatusBadRequest, fmt.Errorf("invalid order source: %w", err))
		}
		source = parsedSource
	} else if headerSource := strings.TrimSpace(strings.ToLower(r.Header.Get(orderSourceHeader))); headerSource != "" {
		parsedSource, err := enums.NewOrderSource(headerSource)
		if err != nil {
			return core.Err(http.StatusBadRequest, fmt.Errorf("invalid order source header: %w", err))
		}
		source = parsedSource
	}

	uid, err := uuid.NewV7()
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to generate uuid v7: %w", err))
	}

	order := &models.Order{
		UUID:      uid,
		Status:    enums.RequestStatusCreated,
		Source:    source,
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

type updateOrderStatusRequest struct {
	Status string `json:"status" mold:"trim,lcase" validate:"required,oneof=created in_progress reviewed"`
}

func (s *Service) updateOrderStatus(r *http.Request, user *models.User) core.Response {
	res := allowForModeratorOrAdmin(user)
	if res != nil {
		return res
	}

	uid, err := uuid.Parse(r.PathValue("uuid"))
	if err != nil {
		return core.Err(http.StatusBadRequest, fmt.Errorf("invalid order uuid: %w", err))
	}

	req, res := core.Validate[updateOrderStatusRequest](r, s.conform, s.validate)
	if res != nil {
		return res
	}

	status, err := enums.NewRequestStatus(req.Status)
	if err != nil {
		return core.Err(http.StatusBadRequest, fmt.Errorf("invalid order status: %w", err))
	}

	order, err := s.store.GetOrderByUUID(r.Context(), uid)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return core.Err(http.StatusNotFound, fmt.Errorf("order not found"))
		}
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to get order: %w", err))
	}

	err = s.store.UpdateOrderStatus(r.Context(), order.ID, status)
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to update order status: %w", err))
	}

	updated, err := s.store.GetOrderByUUID(r.Context(), uid)
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to load updated order: %w", err))
	}

	return core.JSON(http.StatusOK, updated)
}
