package api

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	model "github.com/zagvozdeen/ola/internal/store/models"
)

func (s *Service) registerListeners() {
	s.eventBus.OrderCreated.Subscribe(func(ctx context.Context, order *model.Order) error {
		if s.bot == nil {
			return nil
		}
		_, err := s.bot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    s.cfg.Telegram.GroupID,
			ParseMode: models.ParseModeMarkdown,
			Text:      "*Пришёл новый заказ 🎈*\n\nНажмите на кнопку ниже, чтобы посмотреть заказ\\!",
			ReplyMarkup: models.InlineKeyboardMarkup{InlineKeyboard: [][]models.InlineKeyboardButton{{{
				Text: "Посмотреть заказ",
				URL:  "https://t.me/ola_studio_bot?startapp",
			}}}},
		})
		if err != nil {
			return fmt.Errorf("failed to send telegram message: %w", err)
		}
		return nil
	})

	s.eventBus.FeedbackCreated.Subscribe(func(ctx context.Context, feedback *model.Feedback) error {
		if s.bot == nil {
			return nil
		}
		_, err := s.bot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    s.cfg.Telegram.GroupID,
			ParseMode: models.ParseModeMarkdown,
			Text:      "*Пришла новая заявка на обслуживание 🎈*\n\nНажмите на кнопку ниже, чтобы посмотреть заявку\\!",
			ReplyMarkup: models.InlineKeyboardMarkup{InlineKeyboard: [][]models.InlineKeyboardButton{{{
				Text: "Посмотреть заявку",
				URL:  "https://t.me/ola_studio_bot?startapp",
			}}}},
		})
		if err != nil {
			return fmt.Errorf("failed to send telegram message: %w", err)
		}
		return nil
	})
}
