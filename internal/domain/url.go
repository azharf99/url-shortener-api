package domain

import (
	"context"
	"time"
)

type URL struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	OriginalURL string    `json:"original_url" gorm:"not null"`
	ShortCode   string    `json:"short_code" gorm:"unique;not null"`
	UserID      uint      `json:"user_id" gorm:"not null"`
	User        User      `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type URLRepository interface {
	Create(ctx context.Context, url *URL) error
	GetByShortCode(ctx context.Context, shortCode string) (*URL, error)
	GetByID(ctx context.Context, id uint) (*URL, error)
	Update(ctx context.Context, url *URL) error
	Delete(ctx context.Context, id uint) error
	ListByUserID(ctx context.Context, userID uint) ([]URL, error)
	ListAll(ctx context.Context) ([]URL, error)
}

type URLUsecase interface {
	Shorten(ctx context.Context, userID uint, originalURL string) (*URL, error)
	GetOriginalURL(ctx context.Context, shortCode string) (string, error)
	UpdateURL(ctx context.Context, userID uint, role Role, urlID uint, originalURL string) error
	DeleteURL(ctx context.Context, userID uint, role Role, urlID uint) error
	ListURLs(ctx context.Context, userID uint, role Role) ([]URL, error)
}
