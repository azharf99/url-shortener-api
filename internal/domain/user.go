package domain

import (
	"context"
	"time"
)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type User struct {
	ID              uint
	Username        string
	Email           string
	Password        string
	Role            Role
	GoogleID        string
	IsPremium       bool
	SubscriptionEnd time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uint) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByUsernameOrEmail(ctx context.Context, identifier string) (*User, error)
	GetByGoogleID(ctx context.Context, googleID string) (*User, error)
	List(ctx context.Context, search string, offset, limit int) ([]User, int64, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uint) error
}

type UserUsecase interface {
	Register(ctx context.Context, username, email, password string) error
	Login(ctx context.Context, username, password string) (string, error)
	GoogleLogin(ctx context.Context, googleID, email, name string) (string, error)
	ListUsers(ctx context.Context, search string, page, limit int) ([]User, int64, error)
	GetUserByID(ctx context.Context, id uint) (*User, error)
	UpdateUser(ctx context.Context, id uint, username, email string, role Role) error
	DeleteUser(ctx context.Context, id uint) error
	AdminCreateUser(ctx context.Context, username, email, password string, role Role) error
}
