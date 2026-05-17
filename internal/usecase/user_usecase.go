package usecase

import (
	"context"
	"errors"

	"github.com/azharf99/url-shortener-api/internal/domain"
	"github.com/azharf99/url-shortener-api/internal/utils"
)

type userUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(userRepo domain.UserRepository) domain.UserUsecase {
	return &userUsecase{userRepo}
}

func (u *userUsecase) Register(ctx context.Context, username, email, password string) error {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	user := &domain.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
		Role:     domain.RoleUser,
	}

	return u.userRepo.Create(ctx, user)
}

func (u *userUsecase) Login(ctx context.Context, username, password string) (string, error) {
	user, err := u.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	return utils.GenerateToken(user.ID, user.Role)
}
