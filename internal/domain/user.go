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
	ID        uint
	Username  string
	Email     string
	Password  string
	Role      Role
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uint) (*User, error)
	GetByUsernameOrEmail(ctx context.Context, identifier string) (*User, error)
	GetAll(ctx context.Context) ([]User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uint) error
}

type UserUsecase interface {
	Register(ctx context.Context, username, email, password string) error
	Login(ctx context.Context, username, password string) (string, error)
	GetAllUsers(ctx context.Context) ([]User, error)
	GetUserByID(ctx context.Context, id uint) (*User, error)
	UpdateUser(ctx context.Context, id uint, username, email string, role Role) error
	DeleteUser(ctx context.Context, id uint) error
	AdminCreateUser(ctx context.Context, username, email, password string, role Role) error
}
