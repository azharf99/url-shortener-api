package repository

import (
	"time"

	"github.com/azharf99/url-shortener-api/internal/domain"
)

type UserModel struct {
	ID              uint        `gorm:"primaryKey"`
	Username        string      `gorm:"unique;not null"`
	Email           string      `gorm:"unique;not null"`
	Password        string      `gorm:"default:null"`
	Role            domain.Role `gorm:"type:varchar(20);default:'user'"`
	GoogleID        string      `gorm:"unique;index;default:null"`
	IsPremium       bool        `gorm:"default:false"`
	SubscriptionEnd time.Time   `gorm:"default:null"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (UserModel) TableName() string {
	return "users"
}

type URLModel struct {
	ID          uint   `gorm:"primaryKey"`
	OriginalURL string `gorm:"not null"`
	ShortCode   string `gorm:"unique;not null"`
	UserID      uint   `gorm:"not null"`
	Clicks      int64  `gorm:"default:0"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (URLModel) TableName() string {
	return "urls"
}

// Mappers
func ToUserEntity(m *UserModel) *domain.User {
	if m == nil {
		return nil
	}
	return &domain.User{
		ID:              m.ID,
		Username:        m.Username,
		Email:           m.Email,
		Password:        m.Password,
		Role:            m.Role,
		GoogleID:        m.GoogleID,
		IsPremium:       m.IsPremium,
		SubscriptionEnd: m.SubscriptionEnd,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

func FromUserEntity(e *domain.User) *UserModel {
	if e == nil {
		return nil
	}
	return &UserModel{
		ID:              e.ID,
		Username:        e.Username,
		Email:           e.Email,
		Password:        e.Password,
		Role:            e.Role,
		GoogleID:        e.GoogleID,
		IsPremium:       e.IsPremium,
		SubscriptionEnd: e.SubscriptionEnd,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
	}
}

func ToURLEntity(m *URLModel) *domain.URL {
	if m == nil {
		return nil
	}
	return &domain.URL{
		ID:          m.ID,
		OriginalURL: m.OriginalURL,
		ShortCode:   m.ShortCode,
		UserID:      m.UserID,
		Clicks:      m.Clicks,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func FromURLEntity(e *domain.URL) *URLModel {
	if e == nil {
		return nil
	}
	return &URLModel{
		ID:          e.ID,
		OriginalURL: e.OriginalURL,
		ShortCode:   e.ShortCode,
		UserID:      e.UserID,
		Clicks:      e.Clicks,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}
