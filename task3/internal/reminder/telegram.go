package reminder

import (
	"context"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pozedorum/WB_project_4/task3/internal/interfaces"
)

type telegramClient struct {
	bot *tgbotapi.BotAPI
}

func NewTelegramClient(token string) (interfaces.TelegramClient, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create Telegram bot: %w", err)
	}

	return &telegramClient{
		bot: bot,
	}, nil
}

func (c *telegramClient) SendMessage(ctx context.Context, chatID int64, text string) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("send message cancelled: %w", ctx.Err())
	default:
	}

	msg := tgbotapi.NewMessage(chatID, text)

	done := make(chan error, 1)
	go func() {
		_, err := c.bot.Send(msg)
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("failed to send Telegram message: %w", err)
		}
		return nil
	case <-ctx.Done():
		return fmt.Errorf("send message timeout: %w", ctx.Err())
	case <-time.After(10 * time.Second):
		return fmt.Errorf("send message timeout after 10 seconds")
	}
}

func (c *telegramClient) GetBotInfo() string {
	return c.bot.Self.UserName
}

func (c *telegramClient) Close() {
	// Telegram client doesn't require explicit closing
}
