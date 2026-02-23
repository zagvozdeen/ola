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
	Password  *string        `json:"-"`
	Role      enums.UserRole `json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type File struct {
	ID         int       `json:"id"`
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
	PriceTo     *int              `json:"price_to,omitempty"`
	Type        enums.ProductType `json:"type"`
	FileID      int               `json:"file_id"`
	FileContent *string           `json:"file_content,omitempty"`
	UserID      int               `json:"user_id"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

//type Service struct {
//	ID          int       `json:"id"`
//	UUID        string    `json:"uuid"`
//	Name        string    `json:"name"`
//	Description string    `json:"description"`
//	PriceFrom   int       `json:"price_from"`
//	PriceTo     *int      `json:"price_to,omitempty"`
//	FileID      int       `json:"file_id"`
//	FileContent *string   `json:"file_content,omitempty"`
//	UserID      int       `json:"user_id"`
//	CreatedAt   time.Time `json:"created_at"`
//	UpdatedAt   time.Time `json:"updated_at"`
//}

type Review struct {
	ID          int       `json:"id"`
	UUID        uuid.UUID `json:"uuid"`
	Name        string    `json:"name"`
	Content     string    `json:"content"`
	FileID      int       `json:"file_id"`
	FileContent *string   `json:"file_content,omitempty"`
	UserID      int       `json:"user_id"`
	PublishedAt time.Time `json:"published_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Order struct {
	ID        int       `json:"id"`
	UUID      string    `json:"uuid"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Content   string    `json:"content"`
	UserID    *int      `json:"user_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Feedback struct {
	ID        int       `json:"id"`
	UUID      uuid.UUID `json:"uuid"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Content   string    `json:"content"`
	UserID    *int      `json:"user_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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

type CategoryService struct {
	CategoryID int `json:"category_id"`
	ServiceID  int `json:"service_id"`
}
