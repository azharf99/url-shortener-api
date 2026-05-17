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
	user, err := u.userRepo.GetByUsernameOrEmail(ctx, username)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	return utils.GenerateToken(user.ID, user.Role)
}

func (u *userUsecase) ListUsers(ctx context.Context, search string, page, limit int) ([]domain.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	return u.userRepo.List(ctx, search, offset, limit)
}

func (u *userUsecase) GetUserByID(ctx context.Context, id uint) (*domain.User, error) {
	return u.userRepo.GetByID(ctx, id)
}

func (u *userUsecase) UpdateUser(ctx context.Context, id uint, username, email string, role domain.Role) error {
	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	user.Username = username
	user.Email = email
	user.Role = role

	return u.userRepo.Update(ctx, user)
}

func (u *userUsecase) DeleteUser(ctx context.Context, id uint) error {
	return u.userRepo.Delete(ctx, id)
}

func (u *userUsecase) AdminCreateUser(ctx context.Context, username, email, password string, role domain.Role) error {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	user := &domain.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
		Role:     role,
	}

	return u.userRepo.Create(ctx, user)
}
