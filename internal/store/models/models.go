package models

import (
	"time"

	"github.com/zagvozdeen/ola/internal/store/enums"
)

type User struct {
	ID        int            `json:"id"`
	TID       *int64         `json:"tid,omitempty"`
	UUID      string         `json:"uuid"`
	FirstName string         `json:"first_name"`
	LastName  *string        `json:"last_name,omitempty"`
	Username  *string        `json:"username,omitempty"`
	Email     *string        `json:"email,omitempty"`
	Password  *string        `json:"password,omitempty"`
	Role      enums.UserRole `json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type File struct {
	ID         int       `json:"id"`
	UUID       string    `json:"uuid"`
	Content    string    `json:"content"`
	Size       int       `json:"size"`
	MimeType   string    `json:"mime_type"`
	OriginName string    `json:"origin_name"`
	UserID     int       `json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type Product struct {
	ID        string    `json:"id"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Service struct {
	ID        string    `json:"id"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Order struct {
	ID        string    `json:"id"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Feedback struct {
	ID        string    `json:"id"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
