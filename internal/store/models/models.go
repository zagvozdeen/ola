package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/zagvozdeen/ola/internal/store/enums"
)

type User struct {
	ID        int            `json:"id"`
	TID       *int64         `json:"tid"`
	UUID      uuid.UUID      `json:"uuid"`
	FirstName string         `json:"first_name"`
	LastName  *string        `json:"last_name"`
	Username  *string        `json:"username"`
	Email     *string        `json:"email"`
	Phone     *string        `json:"phone"`
	Password  *string        `json:"-"`
	Role      enums.UserRole `json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type File struct {
	//ID         int       `json:"id"`
	UUID       uuid.UUID `json:"uuid"`
	Content    string    `json:"content"`
	Size       int64     `json:"size"`
	MimeType   string    `json:"mime_type"`
	OriginName string    `json:"origin_name"`
	UserID     int       `json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type Product struct {
	ID          int               `json:"id"`
	UUID        uuid.UUID         `json:"uuid"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	PriceFrom   int               `json:"price_from"`
	PriceTo     *int              `json:"price_to"`
	Type        enums.ProductType `json:"type"`
	IsMain      bool              `json:"is_main"`
	FileContent string            `json:"file_content"`
	UserID      int               `json:"user_id"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type Cart struct {
	ID        int       `json:"id"`
	UUID      uuid.UUID `json:"uuid"`
	UserID    *int      `json:"user_id"`
	SessionID *string   `json:"session_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CartItem struct {
	ProductID   int               `json:"product_id"`
	ProductUUID uuid.UUID         `json:"product_uuid"`
	ProductName string            `json:"product_name"`
	PriceFrom   int               `json:"price_from"`
	PriceTo     *int              `json:"price_to,omitempty"`
	Type        enums.ProductType `json:"type"`
	FileContent *string           `json:"file_content,omitempty"`
	Qty         int               `json:"qty"`
}

type Order struct {
	ID        int                 `json:"id"`
	UUID      uuid.UUID           `json:"uuid"`
	Status    enums.RequestStatus `json:"status"`
	Source    enums.OrderSource   `json:"source"`
	Name      string              `json:"name"`
	Phone     string              `json:"phone"`
	Content   string              `json:"content"`
	UserID    *int                `json:"user_id"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
}

type Feedback struct {
	ID        int                 `json:"id"`
	UUID      uuid.UUID           `json:"uuid"`
	Status    enums.RequestStatus `json:"status"`
	Source    enums.OrderSource   `json:"source"`
	Type      enums.FeedbackType  `json:"type"`
	Name      string              `json:"name"`
	Phone     string              `json:"phone"`
	Content   string              `json:"content"`
	UserID    int                 `json:"user_id"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
}

type Category struct {
	ID        int       `json:"id"`
	UUID      uuid.UUID `json:"uuid"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CategoryProduct struct {
	CategoryID int `json:"category_id"`
	ProductID  int `json:"product_id"`
}

type OrderTelegramMessage struct {
	OrderID   int   `json:"order_id"`
	ChatID    int64 `json:"chat_id"`
	MessageID int   `json:"message_id"`
}

type FeedbackTelegramMessage struct {
	FeedbackID int   `json:"feedback_id"`
	ChatID     int64 `json:"chat_id"`
	MessageID  int   `json:"message_id"`
}

type OrderComment struct {
	ID        int
	UUID      uuid.UUID
	Content   string
	OrderID   int
	UserID    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Action struct {
	ID        int
	Content   string
	CreatedAt time.Time
}
