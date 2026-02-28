package api

import (
	"context"
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

func sourceFromAuthHeader(authorization string) enums.OrderSource {
	if strings.HasPrefix(authorization, "tma ") {
		return enums.OrderSourceTMA
	}

	return enums.OrderSourceSPA
}

func (s *Service) getOrders(r *http.Request, user *models.User) core.Response {
	res := allowForOrderManager(user)
	if res != nil {
		return res
	}

	orders, err := s.store.GetAllOrders(r.Context())
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to get orders: %w", err))
	}

	err = s.attachOrderDetails(r.Context(), orders)
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to load order details: %w", err))
	}

	return core.JSON(http.StatusOK, orders)
}

func (s *Service) getOrder(r *http.Request, user *models.User) core.Response {
	res := allowForOrderManager(user)
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

	orderList := []models.Order{*order}
	err = s.attachOrderDetails(r.Context(), orderList)
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to load order details: %w", err))
	}
	*order = orderList[0]

	return core.JSON(http.StatusOK, order)
}

func (s *Service) createOrder(r *http.Request, user *models.User) core.Response {
	req, res := core.Validate[createOrderRequest](r, s.conform, s.validate)
	if res != nil {
		return res
	}

	user.Phone = new(req.Phone)
	user.UpdatedAt = time.Now()
	err := s.store.UpdateUserPhone(r.Context(), user)
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to update user phone: %w", err))
	}

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

	s.eventBus.OrderCreated.Publish(context.WithoutCancel(r.Context()), order)

	return core.JSON(http.StatusCreated, order)
}

func (s *Service) createOrderFromCart(r *http.Request, user *models.User) core.Response {
	req, res := core.Validate[createOrderRequest](r, s.conform, s.validate)
	if res != nil {
		return res
	}

	ctx, err := s.store.Begin(r.Context())
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to begin transaction: %w", err))
	}
	defer s.store.Rollback(ctx)

	user.Phone = new(req.Phone)
	user.UpdatedAt = time.Now()
	err = s.store.UpdateUserPhone(ctx, user)
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to update user phone: %w", err))
	}

	order, err := s.store.CreateOrderFromUserCart(
		ctx,
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

	s.store.Commit(ctx)

	s.eventBus.OrderCreated.Publish(context.WithoutCancel(r.Context()), order)

	return core.JSON(http.StatusCreated, order)
}

type createGuestOrderRequest struct {
	Name    string `json:"name" mold:"trim" validate:"required,max=255"`
	Phone   string `json:"phone" mold:"trim" validate:"required,max=255,ru_phone"`
	Content string `json:"content" mold:"trim" validate:"required,max=3000"`
	Consent bool   `json:"consent" validate:"required"`
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
		UUID:      uid,
		Status:    enums.RequestStatusCreated,
		Source:    enums.OrderSourceLanding,
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

	s.eventBus.OrderCreated.Publish(context.WithoutCancel(r.Context()), order)

	return core.JSON(http.StatusCreated, order)
}

type updateOrderStatusRequest struct {
	Status  *enums.RequestStatus `json:"status"`
	Comment string               `json:"comment" mold:"trim" validate:"omitempty,max=3000"`
}

func (s *Service) updateOrderStatus(r *http.Request, user *models.User) core.Response {
	res := allowForOrderManager(user)
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

	if req.Status == nil && req.Comment == "" {
		return core.Err(http.StatusBadRequest, fmt.Errorf("nothing to update"))
	}

	ctx, err := s.store.Begin(r.Context())
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to begin transaction: %w", err))
	}
	defer s.store.Rollback(ctx)

	order, err := s.store.GetOrderByUUID(ctx, uid)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return core.Err(http.StatusNotFound, fmt.Errorf("order not found"))
		}
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to get order: %w", err))
	}

	now := time.Now()
	if req.Status != nil {
		order.Status = *req.Status
	}

	order.UpdatedAt = now
	err = s.store.UpdateOrderStatus(ctx, order)
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to update order status: %w", err))
	}

	if req.Comment != "" {
		commentUUID, commentErr := uuid.NewV7()
		if commentErr != nil {
			return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to generate order comment uuid: %w", commentErr))
		}

		comment := &models.OrderComment{
			UUID:      commentUUID,
			Content:   req.Comment,
			OrderID:   order.ID,
			UserID:    user.ID,
			CreatedAt: now,
			UpdatedAt: now,
		}

		err = s.store.CreateOrderComment(ctx, comment)
		if err != nil {
			return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to create order comment: %w", err))
		}
	}

	orderList := []models.Order{*order}
	err = s.attachOrderDetails(ctx, orderList)
	if err != nil {
		return core.Err(http.StatusInternalServerError, fmt.Errorf("failed to load order details: %w", err))
	}
	*order = orderList[0]

	s.store.Commit(ctx)

	s.eventBus.OrderChanged.Publish(context.WithoutCancel(r.Context()), order)

	return core.JSON(http.StatusOK, order)
}

func (s *Service) attachOrderDetails(ctx context.Context, orders []models.Order) error {
	if len(orders) == 0 {
		return nil
	}

	orderIDs := make([]int, 0, len(orders))
	for i := range orders {
		orders[i].Items = make([]models.OrderItem, 0)
		orders[i].Comments = make([]models.OrderComment, 0)
		orderIDs = append(orderIDs, orders[i].ID)
	}

	itemsByOrderID, err := s.store.GetOrderItemsByOrderIDs(ctx, orderIDs)
	if err != nil {
		return err
	}

	commentsByOrderID, err := s.store.GetOrderCommentsByOrderIDs(ctx, orderIDs)
	if err != nil {
		return err
	}

	for i := range orders {
		if items, ok := itemsByOrderID[orders[i].ID]; ok {
			orders[i].Items = items
		}
		if comments, ok := commentsByOrderID[orders[i].ID]; ok {
			orders[i].Comments = comments
		}
	}

	return nil
}
