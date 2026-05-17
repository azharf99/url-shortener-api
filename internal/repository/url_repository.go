package repository

import (
	"context"
	"errors"

	"github.com/azharf99/url-shortener-api/internal/domain"
	"gorm.io/gorm"
)

type urlRepository struct {
	db *gorm.DB
}

func NewURLRepository(db *gorm.DB) domain.URLRepository {
	return &urlRepository{db}
}

func (r *urlRepository) Create(ctx context.Context, url *domain.URL) error {
	model := FromURLEntity(url)
	err := r.db.WithContext(ctx).Create(model).Error
	if err == nil {
		url.ID = model.ID
		url.CreatedAt = model.CreatedAt
		url.UpdatedAt = model.UpdatedAt
	}
	return err
}

func (r *urlRepository) GetByShortCode(ctx context.Context, shortCode string) (*domain.URL, error) {
	var model URLModel
	err := r.db.WithContext(ctx).Where("short_code = ?", shortCode).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return ToURLEntity(&model), nil
}

func (r *urlRepository) GetByID(ctx context.Context, id uint) (*domain.URL, error) {
	var model URLModel
	err := r.db.WithContext(ctx).First(&model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return ToURLEntity(&model), nil
}

func (r *urlRepository) Update(ctx context.Context, url *domain.URL) error {
	model := FromURLEntity(url)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *urlRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&URLModel{}, id).Error
}

func (r *urlRepository) ListByUserID(ctx context.Context, userID uint) ([]domain.URL, error) {
	var models []URLModel
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&models).Error
	if err != nil {
		return nil, err
	}

	urls := make([]domain.URL, len(models))
	for i, m := range models {
		urls[i] = *ToURLEntity(&m)
	}
	return urls, nil
}

func (r *urlRepository) ListAll(ctx context.Context) ([]domain.URL, error) {
	var models []URLModel
	err := r.db.WithContext(ctx).Find(&models).Error
	if err != nil {
		return nil, err
	}

	urls := make([]domain.URL, len(models))
	for i, m := range models {
		urls[i] = *ToURLEntity(&m)
	}
	return urls, nil
}

func (r *urlRepository) IncrementClick(ctx context.Context, shortCode string) error {
	return r.db.WithContext(ctx).Model(&URLModel{}).
		Where("short_code = ?", shortCode).
		Update("clicks", gorm.Expr("clicks + ?", 1)).Error
}
