package api

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/google/uuid"
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

		var user *model.User
		var err error
		if order.UserID != nil {
			user, err = s.store.GetUserByID(ctx, *order.UserID)
			if err != nil {
				return fmt.Errorf("failed to get user: %w", err)
			}
		}

		message, err := s.bot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      s.cfg.Telegram.GroupID,
			ParseMode:   models.ParseModeMarkdown,
			Text:        buildOrderTelegramText(order, user),
			ReplyMarkup: getKeyboard(order.Status, orderCallbackPrefix, order.ID, order.UUID),
		})
		if err != nil {
			return fmt.Errorf("failed to send order telegram message: %w", err)
		}

		err = s.store.CreateOrderTelegramMessage(ctx, &model.OrderTelegramMessage{
			OrderID:   order.ID,
			ChatID:    message.Chat.ID,
			MessageID: int64(message.ID),
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

		var user *model.User
		var err error
		if order.UserID != nil {
			user, err = s.store.GetUserByID(ctx, *order.UserID)
			if err != nil {
				return fmt.Errorf("failed to get user: %w", err)
			}
		}

		message, err := s.store.GetOrderTelegramMessageByOrderID(ctx, order.ID)
		if err != nil {
			if errors.Is(err, model.ErrNotFound) {
				return nil
			}
			return fmt.Errorf("failed to load order telegram message: %w", err)
		}

		_, err = s.bot.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      message.ChatID,
			MessageID:   int(message.MessageID),
			ParseMode:   models.ParseModeMarkdown,
			Text:        buildOrderTelegramText(order, user),
			ReplyMarkup: getKeyboard(order.Status, orderCallbackPrefix, order.ID, order.UUID),
		})
		if err != nil {
			return fmt.Errorf("failed to edit order telegram message: %w", err)
		}

		return nil
	})

	s.eventBus.FeedbackCreated.Subscribe(func(ctx context.Context, feedback *model.Feedback) error {
		if feedback == nil || s.bot == nil {
			return nil
		}

		var user *model.User
		var err error
		user, err = s.store.GetUserByID(ctx, feedback.UserID)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}

		message, err := s.bot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      s.cfg.Telegram.GroupID,
			ParseMode:   models.ParseModeMarkdown,
			Text:        buildFeedbackTelegramText(feedback, user),
			ReplyMarkup: getKeyboard(feedback.Status, feedbackCallbackPrefix, feedback.ID, feedback.UUID),
		})
		if err != nil {
			return fmt.Errorf("failed to send telegram message: %w", err)
		}

		err = s.store.CreateFeedbackTelegramMessage(ctx, &model.FeedbackTelegramMessage{
			FeedbackID: feedback.ID,
			ChatID:     message.Chat.ID,
			MessageID:  int64(message.ID),
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

		var user *model.User
		var err error
		user, err = s.store.GetUserByID(ctx, feedback.UserID)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}

		message, err := s.store.GetFeedbackTelegramMessageByFeedbackID(ctx, feedback.ID)
		if err != nil {
			if errors.Is(err, model.ErrNotFound) {
				return nil
			}
			return fmt.Errorf("failed to load order telegram message: %w", err)
		}

		_, err = s.bot.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      message.ChatID,
			MessageID:   int(message.MessageID),
			ParseMode:   models.ParseModeMarkdown,
			Text:        buildFeedbackTelegramText(feedback, user),
			ReplyMarkup: getKeyboard(feedback.Status, feedbackCallbackPrefix, feedback.ID, feedback.UUID),
		})
		if err != nil {
			return fmt.Errorf("failed to edit feedback telegram message: %w", err)
		}

		return nil
	})
}

func getKeyboard(status enums.RequestStatus, prefix string, id int, uuid uuid.UUID) models.ReplyMarkup {
	var text, value string
	switch prefix {
	case orderCallbackPrefix:
		text = "Посмотреть заказ"
		value = base64.URLEncoding.EncodeToString([]byte("order:" + uuid.String()))
	case feedbackCallbackPrefix:
		text = "Посмотреть заявку"
		value = base64.URLEncoding.EncodeToString([]byte("feedback:" + uuid.String()))
	default:
		return nil
	}
	var keyboard []models.InlineKeyboardButton
	switch status {
	case enums.RequestStatusCreated:
		keyboard = []models.InlineKeyboardButton{
			{Text: "Взять в работу", CallbackData: fmt.Sprintf("%s:%d:%s", prefix, id, enums.RequestStatusInProgress)},
			{Text: "Завершить", CallbackData: fmt.Sprintf("%s:%d:%s", prefix, id, enums.RequestStatusReviewed)},
		}
	case enums.RequestStatusInProgress:
		keyboard = []models.InlineKeyboardButton{
			{Text: "Открыть", CallbackData: fmt.Sprintf("%s:%d:%s", prefix, id, enums.RequestStatusCreated)},
			{Text: "Завершить", CallbackData: fmt.Sprintf("%s:%d:%s", prefix, id, enums.RequestStatusReviewed)},
		}
	case enums.RequestStatusReviewed:
		keyboard = []models.InlineKeyboardButton{
			{Text: "Открыть", CallbackData: fmt.Sprintf("%s:%d:%s", prefix, id, enums.RequestStatusCreated)},
			{Text: "Взять в работу", CallbackData: fmt.Sprintf("%s:%d:%s", prefix, id, enums.RequestStatusInProgress)},
		}
	default:
		return nil
	}
	return models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			keyboard,
			{{Text: text, URL: "https://t.me/ola_studio_bot?startapp=" + value}},
		},
	}
}

func buildOrderTelegramText(order *model.Order, user *model.User) string {
	name := bot.EscapeMarkdown(order.Name)
	if user != nil && user.Username != nil {
		name = fmt.Sprintf("[%s](%s)", name, bot.EscapeMarkdown("https://t.me/"+*user.Username))
	}
	return fmt.Sprintf(
		"%s Заказ \\#%s\n\n*– UUID\\:* %s\n*– Статус\\:* %s\n*– Имя\\:* %s\n*– Телефон\\:* %s\n*– Комментарий\\:* %s",
		order.Status.Emoji(),
		bot.EscapeMarkdown(strconv.Itoa(order.ID)),
		bot.EscapeMarkdown(order.UUID.String()),
		bot.EscapeMarkdown(order.Status.Label()),
		name,
		bot.EscapeMarkdown(order.Phone),
		bot.EscapeMarkdown(order.Content),
	)
}

func buildFeedbackTelegramText(feedback *model.Feedback, user *model.User) string {
	name := bot.EscapeMarkdown(feedback.Name)
	if user != nil && user.Username != nil {
		name = fmt.Sprintf("[%s](%s)", name, bot.EscapeMarkdown("https://t.me/"+*user.Username))
	}
	return fmt.Sprintf(
		"%s Обратная связь \\#%s\n\n*– UUID\\:* %s\n*– Статус\\:* %s\n*– Тип\\:* %s\n*– Имя\\:* %s\n*– Телефон\\:* %s\n*– Комментарий\\:* %s",
		feedback.Status.Emoji(),
		bot.EscapeMarkdown(strconv.Itoa(feedback.ID)),
		bot.EscapeMarkdown(feedback.UUID.String()),
		bot.EscapeMarkdown(feedback.Status.Label()),
		bot.EscapeMarkdown(feedback.Type.Label()),
		name,
		bot.EscapeMarkdown(feedback.Phone),
		bot.EscapeMarkdown(feedback.Content),
	)
}
