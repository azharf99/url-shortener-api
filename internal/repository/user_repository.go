package repository

import (
	"context"
	"errors"

	"github.com/azharf99/url-shortener-api/internal/domain"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	model := FromUserEntity(user)
	err := r.db.WithContext(ctx).Create(model).Error
	if err == nil {
		user.ID = model.ID
		user.CreatedAt = model.CreatedAt
		user.UpdatedAt = model.UpdatedAt
	}
	return err
}

func (r *userRepository) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	var model UserModel
	err := r.db.WithContext(ctx).First(&model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return ToUserEntity(&model), nil
}

func (r *userRepository) GetByUsernameOrEmail(ctx context.Context, identifier string) (*domain.User, error) {
	var model UserModel
	err := r.db.WithContext(ctx).Where("username = ? OR email = ?", identifier, identifier).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return ToUserEntity(&model), nil
}

func (r *userRepository) GetAll(ctx context.Context) ([]domain.User, error) {
	var models []UserModel
	err := r.db.WithContext(ctx).Find(&models).Error
	if err != nil {
		return nil, err
	}

	users := make([]domain.User, len(models))
	for i, m := range models {
		users[i] = *ToUserEntity(&m)
	}
	return users, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	model := FromUserEntity(user)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&UserModel{}, id).Error
}
