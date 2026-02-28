package store

import (
	"context"
	"fmt"
	"strings"

	"github.com/zagvozdeen/ola/internal/store/models"
)

func (s *Store) GetOrderCommentsByOrderIDs(ctx context.Context, orderIDs []int) (map[int][]models.OrderComment, error) {
	commentsByOrderID := make(map[int][]models.OrderComment, len(orderIDs))
	if len(orderIDs) == 0 {
		return commentsByOrderID, nil
	}

	placeholders := make([]string, 0, len(orderIDs))
	args := make([]any, 0, len(orderIDs))
	for _, orderID := range orderIDs {
		args = append(args, orderID)
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(args)))
	}

	rows, err := s.querier(ctx).Query(
		ctx,
		`SELECT
			oc.id,
			oc.uuid,
			oc.content,
			oc.order_id,
			oc.user_id,
			oc.created_at,
			oc.updated_at,
			u.id,
			u.uuid,
			u.first_name,
			u.last_name,
			u.username
		FROM order_comments oc
		JOIN users u ON u.id = oc.user_id
		WHERE oc.order_id IN (`+strings.Join(placeholders, ", ")+`)
		ORDER BY oc.order_id, oc.created_at, oc.id`,
		args...,
	)
	if err != nil {
		return nil, wrapDBError(err)
	}
	defer rows.Close()

	for rows.Next() {
		comment := models.OrderComment{
			Author: &models.OrderCommentAuthor{},
		}
		err = rows.Scan(
			&comment.ID,
			&comment.UUID,
			&comment.Content,
			&comment.OrderID,
			&comment.UserID,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.Author.ID,
			&comment.Author.UUID,
			&comment.Author.FirstName,
			&comment.Author.LastName,
			&comment.Author.Username,
		)
		if err != nil {
			return nil, wrapDBError(err)
		}

		commentsByOrderID[comment.OrderID] = append(commentsByOrderID[comment.OrderID], comment)
	}
	if err = rows.Err(); err != nil {
		return nil, wrapDBError(err)
	}

	return commentsByOrderID, nil
}

func (s *Store) CreateOrderComment(ctx context.Context, comment *models.OrderComment) error {
	err := s.querier(ctx).QueryRow(
		ctx,
		"INSERT INTO order_comments (uuid, content, order_id, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		comment.UUID,
		comment.Content,
		comment.OrderID,
		comment.UserID,
		comment.CreatedAt,
		comment.UpdatedAt,
	).Scan(&comment.ID)
	return wrapDBError(err)
}
