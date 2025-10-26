// service_test.go
package benchmarks

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pozedorum/WB_project_4/task5/internal/models"
	"github.com/pozedorum/WB_project_4/task5/internal/service"
	"github.com/pozedorum/WB_project_4/task5/internal/utils"
)

// Устанавливаем глобально для всех тестов в этом пакете
func TestMain(m *testing.M) {
	gin.SetMode(gin.ReleaseMode)
	m.Run()
}

// Mock репозиторий для тестирования с поддержкой батчинга
type mockRepo struct {
	urls        map[string]*models.ShortURL
	clicks      []*models.ClickAnalyticsEntry
	batchClicks []models.ClickTask
}

func newMockRepo() *mockRepo {
	return &mockRepo{
		urls:        make(map[string]*models.ShortURL),
		clicks:      make([]*models.ClickAnalyticsEntry, 0),
		batchClicks: make([]models.ClickTask, 0),
	}
}

func (m *mockRepo) CreateShortURL(ctx context.Context, n *models.ShortURL) error {
	m.urls[n.ShortCode] = n
	return nil
}

func (m *mockRepo) GetOriginalURLIfExists(ctx context.Context, shortCode string) (*models.ShortURL, error) {
	if url, exists := m.urls[shortCode]; exists {
		return url, nil
	}
	return nil, models.ErrShortURLNotFound
}

func (m *mockRepo) GetStatisticsByShortCode(ctx context.Context, shortCode string, period string, groupBy string) (*models.AnalyticsResponse, error) {
	return &models.AnalyticsResponse{
		TotalClicks:    100,
		UniqueVisitors: 50,
	}, nil
}

func (m *mockRepo) RegisterClick(ctx context.Context, click *models.ClickAnalyticsEntry) error {
	m.clicks = append(m.clicks, click)
	return nil
}

func (m *mockRepo) BatchRegisterClicks(ctx context.Context, clicks []models.ClickTask) error {
	m.batchClicks = append(m.batchClicks, clicks...)

	// Обновляем счетчики кликов для каждого short_code
	counts := make(map[string]int)
	for _, click := range clicks {
		counts[click.ShortCode]++
	}

	for shortCode, count := range counts {
		if url, exists := m.urls[shortCode]; exists {
			url.ClicksCount += count
		}
	}

	return nil
}

func BenchmarkCreateShortURL(b *testing.B) {
	repo := newMockRepo()
	svc := service.New(repo)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		url := fmt.Sprintf("https://example.com/%d", i)
		_, _ = svc.CreateShortURL(ctx, url, "")
	}

	svc.Stop()
}

func BenchmarkCreateShortURLWithCollision(b *testing.B) {
	repo := newMockRepo()
	svc := service.New(repo)
	ctx := context.Background()

	// Создаем URL, которые будут генерировать коллизии
	baseURL := "https://example.com/same-content"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Добавляем параметр для создания разных URL с одинаковым хэшем
		url := fmt.Sprintf("%s?param=%d", baseURL, i)
		_, _ = svc.CreateShortURL(ctx, url, "")
	}

	svc.Stop()
}

func BenchmarkGenerateShortURL(b *testing.B) {
	urls := []string{
		"https://example.com/page1",
		"https://example.com/page2",
		"https://example.com/page3",
		"https://google.com/search",
		"https://github.com/user/repo",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = utils.GenerateShortURL(urls[i%len(urls)])
	}
}

func BenchmarkGenerateShortURLWithSalt(b *testing.B) {
	urls := []string{
		"https://example.com/page1",
		"https://example.com/page2",
	}

	salts := []string{"1", "2", "3", "4", "5"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = utils.GenerateShortURLWithSalt(urls[i%len(urls)], salts[i%len(salts)])
	}
}

func BenchmarkServiceRedirect(b *testing.B) {
	repo := newMockRepo()
	svc := service.New(repo)
	ctx := context.Background()

	// Сначала создаем тестовую ссылку
	shortURL, err := svc.CreateShortURL(ctx, "https://example.com", "")
	if err != nil {
		b.Fatalf("Failed to create short URL: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = svc.Redirect(ctx, shortURL.ShortCode, "Mozilla/5.0", "192.168.1.1")
	}

	// Даем время батчеру обработать клики
	time.Sleep(100 * time.Millisecond)
	svc.Stop()
}

func BenchmarkServiceRedirectParallel(b *testing.B) {
	repo := newMockRepo()
	svc := service.New(repo)
	ctx := context.Background()

	// Создаем несколько ссылок для параллельного тестирования
	shortURLs := make([]string, 10)
	for i := 0; i < 10; i++ {
		url, err := svc.CreateShortURL(ctx, fmt.Sprintf("https://example.com/%d", i), "")
		if err != nil {
			b.Fatalf("Failed to create short URL %d: %v", i, err)
		}
		shortURLs[i] = url.ShortCode
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			url := shortURLs[i%len(shortURLs)]
			_, _ = svc.Redirect(ctx, url, "Mozilla/5.0", "192.168.1.1")
			i++
		}
	})

	time.Sleep(100 * time.Millisecond)
	svc.Stop()
}
