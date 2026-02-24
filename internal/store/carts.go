package store

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/zagvozdeen/ola/internal/store/models"
)

func (s *Store) getUserCartByUserID(ctx context.Context, userID int) (*models.Cart, error) {
	cart := &models.Cart{}

	err := s.querier(ctx).QueryRow(
		ctx,
		"SELECT id, uuid, user_id, session_id::text, created_at, updated_at FROM carts WHERE user_id = $1",
		userID,
	).Scan(&cart.ID, &cart.UUID, &cart.UserID, &cart.SessionID, &cart.CreatedAt, &cart.UpdatedAt)
	if err != nil {
		return nil, wrapDBError(err)
	}

	return cart, nil
}

func (s *Store) getOrCreateUserCartByUserID(ctx context.Context, userID int) (*models.Cart, error) {
	cart, err := s.getUserCartByUserID(ctx, userID)
	if err == nil {
		return cart, nil
	}

	if !errors.Is(err, models.ErrNotFound) {
		return nil, err
	}

	uid, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	cart = &models.Cart{
		UUID:      uid,
		UserID:    &userID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = s.querier(ctx).QueryRow(
		ctx,
		"INSERT INTO carts (uuid, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id",
		cart.UUID,
		cart.UserID,
		cart.CreatedAt,
		cart.UpdatedAt,
	).Scan(&cart.ID)
	if err != nil {
		err = wrapDBError(err)
		if errors.Is(err, models.ErrUniqueViolation) {
			return s.getUserCartByUserID(ctx, userID)
		}

		return nil, err
	}

	return cart, nil
}

func (s *Store) GetUserCartItems(ctx context.Context, userID int) ([]models.CartItem, error) {
	cart, err := s.getUserCartByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return []models.CartItem{}, nil
		}

		return nil, err
	}

	rows, err := s.querier(ctx).Query(
		ctx,
		"SELECT ci.product_id, p.uuid, p.name, p.price_from, p.price_to, p.type, f.content, ci.qty FROM cart_items ci JOIN products p ON p.id = ci.product_id JOIN files f ON f.id = p.file_id WHERE ci.cart_id = $1 ORDER BY p.created_at DESC",
		cart.ID,
	)
	if err != nil {
		return nil, wrapDBError(err)
	}
	defer rows.Close()

	items := make([]models.CartItem, 0)
	for rows.Next() {
		item := models.CartItem{}
		err = rows.Scan(
			&item.ProductID,
			&item.ProductUUID,
			&item.ProductName,
			&item.PriceFrom,
			&item.PriceTo,
			&item.Type,
			&item.FileContent,
			&item.Qty,
		)
		if err != nil {
			return nil, wrapDBError(err)
		}

		items = append(items, item)
	}
	if err = rows.Err(); err != nil {
		return nil, wrapDBError(err)
	}

	return items, nil
}

func (s *Store) UpsertUserCartItem(ctx context.Context, userID, productID, qty int) error {
	cart, err := s.getOrCreateUserCartByUserID(ctx, userID)
	if err != nil {
		return err
	}

	_, err = s.querier(ctx).Exec(
		ctx,
		"INSERT INTO cart_items (cart_id, product_id, qty) VALUES ($1, $2, $3) ON CONFLICT (cart_id, product_id) DO UPDATE SET qty = EXCLUDED.qty",
		cart.ID,
		productID,
		qty,
	)
	return wrapDBError(err)
}

func (s *Store) DeleteUserCartItem(ctx context.Context, userID int, productUUID uuid.UUID) error {
	cart, err := s.getUserCartByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil
		}

		return err
	}

	_, err = s.querier(ctx).Exec(
		ctx,
		"DELETE FROM cart_items ci USING products p WHERE ci.cart_id = $1 AND ci.product_id = p.id AND p.uuid = $2",
		cart.ID,
		productUUID,
	)
	return wrapDBError(err)
}
