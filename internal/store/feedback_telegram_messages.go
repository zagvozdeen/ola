package store

import (
	"context"

	"github.com/zagvozdeen/ola/internal/store/models"
)

func (s *Store) GetFeedbackTelegramMessageByFeedbackID(ctx context.Context, id int) (*models.FeedbackTelegramMessage, error) {
	message := &models.FeedbackTelegramMessage{}
	err := s.querier(ctx).QueryRow(
		ctx,
		"SELECT feedback_id, chat_id, message_id FROM feedback_telegram_messages WHERE feedback_id = $1 ORDER BY message_id DESC LIMIT 1",
		id,
	).Scan(&message.FeedbackID, &message.ChatID, &message.MessageID)
	if err != nil {
		return nil, wrapDBError(err)
	}
	return message, nil
}

func (s *Store) CreateFeedbackTelegramMessage(ctx context.Context, message *models.FeedbackTelegramMessage) error {
	_, err := s.querier(ctx).Exec(
		ctx,
		"INSERT INTO feedback_telegram_messages (feedback_id, chat_id, message_id) VALUES ($1, $2, $3)",
		message.FeedbackID,
		message.ChatID,
		message.MessageID,
	)
	return wrapDBError(err)
}
