package repository

import (
	"context"

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
	return r.db.WithContext(ctx).Create(url).Error
}

func (r *urlRepository) GetByShortCode(ctx context.Context, shortCode string) (*domain.URL, error) {
	var url domain.URL
	err := r.db.WithContext(ctx).Where("short_code = ?", shortCode).First(&url).Error
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (r *urlRepository) GetByID(ctx context.Context, id uint) (*domain.URL, error) {
	var url domain.URL
	err := r.db.WithContext(ctx).First(&url, id).Error
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (r *urlRepository) Update(ctx context.Context, url *domain.URL) error {
	return r.db.WithContext(ctx).Save(url).Error
}

func (r *urlRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.URL{}, id).Error
}

func (r *urlRepository) ListByUserID(ctx context.Context, userID uint) ([]domain.URL, error) {
	var urls []domain.URL
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&urls).Error
	return urls, err
}

func (r *urlRepository) ListAll(ctx context.Context) ([]domain.URL, error) {
	var urls []domain.URL
	err := r.db.WithContext(ctx).Find(&urls).Error
	return urls, err
}
