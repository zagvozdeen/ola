package api

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/google/uuid"
	"github.com/zagvozdeen/ola/internal/store/enums"
	model "github.com/zagvozdeen/ola/internal/store/models"
)

var errTelegramBotDisabled = errors.New("telegram bot disabled")

func (s *Service) startBot(ctx context.Context) error {
	if !s.cfg.Telegram.BotEnabled {
		s.log.Info("Telegram bot disabled")
		return errTelegramBotDisabled
	}
	b, err := bot.New(
		s.cfg.Telegram.BotToken,
		bot.WithDefaultHandler(s.defaultHandler),
		bot.WithCallbackQueryDataHandler(orderCallbackPrefix, bot.MatchTypePrefix, s.callbackQueryHandler(s.handleOrderStatusCallback, enums.UserRoleModerator, enums.UserRoleAdmin)),
		bot.WithCallbackQueryDataHandler(feedbackCallbackPrefix, bot.MatchTypePrefix, s.callbackQueryHandler(s.handleFeedbackStatusCallback, enums.UserRoleModerator, enums.UserRoleAdmin)),
	)
	if err != nil {
		return err
	}
	s.bot = b
	b.Start(ctx)
	return nil
}

func (s *Service) createUserIfNotExists(ctx context.Context, from models.User) (*model.User, error) {
	user, err := s.store.GetUserByTID(ctx, from.ID)
	if err != nil {
		if !errors.Is(err, model.ErrNotFound) {
			return nil, fmt.Errorf("failed to get user by TID: %w", err)
		}
		var uid uuid.UUID
		uid, err = uuid.NewV7()
		if err != nil {
			return nil, fmt.Errorf("failed to generate uuid: %w", err)
		}
		user = &model.User{
			TID:       new(from.ID),
			UUID:      uid,
			FirstName: from.FirstName,
			LastName:  new(from.LastName),
			Username:  new(from.Username),
			Role:      enums.UserRoleUser,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err = s.store.CreateUser(ctx, user)
		if err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	}
	return user, nil
}

func (s *Service) defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update == nil || update.Message == nil {
		return
	}

	_, err := s.createUserIfNotExists(ctx, *update.Message.From)
	if err != nil {
		s.log.Error("Failed to create user", err)
		return
	}

	if update.Message.Chat.Type == models.ChatTypePrivate {
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			ParseMode: models.ParseModeMarkdown,
			Text:      "*Давайте сделаем заказ 🎈*\n\nНажмите на кнопку ниже, чтобы сделать ваш праздник\\!",
			ReplyMarkup: models.InlineKeyboardMarkup{InlineKeyboard: [][]models.InlineKeyboardButton{{{
				Text:   "Заказать продукт",
				WebApp: &models.WebAppInfo{URL: s.cfg.Telegram.MiniAppURL},
			}, {
				Text:   "Стать партнёром",
				WebApp: &models.WebAppInfo{URL: "https://ola.creavo.ru/spa/settings/partnership"},
			}}}},
		})
		if err != nil {
			s.log.Error("Failed to send telegram message", err)
			return
		}
	}
}

func (s *Service) callbackQueryHandler(fn func(context.Context, *bot.Bot, *models.CallbackQuery, *model.User) (string, error), roles ...enums.UserRole) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		callback := update.CallbackQuery
		if callback == nil {
			return
		}

		user, err := s.createUserIfNotExists(ctx, callback.From)
		if err != nil {
			s.log.Error("Failed to create user", err)
			return
		}
		if !slices.Contains(roles, user.Role) {
			s.answerOrderStatusCallback(ctx, b, callback.ID, "Это действие недоступно для вас")
			return
		}

		text, err := fn(ctx, b, callback, user)
		if err != nil {
			s.log.Error("Failed to callback query", err)
		}
		s.answerOrderStatusCallback(ctx, b, callback.ID, text)
	}
}

func (s *Service) handleOrderStatusCallback(ctx context.Context, b *bot.Bot, callback *models.CallbackQuery, user *model.User) (string, error) {
	orderID, status, ok := parseStatusCallbackData(callback.Data, orderCallbackPrefix)
	if !ok {
		return "Не удалось распарсить данные", nil
	}

	order, err := s.store.GetOrderByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return "Заказ не найден", nil
		}
		return "Не удалось получить заказ", fmt.Errorf("failed to load order after telegram callback: %w", err)
	}

	order.Status = status
	order.UpdatedAt = time.Now()
	err = s.store.UpdateOrderStatus(ctx, order)
	if err != nil {
		return "Не удалось обновить статус", fmt.Errorf("failed to update order status from telegram callback: %w", err)
	}

	s.eventBus.OrderChanged.Publish(context.WithoutCancel(ctx), order)
	return fmt.Sprintf("Статус: %s", status.Label()), nil
}

func (s *Service) handleFeedbackStatusCallback(ctx context.Context, b *bot.Bot, callback *models.CallbackQuery, user *model.User) (string, error) {
	feedbackID, status, ok := parseStatusCallbackData(callback.Data, feedbackCallbackPrefix)
	if !ok {
		return "Не удалось распарсить данные", nil
	}

	feedback, err := s.store.GetFeedbackByID(ctx, feedbackID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return "Обратная связь не найден", nil
		}
		return "Не удалось получить обратную связь", fmt.Errorf("failed to load order after telegram callback: %w", err)
	}

	feedback.Status = status
	feedback.UpdatedAt = time.Now()
	err = s.store.UpdateFeedbackStatus(ctx, feedback)
	if err != nil {
		return "Не удалось обновить обратную связь", fmt.Errorf("failed to update order status from telegram callback: %w", err)
	}

	s.eventBus.FeedbackChanged.Publish(context.WithoutCancel(ctx), feedback)
	return fmt.Sprintf("Статус: %s", status.Label()), nil
}

func (s *Service) answerOrderStatusCallback(ctx context.Context, b *bot.Bot, callbackID string, text string) {
	_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: callbackID,
		Text:            text,
		ShowAlert:       false,
	})
	if err != nil {
		s.log.Error("Failed to answer callback query", err)
	}
}

func parseStatusCallbackData(data string, prefix string) (int, enums.RequestStatus, bool) {
	parts := strings.Split(data, ":")
	if len(parts) != 3 || parts[0] != prefix {
		return 0, enums.RequestStatus{}, false
	}

	orderID, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, enums.RequestStatus{}, false
	}

	status, err := enums.NewRequestStatus(parts[2])
	if err != nil {
		return 0, enums.RequestStatus{}, false
	}

	return orderID, status, true
}
