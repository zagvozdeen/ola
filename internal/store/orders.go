package store

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/zagvozdeen/ola/internal/store/enums"
	"github.com/zagvozdeen/ola/internal/store/models"
)

func (s *Store) GetAllOrders(ctx context.Context) ([]models.Order, error) {
	rows, err := s.querier(ctx).Query(ctx, "SELECT id, uuid, status, source, name, phone, content, user_id, created_at, updated_at FROM orders ORDER BY updated_at DESC, created_at DESC")
	if err != nil {
		return nil, wrapDBError(err)
	}
	defer rows.Close()

	orders := make([]models.Order, 0)
	for rows.Next() {
		order := models.Order{}
		err = rows.Scan(&order.ID, &order.UUID, &order.Status, &order.Source, &order.Name, &order.Phone, &order.Content, &order.UserID, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, wrapDBError(err)
		}
		orders = append(orders, order)
	}
	if err = rows.Err(); err != nil {
		return nil, wrapDBError(err)
	}

	return orders, nil
}

func (s *Store) GetOrderByUUID(ctx context.Context, orderUUID uuid.UUID) (*models.Order, error) {
	order := &models.Order{}
	err := s.querier(ctx).QueryRow(
		ctx,
		"SELECT id, uuid, status, source, name, phone, content, user_id, created_at, updated_at FROM orders WHERE uuid = $1",
		orderUUID,
	).Scan(
		&order.ID, &order.UUID, &order.Status, &order.Source, &order.Name, &order.Phone, &order.Content, &order.UserID, &order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		return nil, wrapDBError(err)
	}

	return order, nil
}

func (s *Store) GetOrderByID(ctx context.Context, orderID int) (*models.Order, error) {
	order := &models.Order{}
	err := s.querier(ctx).QueryRow(
		ctx,
		"SELECT id, uuid, status, source, name, phone, content, user_id, created_at, updated_at FROM orders WHERE id = $1",
		orderID,
	).Scan(
		&order.ID, &order.UUID, &order.Status, &order.Source, &order.Name, &order.Phone, &order.Content, &order.UserID, &order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		return nil, wrapDBError(err)
	}

	return order, nil
}

func (s *Store) CreateOrder(ctx context.Context, order *models.Order) error {
	err := s.querier(ctx).QueryRow(
		ctx,
		"INSERT INTO orders (uuid, status, source, name, phone, content, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id",
		order.UUID, order.Status, order.Source, order.Name, order.Phone, order.Content, order.UserID, order.CreatedAt, order.UpdatedAt,
	).Scan(&order.ID)
	return wrapDBError(err)
}

func (s *Store) UpdateOrderStatus(ctx context.Context, order *models.Order) error {
	_, err := s.querier(ctx).Exec(
		ctx,
		"UPDATE orders SET status = $1, updated_at = $2 WHERE id = $3",
		order.Status, order.UpdatedAt, order.ID,
	)
	return wrapDBError(err)
}

func (s *Store) GetOrderItemsByOrderIDs(ctx context.Context, orderIDs []int) (map[int][]models.OrderItem, error) {
	itemsByOrderID := make(map[int][]models.OrderItem, len(orderIDs))
	if len(orderIDs) == 0 {
		return itemsByOrderID, nil
	}

	placeholders := make([]string, 0, len(orderIDs))
	args := make([]any, 0, len(orderIDs))
	for _, orderID := range orderIDs {
		args = append(args, orderID)
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(args)))
	}

	rows, err := s.querier(ctx).Query(
		ctx,
		`SELECT oi.order_id, oi.product_id, oi.product_name, oi.price_from, oi.price_to, oi.qty, p.file_content
		FROM order_items oi
		LEFT JOIN products p ON p.id = oi.product_id
		WHERE oi.order_id IN (`+strings.Join(placeholders, ", ")+`)
		ORDER BY oi.order_id, oi.product_name, oi.product_id`,
		args...,
	)
	if err != nil {
		return nil, wrapDBError(err)
	}
	defer rows.Close()

	for rows.Next() {
		item := models.OrderItem{}
		err = rows.Scan(
			&item.OrderID,
			&item.ProductID,
			&item.ProductName,
			&item.PriceFrom,
			&item.PriceTo,
			&item.Qty,
			&item.FileContent,
		)
		if err != nil {
			return nil, wrapDBError(err)
		}

		itemsByOrderID[item.OrderID] = append(itemsByOrderID[item.OrderID], item)
	}
	if err = rows.Err(); err != nil {
		return nil, wrapDBError(err)
	}

	return itemsByOrderID, nil
}

func (s *Store) CreateOrderFromUserCart(ctx context.Context, userID int, source enums.OrderSource, name, phone, content string) (*models.Order, error) {
	cart, err := s.getOrCreateUserCartByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	uid, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	order := &models.Order{
		UUID:      uid,
		Status:    enums.RequestStatusCreated,
		Source:    source,
		Name:      name,
		Phone:     phone,
		Content:   content,
		UserID:    &userID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = s.querier(ctx).QueryRow(
		ctx,
		"INSERT INTO orders (uuid, status, source, name, phone, content, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id",
		order.UUID, order.Status, order.Source, order.Name, order.Phone, order.Content, order.UserID, order.CreatedAt, order.UpdatedAt,
	).Scan(&order.ID)
	if err != nil {
		return nil, wrapDBError(err)
	}

	tag, err := s.querier(ctx).Exec(
		ctx,
		"INSERT INTO order_items (order_id, product_id, product_name, price_from, price_to, qty) SELECT $1, p.id, p.name, p.price_from, p.price_to, ci.qty FROM cart_items ci JOIN products p ON p.id = ci.product_id WHERE ci.cart_id = $2",
		order.ID,
		cart.ID,
	)
	if err != nil {
		return nil, wrapDBError(err)
	}
	if tag.RowsAffected() == 0 {
		return nil, models.ErrCartEmpty
	}

	_, err = s.querier(ctx).Exec(ctx, "DELETE FROM cart_items WHERE cart_id = $1", cart.ID)
	if err != nil {
		return nil, wrapDBError(err)
	}

	return order, nil
}
