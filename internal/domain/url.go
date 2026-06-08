package domain

import (
	"context"
	"time"
)

type URL struct {
	ID          uint
	OriginalURL string
	ShortCode   string
	UserID      uint
	Clicks      int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type URLRepository interface {
	Create(ctx context.Context, url *URL) error
	GetByShortCode(ctx context.Context, shortCode string) (*URL, error)
	GetByID(ctx context.Context, id uint) (*URL, error)
	Update(ctx context.Context, url *URL) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, userID uint, role Role, search string, offset, limit int) ([]URL, int64, error)
	IncrementClick(ctx context.Context, shortCode string) error
}

type URLUsecase interface {
	Shorten(ctx context.Context, userID uint, originalURL string, customShortCode string) (*URL, error)
	GetOriginalURL(ctx context.Context, shortCode string) (string, error)
	UpdateURL(ctx context.Context, userID uint, role Role, urlID uint, originalURL string) error
	DeleteURL(ctx context.Context, userID uint, role Role, urlID uint) error
	ListURLs(ctx context.Context, userID uint, role Role, search string, page, limit int) ([]URL, int64, error)
}
