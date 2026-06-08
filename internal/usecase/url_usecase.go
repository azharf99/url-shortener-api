package usecase

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/azharf99/url-shortener-api/internal/domain"
)

type urlUsecase struct {
	urlRepo  domain.URLRepository
	userRepo domain.UserRepository
}

func NewURLUsecase(urlRepo domain.URLRepository, userRepo domain.UserRepository) domain.URLUsecase {
	return &urlUsecase{urlRepo, userRepo}
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateShortCode(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func (u *urlUsecase) Shorten(ctx context.Context, userID uint, originalURL string, customShortCode string) (*domain.URL, error) {
	var shortCode string

	if customShortCode != "" {
		// Fetch user to verify premium status
		user, err := u.userRepo.GetByID(ctx, userID)
		if err != nil {
			return nil, err
		}
		if user == nil {
			return nil, errors.New("user not found")
		}

		isPremium := user.IsPremium && user.SubscriptionEnd.After(time.Now())
		if !isPremium {
			return nil, errors.New("custom shortcode is only available for premium subscribers")
		}

		if len(customShortCode) > 20 {
			return nil, errors.New("custom shortcode exceeds maximum length of 20 characters")
		}
		
		existing, err := u.urlRepo.GetByShortCode(ctx, customShortCode)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, errors.New("custom shortcode already exists")
		}
		shortCode = customShortCode
	} else {
		shortCode = generateShortCode(6)
	}

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
		return "", errors.New("url not found")
	}

	// Increment click count (fire and forget for performance, or handle error)
	_ = u.urlRepo.IncrementClick(ctx, shortCode)

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

func (u *urlUsecase) ListURLs(ctx context.Context, userID uint, role domain.Role, search string, page, limit int) ([]domain.URL, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	return u.urlRepo.List(ctx, userID, role, search, offset, limit)
}
