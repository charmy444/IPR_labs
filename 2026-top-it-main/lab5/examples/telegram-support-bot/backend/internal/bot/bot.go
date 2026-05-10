package bot

import (
	"context"
	"log/slog"

	"telegram-support-bot/internal/appmetrics"
	"telegram-support-bot/internal/models"
	"telegram-support-bot/internal/repository"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api  *tgbotapi.BotAPI
	repo *repository.Repository
}

func New(api *tgbotapi.BotAPI, repo *repository.Repository) *Bot {
	return &Bot{
		api:  api,
		repo: repo,
	}
}

func (b *Bot) Start() error {
	slog.Info("Starting Telegram bot")

	// Get updates
	updates := b.api.GetUpdatesChan(tgbotapi.UpdateConfig{Timeout: 60})

	for update := range updates {
		if update.Message == nil {
			continue
		}

		// Handle /start command
		if update.Message.Text == "/start" {
			b.handleStartCommand(update.Message)
			continue
		}

		// Handle /help command
		if update.Message.Text == "/help" {
			b.handleHelpCommand(update.Message)
			continue
		}

		// Handle regular messages
		b.handleMessage(update.Message)
	}

	return nil
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) {
	ctx := context.Background()
	user, err := b.repo.GetUserByID(ctx, message.From.ID)
	if err != nil {
		// Create new user
		user = &models.User{
			ID:        int64(message.From.ID),
			Username:  message.From.UserName,
			FirstName: message.From.FirstName,
			LastName:  message.From.LastName,
			IsStaff:   false,
		}

		if err := b.repo.CreateUser(ctx, user); err != nil {
			slog.Error("Failed to create user", "error", err, "user_id", message.From.ID)
			b.sendMessage(message.Chat.ID, "Failed to register you. Please try again.")
			return
		}

		b.sendMessage(message.Chat.ID, "Welcome! You have been registered. Send /help to see available commands.")
		slog.Info("User registered", "user_id", message.From.ID, "username", message.From.UserName)
	} else {
		b.sendMessage(message.Chat.ID, "You are already registered. Send /help to see available commands.")
	}
}

func (b *Bot) handleHelpCommand(message *tgbotapi.Message) {
	helpText := `Available commands:
/start - Register as a new user
/help - Show this help message

Just send any message to create a support ticket.`

	b.sendMessage(message.Chat.ID, helpText)
}

func (b *Bot) handleMessage(message *tgbotapi.Message) {
	// Create user if not exists
	ctx := context.Background()
	user, err := b.repo.GetUserByID(ctx, message.From.ID)
	if err != nil {
		user = &models.User{
			ID:        int64(message.From.ID),
			Username:  message.From.UserName,
			FirstName: message.From.FirstName,
			LastName:  message.From.LastName,
			IsStaff:   false,
		}

		if err := b.repo.CreateUser(ctx, user); err != nil {
			slog.Error("Failed to create user", "error", err, "user_id", message.From.ID)
			b.sendMessage(message.Chat.ID, "Failed to register you. Please try again.")
			return
		}
	}

	// Save message to database
	msg := &models.Message{
		ID:      int64(message.MessageID),
		UserID:  user.ID,
		Content: message.Text,
		IsRead:  false,
	}

	if err := b.repo.CreateMessage(ctx, msg); err != nil {
		slog.Error("Failed to create message", "error", err, "message_id", message.MessageID)
		b.sendMessage(message.Chat.ID, "Failed to save your message. Please try again.")
		return
	}

	appmetrics.SupportTicketsCreated.Inc()
	slog.Info("Message received", "message_id", message.MessageID, "content", message.Text)
	b.sendMessage(message.Chat.ID, "Your message has been received. Support staff will respond shortly.")
}

func (b *Bot) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := b.api.Send(msg); err != nil {
		slog.Error("Failed to send message", "error", err, "chat_id", chatID)
	}
}

// SendMessageToUser отправляет сообщение пользователю (просто ответ, без reply)
func (b *Bot) SendMessageToUser(telegramID int64, text string) (int64, error) {
	msg := tgbotapi.NewMessage(telegramID, text)
	sentMsg, err := b.api.Send(msg)
	if err != nil {
		slog.Error("Failed to send message to user", "error", err, "telegram_id", telegramID)
		return 0, err
	}
	return int64(sentMsg.MessageID), nil
}
