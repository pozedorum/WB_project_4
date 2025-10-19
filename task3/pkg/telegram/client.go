// Package telegram предоставляет клиент для работы с Telegram Bot API.
package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Client представляет клиент для работы с Telegram Bot API.
type Client struct {
	bot     *tgbotapi.BotAPI
	timeout int
}

// Config представляет конфигурацию для Telegram клиента.
type Config struct {
	Token   string
	Timeout int // timeout в секундах
}

// NewClient создает новый Telegram клиент.
func NewClient(cfg Config) (*Client, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to create Telegram bot: %w", err)
	}

	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 30 // дефолтный timeout
	}

	return &Client{
		bot:     bot,
		timeout: timeout,
	}, nil
}

// SendMessage отправляет сообщение пользователю.
func (c *Client) SendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := c.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send Telegram message: %w", err)
	}
	return nil
}

// SendMessageWithKeyboard отправляет сообщение с клавиатурой.
func (c *Client) SendMessageWithKeyboard(chatID int64, text string, keyboard tgbotapi.ReplyKeyboardMarkup) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	_, err := c.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send Telegram message with keyboard: %w", err)
	}
	return nil
}

// GetBotInfo возвращает информацию о боте.
func (c *Client) GetBotInfo() tgbotapi.User {
	return c.bot.Self
}

// SetWebhook устанавливает webhook для бота.
func (c *Client) SetWebhook(url string) error {
	wh := tgbotapi.NewWebhook(url)
	_, err := c.bot.SetWebhook(wh)
	if err != nil {
		return fmt.Errorf("failed to set webhook: %w", err)
	}
	return nil
}

// Close закрывает клиент (если требуется).
func (c *Client) Close() {
	// В текущей реализации не требует закрытия
}
