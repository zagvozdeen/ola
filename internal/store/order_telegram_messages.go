package store

import (
	"context"

	"github.com/zagvozdeen/ola/internal/store/models"
)

func (s *Store) GetOrderTelegramMessageByOrderID(ctx context.Context, orderID int) (*models.OrderTelegramMessage, error) {
	message := &models.OrderTelegramMessage{}
	err := s.querier(ctx).QueryRow(
		ctx,
		"SELECT order_id, chat_id, message_id FROM order_telegram_messages WHERE order_id = $1 ORDER BY message_id DESC LIMIT 1",
		orderID,
	).Scan(&message.OrderID, &message.ChatID, &message.MessageID)
	if err != nil {
		return nil, wrapDBError(err)
	}
	return message, nil
}

func (s *Store) CreateOrderTelegramMessage(ctx context.Context, message *models.OrderTelegramMessage) error {
	//_, err := s.querier(ctx).Exec(
	//	ctx,
	//	"DELETE FROM order_telegram_messages WHERE order_id = $1",
	//	message.OrderID,
	//)
	//if err != nil {
	//	return wrapDBError(err)
	//}

	_, err := s.querier(ctx).Exec(
		ctx,
		"INSERT INTO order_telegram_messages (order_id, chat_id, message_id) VALUES ($1, $2, $3)",
		message.OrderID,
		message.ChatID,
		message.MessageID,
	)
	return wrapDBError(err)
}
