package service

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/pozedorum/WB_project_4/task5/internal/models"

	"github.com/pozedorum/WB_project_4/task5/internal/utils"
	"github.com/pozedorum/wbf/zlog"
)

const (
	attemptsCount = 5
)

type ShortURLService struct {
	repo         Repository
	clickBatcher *ClickBatcher
}

func New(repo Repository) *ShortURLService {
	zlog.Logger.Info().Msg("Creating short url service")
	batcher := NewClickBatcher(repo, 100, 1*time.Second)
	batcher.Start()
	return &ShortURLService{repo: repo, clickBatcher: batcher}
}

func (s *ShortURLService) CreateShortURL(ctx context.Context, originalURL string, customCode string) (*models.ShortURL, error) {
	if err := validateURL(originalURL); err != nil {
		return nil, err
	}

	var shortCode string

	if customCode != "" {
		shortCode = customCode
		existingURL, err := s.repo.GetOriginalURLIfExists(ctx, shortCode)
		if err == nil {
			if existingURL.OriginalURL != originalURL {
				// Кастомный код уже существует и связан с другим URL
				return nil, models.ErrDuplicateShortCode
			} else {
				// Кастомный код существует и связан с правильным URL
				return existingURL, nil
			}
		}
		if !errors.Is(err, models.ErrShortURLNotFound) {
			return nil, err
		}

	} else {
		shortCode = utils.GenerateShortURL(originalURL)
		uniqueShortCode, ok, err := s.ensureUniqueShortCode(ctx, originalURL, shortCode)
		if err != nil {
			return nil, err
		}
		if ok {
			shortURL, err := s.repo.GetOriginalURLIfExists(ctx, shortCode)
			if err != nil {
				return nil, err
			}
			return shortURL, nil
		}
		shortCode = uniqueShortCode
	}

	shortURL := &models.ShortURL{
		ShortCode:   shortCode,
		OriginalURL: originalURL,
		CreatedAt:   time.Now(),
		ClicksCount: 0,
	}

	if err := s.repo.CreateShortURL(ctx, shortURL); err != nil {
		if errors.Is(err, models.ErrDuplicateShortCode) {
			return nil, models.ErrDuplicateShortCode
		}
		return nil, err
	}

	zlog.Logger.Info().
		Str("short_code", shortCode).
		Msg("Short URL created")

	return shortURL, nil
}

func (s *ShortURLService) Stop() {
	s.clickBatcher.Stop()
	zlog.Logger.Info().Msg("Click batcher stopped")
}

func (s *ShortURLService) Redirect(ctx context.Context, shortCode string, userAgent, ip string) (string, error) {
	// Синхронно получаем URL для редиректа
	shortURL, err := s.repo.GetOriginalURLIfExists(ctx, shortCode)
	if err != nil {
		if errors.Is(err, models.ErrShortURLNotFound) {
			return "", models.ErrShortURLNotFound
		}
		return "", err
	}

	// Асинхронно добавляем клик в батчер (не блокирует ответ)
	clickTask := models.ClickTask{
		ShortCode: shortCode,
		UserAgent: userAgent,
		IPAddress: ip,
		CreatedAt: time.Now(),
	}

	if !s.clickBatcher.AddClick(clickTask) {
		// Логируем только если очередь переполнена, но не прерываем редирект
		zlog.Logger.Warn().Str("short_code", shortCode).Msg("Click queue full, analytics skipped")
	}

	zlog.Logger.Debug().
		Str("short_code", shortCode).
		Str("original_url", shortURL.OriginalURL).
		Msg("Redirect processed")

	return shortURL.OriginalURL, nil
}
func (s *ShortURLService) GetStatByShortCode(ctx context.Context, shortCode string, period string, groupBy string) (*models.AnalyticsResponse, error) {
	// Проверяем существование ссылки
	_, err := s.repo.GetOriginalURLIfExists(ctx, shortCode)
	if err != nil {
		if errors.Is(err, models.ErrShortURLNotFound) {
			return nil, models.ErrShortURLNotFound
		}
		return nil, err
	}

	return s.repo.GetStatisticsByShortCode(ctx, shortCode, period, groupBy)
}

func (s *ShortURLService) ensureUniqueShortCode(ctx context.Context, originalURL, baseShortCode string) (string, bool, error) {
	// 1. Атомарно проверяем существование и получаем данные
	existingShortURL, err := s.repo.GetOriginalURLIfExists(ctx, baseShortCode)
	if err != nil {
		if errors.Is(err, models.ErrShortURLNotFound) {
			// Код свободен - возвращаем как есть
			return baseShortCode, false, nil
		}
		return "", false, fmt.Errorf("failed to check short code existence: %w", err)
	}

	// 2. Если URL совпадает - возвращаем существующий код
	if existingShortURL.OriginalURL == originalURL {
		zlog.Logger.Info().Str("short_code", baseShortCode).Msg("Returning existing short code for same URL")
		return baseShortCode, true, nil
	}

	// 3. Разные URL с одинаковым хэшем - коллизия!
	zlog.Logger.Warn().
		Str("base_short_code", baseShortCode).
		Str("existing_url", existingShortURL.OriginalURL).
		Str("new_url", originalURL).
		Msg("Hash collision detected, generating new code")

	// 4. Генерируем новые коды с солью
	for attempt := 1; attempt <= attemptsCount; attempt++ {
		saltedShortCode := utils.GenerateShortURLWithSalt(originalURL, fmt.Sprintf("%d", attempt))

		// Атомарно проверяем новый код
		existingSaltedURL, err := s.repo.GetOriginalURLIfExists(ctx, saltedShortCode)
		if err != nil {
			if errors.Is(err, models.ErrShortURLNotFound) {
				// Код свободен - используем
				zlog.Logger.Info().
					Str("new_short_code", saltedShortCode).
					Int("attempt", attempt).
					Msg("Generated unique short code after collision")
				return saltedShortCode, false, nil
			}
			return "", false, fmt.Errorf("failed to check salted short code: %w", err)
		}

		// Если нашли существующий URL - проверяем совпадение
		if existingSaltedURL.OriginalURL == originalURL {
			zlog.Logger.Info().
				Str("short_code", saltedShortCode).
				Msg("Found existing short code for the same URL")
			return saltedShortCode, false, nil
		}
	}

	return "", false, fmt.Errorf("failed to generate unique short code after %d attempts for URL: %s", attemptsCount, originalURL)
}

func validateURL(rawURL string) error {
	if rawURL == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("URL scheme must be http or https")
	}

	if parsedURL.Host == "" {
		return fmt.Errorf("URL must contain host")
	}

	return nil
}
