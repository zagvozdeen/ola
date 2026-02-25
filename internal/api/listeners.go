package api

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/zagvozdeen/ola/internal/store/enums"
	model "github.com/zagvozdeen/ola/internal/store/models"
)

const orderCallbackPrefix = "order_status"
const feedbackCallbackPrefix = "feedback_status"

func (s *Service) registerListeners() {
	s.eventBus.OrderCreated.Subscribe(func(ctx context.Context, order *model.Order) error {
		if s.bot == nil || order == nil {
			return nil
		}

		message, err := s.bot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    s.cfg.Telegram.GroupID,
			ParseMode: models.ParseModeMarkdown,
			Text:      buildOrderTelegramText(order),
			ReplyMarkup: models.InlineKeyboardMarkup{
				InlineKeyboard: [][]models.InlineKeyboardButton{
					getKeyboard(order.Status, orderCallbackPrefix, order.ID),
					{{Text: "Посмотреть заказ", URL: "https://t.me/ola_studio_bot?startapp"}},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("failed to send order telegram message: %w", err)
		}

		err = s.store.CreateOrderTelegramMessage(ctx, &model.OrderTelegramMessage{
			OrderID:   order.ID,
			ChatID:    message.Chat.ID,
			MessageID: message.ID,
		})
		if err != nil {
			return fmt.Errorf("failed to save order telegram message: %w", err)
		}

		return nil
	})

	s.eventBus.OrderChanged.Subscribe(func(ctx context.Context, order *model.Order) error {
		if s.bot == nil || order == nil {
			return nil
		}

		message, err := s.store.GetOrderTelegramMessageByOrderID(ctx, order.ID)
		if err != nil {
			if errors.Is(err, model.ErrNotFound) {
				return nil
			}
			return fmt.Errorf("failed to load order telegram message: %w", err)
		}

		_, err = s.bot.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:    message.ChatID,
			MessageID: message.MessageID,
			ParseMode: models.ParseModeMarkdown,
			Text:      buildOrderTelegramText(order),
			ReplyMarkup: models.InlineKeyboardMarkup{
				InlineKeyboard: [][]models.InlineKeyboardButton{
					getKeyboard(order.Status, orderCallbackPrefix, order.ID),
					{{Text: "Посмотреть заказ", URL: "https://t.me/ola_studio_bot?startapp"}},
				},
			},
		})
		if err != nil {
			//if strings.Contains(strings.ToLower(err.Error()), "message is not modified") {
			//	return nil
			//}
			return fmt.Errorf("failed to edit order telegram message: %w", err)
		}

		return nil
	})

	s.eventBus.FeedbackCreated.Subscribe(func(ctx context.Context, feedback *model.Feedback) error {
		if feedback == nil || s.bot == nil {
			return nil
		}
		message, err := s.bot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    s.cfg.Telegram.GroupID,
			ParseMode: models.ParseModeMarkdown,
			Text:      buildFeedbackTelegramText(feedback),
			ReplyMarkup: models.InlineKeyboardMarkup{
				InlineKeyboard: [][]models.InlineKeyboardButton{
					getKeyboard(feedback.Status, feedbackCallbackPrefix, feedback.ID),
					{{Text: "Посмотреть заявку", URL: "https://t.me/ola_studio_bot?startapp"}},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("failed to send telegram message: %w", err)
		}

		err = s.store.CreateFeedbackTelegramMessage(ctx, &model.FeedbackTelegramMessage{
			FeedbackID: feedback.ID,
			ChatID:     message.Chat.ID,
			MessageID:  message.ID,
		})
		if err != nil {
			return fmt.Errorf("failed to save feedback telegram message: %w", err)
		}

		return nil
	})

	s.eventBus.FeedbackChanged.Subscribe(func(ctx context.Context, feedback *model.Feedback) error {
		if s.bot == nil || feedback == nil {
			return nil
		}

		message, err := s.store.GetFeedbackTelegramMessageByFeedbackID(ctx, feedback.ID)
		if err != nil {
			if errors.Is(err, model.ErrNotFound) {
				return nil
			}
			return fmt.Errorf("failed to load order telegram message: %w", err)
		}

		_, err = s.bot.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:    message.ChatID,
			MessageID: message.MessageID,
			ParseMode: models.ParseModeMarkdown,
			Text:      buildFeedbackTelegramText(feedback),
			ReplyMarkup: models.InlineKeyboardMarkup{
				InlineKeyboard: [][]models.InlineKeyboardButton{
					getKeyboard(feedback.Status, feedbackCallbackPrefix, feedback.ID),
					{{Text: "Посмотреть заявку", URL: "https://t.me/ola_studio_bot?startapp"}},
				},
			},
		})
		if err != nil {
			//if strings.Contains(strings.ToLower(err.Error()), "message is not modified") {
			//	return nil
			//}
			return fmt.Errorf("failed to edit feedback telegram message: %w", err)
		}

		return nil
	})
}

func getKeyboard(status enums.RequestStatus, prefix string, id int) []models.InlineKeyboardButton {
	switch status {
	case enums.RequestStatusCreated:
		return []models.InlineKeyboardButton{
			{Text: "Взять в работу", CallbackData: fmt.Sprintf("%s:%d:%s", prefix, id, enums.RequestStatusInProgress)},
			{Text: "Завершить", CallbackData: fmt.Sprintf("%s:%d:%s", prefix, id, enums.RequestStatusReviewed)},
		}
	case enums.RequestStatusInProgress:
		return []models.InlineKeyboardButton{
			{Text: "Открыть", CallbackData: fmt.Sprintf("%s:%d:%s", prefix, id, enums.RequestStatusCreated)},
			{Text: "Завершить", CallbackData: fmt.Sprintf("%s:%d:%s", prefix, id, enums.RequestStatusReviewed)},
		}
	case enums.RequestStatusReviewed:
		return []models.InlineKeyboardButton{
			{Text: "Открыть", CallbackData: fmt.Sprintf("%s:%d:%s", prefix, id, enums.RequestStatusCreated)},
			{Text: "Взять в работу", CallbackData: fmt.Sprintf("%s:%d:%s", prefix, id, enums.RequestStatusInProgress)},
		}
	default:
		return nil
	}
}

func buildOrderTelegramText(order *model.Order) string {
	return fmt.Sprintf(
		"%s Заказ \\#%s\n\n*– UUID\\:* %s\n*– Статус\\:* %s\n*– Имя\\:* %s\n*– Телефон\\:* %s\n*– Комментарий\\:* %s",
		order.Status.Emoji(),
		bot.EscapeMarkdown(strconv.Itoa(order.ID)),
		bot.EscapeMarkdown(order.UUID.String()),
		bot.EscapeMarkdown(order.Status.Label()),
		bot.EscapeMarkdown(order.Name),
		bot.EscapeMarkdown(order.Phone),
		bot.EscapeMarkdown(order.Content),
	)
}

func buildFeedbackTelegramText(feedback *model.Feedback) string {
	return fmt.Sprintf(
		"%s Обратная связь \\#%s\n\n*– UUID\\:* %s\n*– Статус\\:* %s\n*– Имя\\:* %s\n*– Телефон\\:* %s\n*– Комментарий\\:* %s",
		feedback.Status.Emoji(),
		bot.EscapeMarkdown(strconv.Itoa(feedback.ID)),
		bot.EscapeMarkdown(feedback.UUID.String()),
		bot.EscapeMarkdown(feedback.Status.Label()),
		bot.EscapeMarkdown(feedback.Name),
		bot.EscapeMarkdown(feedback.Phone),
		bot.EscapeMarkdown(feedback.Content),
	)
}
