package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"telegram-support-bot/internal/models"
	"time"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, username, first_name, last_name, is_staff)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET
			username = EXCLUDED.username,
			first_name = EXCLUDED.first_name,
			last_name = EXCLUDED.last_name,
			is_staff = EXCLUDED.is_staff,
			updated_at = NOW()
		RETURNING id, username, first_name, last_name, is_staff, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
		user.ID, user.Username, user.FirstName, user.LastName, user.IsStaff,
	).Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName,
		&user.IsStaff, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		slog.Error("Failed to create user", "error", err, "user_id", user.ID)
		return err
	}

	slog.Info("User created", "user_id", user.ID, "is_staff", user.IsStaff)
	return nil
}

func (r *Repository) GetUserByID(ctx context.Context, userID int64) (*models.User, error) {
	query := `
		SELECT id, username, first_name, last_name, is_staff, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user models.User
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID, &user.Username, &user.FirstName, &user.LastName,
		&user.IsStaff, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		slog.Error("Failed to get user", "error", err, "user_id", userID)
		return nil, err
	}

	return &user, nil
}

func (r *Repository) CreateMessage(ctx context.Context, message *models.Message) error {
	query := `
        INSERT INTO messages (id, userId, content, is_read, created_at)
        VALUES ($1, $2, $3, false, NOW())
        RETURNING id, userId, content, is_read, created_at
    `

	err := r.db.QueryRowContext(ctx, query, message.ID, message.UserID, message.Content).Scan(
		&message.ID, &message.UserID, &message.Content, &message.IsRead, &message.CreatedAt,
	)
	if err != nil {
		slog.Error("Failed to create message", "error", err, "user_id", message.UserID)
		return err
	}

	slog.Info("Message created", "message_id", message.ID, "user_id", message.UserID)
	return nil
}

func (r *Repository) GetMessages(ctx context.Context, limit int, offset int) ([]models.Message, error) {
	query := `
		SELECT id, userId, content, created_at, is_read
		FROM messages
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		slog.Error("Failed to get messages", "error", err)
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.ID, &msg.UserID, &msg.Content, &msg.CreatedAt, &msg.IsRead); err != nil {
			slog.Error("Failed to scan message", "error", err)
			continue
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func (r *Repository) MarkMessagesAsRead(ctx context.Context, messageIDs []int64) error {
	if len(messageIDs) == 0 {
		return nil
	}

	placeholders := make([]string, len(messageIDs))
	args := make([]interface{}, len(messageIDs))
	for i, id := range messageIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := `
		UPDATE messages
		SET is_read = true
		WHERE id = ANY(ARRAY[` + strings.Join(placeholders, ",") + `])
	`

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		slog.Error("Failed to mark messages as read", "error", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	slog.Info("Messages marked as read", "count", rowsAffected)
	return nil
}

func (r *Repository) CreateSupportResponse(ctx context.Context, response *models.SupportResponse) error {
	query := `
		INSERT INTO support_responses (message_id, staff_id, content, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id, message_id, staff_id, content, created_at
	`

	var id int64
	var messageID int64
	var staffID int64
	var content string
	var createdAt time.Time
	err := r.db.QueryRowContext(ctx, query, response.MessageID, response.StaffID, response.Content).Scan(&id, &messageID, &staffID, &content, &createdAt)
	if err != nil {
		slog.Error("Failed to create support response", "error", err, "message_id", response.MessageID)
		return err
	}

	slog.Info("Support response created", "response_id", id, "message_id", response.MessageID, "staff_id", response.StaffID)
	return nil
}

func (r *Repository) GetMessageWithResponses(ctx context.Context, messageID int64) (*models.MessageWithUser, error) {
	query := `
        SELECT m.id, m.userId, m.content, m.is_read, m.created_at,
               u.id, u.username, u.first_name, u.last_name, u.is_staff, u.created_at, u.updated_at
        FROM messages m
        JOIN users u ON m.userId = u.id
        WHERE m.id = $1
    `

	var result models.MessageWithUser
	err := r.db.QueryRowContext(ctx, query, messageID).Scan(
		&result.Message.ID, &result.Message.UserID, &result.Message.Content,
		&result.Message.IsRead, &result.Message.CreatedAt,
		&result.User.ID, &result.User.Username,
		&result.User.FirstName, &result.User.LastName, &result.User.IsStaff,
		&result.User.CreatedAt, &result.User.UpdatedAt,
	)
	if err != nil {
		slog.Error("Failed to get message with responses", "error", err, "message_id", messageID)
		return nil, err
	}

	// Get responses separately
	responses, err := r.GetResponsesByMessageID(ctx, messageID)
	if err != nil {
		slog.Error("Failed to get responses", "error", err, "message_id", messageID)
		return nil, err
	}
	result.Responses = responses

	return &result, nil
}

func (r *Repository) GetStats(ctx context.Context) (*models.MessageStats, error) {
	query := `
		SELECT 
			(SELECT COUNT(*) FROM messages) as total_messages,
			(SELECT COUNT(*) FROM messages WHERE is_read = false) as unread_messages,
			(SELECT COUNT(*) FROM support_responses) as total_responses,
			(SELECT COUNT(*) FROM users) as total_users
	`

	var stats models.MessageStats
	err := r.db.QueryRowContext(ctx, query).Scan(
		&stats.TotalMessages, &stats.UnreadMessages, &stats.TotalResponses, &stats.TotalUsers,
	)
	if err != nil {
		slog.Error("Failed to get stats", "error", err)
		return nil, err
	}

	return &stats, nil
}

// GetMessagesWithResponses возвращает сообщения с ответами
func (r *Repository) GetMessagesWithResponses(ctx context.Context, limit int, offset int) ([]models.MessageWithUser, error) {
	query := `
		SELECT m.id, m.userId, m.content, m.created_at, m.is_read,
		       u.id, u.username, u.first_name, u.last_name, u.is_staff, u.created_at, u.updated_at
		FROM messages m
		JOIN users u ON m.userId = u.id
		ORDER BY m.created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		slog.Error("Failed to get messages with responses", "error", err)
		return nil, err
	}
	defer rows.Close()

	var results []models.MessageWithUser
	for rows.Next() {
		var result models.MessageWithUser
		err := rows.Scan(
			&result.Message.ID, &result.Message.UserID, &result.Message.Content,
			&result.Message.CreatedAt, &result.Message.IsRead,
			&result.User.ID, &result.User.Username,
			&result.User.FirstName, &result.User.LastName, &result.User.IsStaff,
			&result.User.CreatedAt, &result.User.UpdatedAt,
		)
		if err != nil {
			slog.Error("Failed to scan message", "error", err)
			continue
		}
		results = append(results, result)
	}

	return results, nil
}

// GetMessageByID возвращает сообщение по ID
func (r *Repository) GetMessageByID(ctx context.Context, messageID int64) (*models.MessageWithUser, error) {
	return r.GetMessageWithResponses(ctx, messageID)
}

// GetResponsesByMessageID возвращает ответы на сообщение
func (r *Repository) GetResponsesByMessageID(ctx context.Context, messageID int64) ([]models.SupportResponse, error) {
	query := `
		SELECT id, message_id, staff_id, content, created_at
		FROM support_responses
		WHERE message_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, messageID)
	if err != nil {
		slog.Error("Failed to get responses", "error", err)
		return nil, err
	}
	defer rows.Close()

	var responses []models.SupportResponse
	for rows.Next() {
		var resp models.SupportResponse
		if err := rows.Scan(&resp.ID, &resp.MessageID, &resp.StaffID, &resp.Content, &resp.CreatedAt); err != nil {
			slog.Error("Failed to scan response", "error", err)
			continue
		}
		responses = append(responses, resp)
	}

	return responses, nil
}

// GetUnreadMessages возвращает непрочитанные сообщения
func (r *Repository) GetUnreadMessages(ctx context.Context) ([]models.MessageWithUser, error) {
	query := `
		SELECT m.id, m.userId, m.content, m.created_at, m.is_read,
		       u.id, u.username, u.first_name, u.last_name, u.is_staff, u.created_at, u.updated_at
		FROM messages m
		JOIN users u ON m.userId = u.id
		WHERE m.is_read = false
		ORDER BY m.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		slog.Error("Failed to get unread messages", "error", err)
		return nil, err
	}
	defer rows.Close()

	var results []models.MessageWithUser
	for rows.Next() {
		var result models.MessageWithUser
		err := rows.Scan(
			&result.Message.ID, &result.Message.UserID, &result.Message.Content,
			&result.Message.CreatedAt, &result.Message.IsRead,
			&result.User.ID, &result.User.Username,
			&result.User.FirstName, &result.User.LastName, &result.User.IsStaff,
			&result.User.CreatedAt, &result.User.UpdatedAt,
		)
		if err != nil {
			slog.Error("Failed to scan message", "error", err)
			continue
		}
		results = append(results, result)
	}

	return results, nil
}

// MarkMessageAsRead помечает сообщение как прочитанное
func (r *Repository) MarkMessageAsRead(ctx context.Context, messageID int64) error {
	query := `
		UPDATE messages
		SET is_read = true
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, messageID)
	if err != nil {
		slog.Error("Failed to mark message as read", "error", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	slog.Info("Message marked as read", "message_id", messageID, "rows_affected", rowsAffected)
	return nil
}

// GetUserByMessageID возвращает пользователя по ID сообщения
func (r *Repository) GetUserByMessageID(ctx context.Context, messageID int64) (*models.User, error) {
	query := `
        SELECT u.id, u.username, u.first_name, u.last_name, u.is_staff, u.created_at, u.updated_at
        FROM users u
        JOIN messages m ON m.userId = u.id
        WHERE m.id = $1
    `

	var user models.User
	err := r.db.QueryRowContext(ctx, query, messageID).Scan(
		&user.ID, &user.Username, &user.FirstName, &user.LastName,
		&user.IsStaff, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		slog.Error("Failed to get user by message ID", "error", err, "message_id", messageID)
		return nil, err
	}

	return &user, nil
}

// GetAllMessages возвращает все сообщения
func (r *Repository) GetAllMessages(ctx context.Context, limit int, offset int) ([]models.Message, error) {
	return r.GetMessages(ctx, limit, offset)
}

// GetMessagesByUserID возвращает сообщения пользователя
func (r *Repository) GetMessagesByUserID(ctx context.Context, userID int64) ([]models.Message, error) {
	query := `
		SELECT id, userId, content, created_at, is_read
		FROM messages
		WHERE userId = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		slog.Error("Failed to get messages by user ID", "error", err)
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.ID, &msg.UserID, &msg.Content, &msg.CreatedAt, &msg.IsRead); err != nil {
			slog.Error("Failed to scan message", "error", err)
			continue
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

// GetMessagesByUserIDWithResponses возвращает сообщения пользователя с данными пользователя и ответами
func (r *Repository) GetMessagesByUserIDWithResponses(ctx context.Context, userID int64) ([]models.MessageWithUser, error) {
	user, err := r.GetUserByID(ctx, userID)
	if err != nil || user == nil {
		return nil, err
	}
	messages, err := r.GetMessagesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	result := make([]models.MessageWithUser, 0, len(messages))
	for _, msg := range messages {
		responses, _ := r.GetResponsesByMessageID(ctx, msg.ID)
		result = append(result, models.MessageWithUser{
			Message:   msg,
			User:      *user,
			Responses: responses,
		})
	}
	return result, nil
}
