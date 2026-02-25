package api

import (
	"context"
	"errors"
	"fmt"
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
		bot.WithCallbackQueryDataHandler(orderCallbackPrefix, bot.MatchTypePrefix, s.handleOrderStatusCallback),
		bot.WithCallbackQueryDataHandler(feedbackCallbackPrefix, bot.MatchTypePrefix, s.handleOrderStatusCallback),
	)
	if err != nil {
		return err
	}
	s.bot = b
	b.Start(ctx)
	return nil
}

func (s *Service) defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update == nil || update.Message == nil {
		return
	}

	_, err := s.store.GetUserByTID(ctx, update.Message.From.ID)
	if err != nil {
		if !errors.Is(err, model.ErrNotFound) {
			s.log.Error("Failed to get user by TID", err)
			return
		}
		var uid uuid.UUID
		uid, err = uuid.NewV7()
		if err != nil {
			s.log.Error("Failed generate uuid", err)
			return
		}
		err = s.store.CreateUser(ctx, &model.User{
			TID:       new(update.Message.From.ID),
			UUID:      uid,
			FirstName: update.Message.From.FirstName,
			LastName:  new(update.Message.From.LastName),
			Username:  new(update.Message.From.Username),
			Role:      enums.UserRoleUser,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
		if err != nil {
			s.log.Error("Failed to create user", err)
			return
		}
	}

	keyboard := [][]models.InlineKeyboardButton{{{
		Text:   "Заказать продукт",
		WebApp: &models.WebAppInfo{URL: s.cfg.Telegram.MiniAppURL},
	}}}
	if update.Message.Chat.Type != models.ChatTypePrivate {
		keyboard = [][]models.InlineKeyboardButton{{{
			Text: "Заказать продукт",
			URL:  "https://t.me/ola_studio_bot?startapp",
		}}}
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		ParseMode:   models.ParseModeMarkdown,
		Text:        "*Давайте сделаем заказ 🎈*\n\nНажмите на кнопку ниже, чтобы сделать ваш праздник\\!",
		ReplyMarkup: models.InlineKeyboardMarkup{InlineKeyboard: keyboard},
	})
	if err != nil {
		s.log.Error("Failed to send telegram message", err)
		return
	}
}

func (s *Service) handleOrderStatusCallback(ctx context.Context, b *bot.Bot, update *models.Update) {
	callback := update.CallbackQuery
	if callback == nil {
		return
	}

	orderID, status, ok := parseOrderStatusCallbackData(callback.Data)
	if !ok {
		return
	}

	err := s.store.UpdateOrderStatus(ctx, orderID, status)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			s.answerOrderStatusCallback(ctx, b, callback.ID, "Заказ не найден")
			return
		}
		s.log.Error("Failed to update order status from telegram callback", err)
		s.answerOrderStatusCallback(ctx, b, callback.ID, "Не удалось обновить статус")
		return
	}

	order, err := s.store.GetOrderByID(ctx, orderID)
	if err != nil {
		s.log.Error("Failed to load order after telegram callback", err)
		s.answerOrderStatusCallback(ctx, b, callback.ID, "Статус обновлён, но не удалось обновить сообщение")
		return
	}

	s.eventBus.OrderChanged.Publish(context.WithoutCancel(ctx), order)
	s.answerOrderStatusCallback(ctx, b, callback.ID, fmt.Sprintf("Статус: %s", status.Label()))
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

func parseOrderStatusCallbackData(data string) (int, enums.RequestStatus, bool) {
	parts := strings.Split(data, ":")
	if len(parts) != 3 || parts[0] != orderCallbackPrefix {
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
