package api

import (
	"context"
	"errors"
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
	b, err := bot.New(s.cfg.Telegram.BotToken, bot.WithDefaultHandler(s.defaultHandler), bot.WithDebug())
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
