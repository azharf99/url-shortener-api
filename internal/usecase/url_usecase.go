package usecase

import (
	"context"
	"errors"
	"math/rand"

	"github.com/azharf99/url-shortener-api/internal/domain"
)

type urlUsecase struct {
	urlRepo domain.URLRepository
}

func NewURLUsecase(urlRepo domain.URLRepository) domain.URLUsecase {
	return &urlUsecase{urlRepo}
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateShortCode(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func (u *urlUsecase) Shorten(ctx context.Context, userID uint, originalURL string) (*domain.URL, error) {
	shortCode := generateShortCode(6)
	url := &domain.URL{
		OriginalURL: originalURL,
		ShortCode:   shortCode,
		UserID:      userID,
	}

	err := u.urlRepo.Create(ctx, url)
	if err != nil {
		return nil, err
	}

	return url, nil
}

func (u *urlUsecase) GetOriginalURL(ctx context.Context, shortCode string) (string, error) {
	url, err := u.urlRepo.GetByShortCode(ctx, shortCode)
	if err != nil {
		return "", err
	}
	if url == nil {
		return "", errors.New("URL not found")
	}
	return url.OriginalURL, nil
}

func (u *urlUsecase) UpdateURL(ctx context.Context, userID uint, role domain.Role, urlID uint, originalURL string) error {
	url, err := u.urlRepo.GetByID(ctx, urlID)
	if err != nil {
		return err
	}
	if url == nil {
		return errors.New("URL not found")
	}

	if role != domain.RoleAdmin && url.UserID != userID {
		return errors.New("unauthorized: you do not own this URL")
	}

	url.OriginalURL = originalURL
	return u.urlRepo.Update(ctx, url)
}

func (u *urlUsecase) DeleteURL(ctx context.Context, userID uint, role domain.Role, urlID uint) error {
	url, err := u.urlRepo.GetByID(ctx, urlID)
	if err != nil {
		return err
	}
	if url == nil {
		return errors.New("URL not found")
	}

	if role != domain.RoleAdmin && url.UserID != userID {
		return errors.New("unauthorized: you do not own this URL")
	}

	return u.urlRepo.Delete(ctx, urlID)
}

func (u *urlUsecase) ListURLs(ctx context.Context, userID uint, role domain.Role) ([]domain.URL, error) {
	if role == domain.RoleAdmin {
		return u.urlRepo.ListAll(ctx)
	}
	return u.urlRepo.ListByUserID(ctx, userID)
}
