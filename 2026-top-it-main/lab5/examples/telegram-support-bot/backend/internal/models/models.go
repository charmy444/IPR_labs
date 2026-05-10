package models

import (
	"time"
)

type User struct {
	ID        int64     `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name" db:"last_name"`
	IsStaff   bool      `json:"is_staff" db:"is_staff"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Message struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"userId" db:"userId"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	IsRead    bool      `json:"is_read" db:"is_read"`
}

type SupportResponse struct {
	ID        int64     `json:"id" db:"id"`
	MessageID int64     `json:"message_id" db:"message_id"`
	StaffID   int64     `json:"staff_id" db:"staff_id"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type MessageWithUser struct {
	Message
	User      User              `json:"user"`
	Responses []SupportResponse `json:"responses"`
}

type MessageStats struct {
	TotalMessages  int64 `json:"total_messages"`
	UnreadMessages int64 `json:"unread_messages"`
	TotalResponses int64 `json:"total_responses"`
	TotalUsers     int64 `json:"total_users"`
}
