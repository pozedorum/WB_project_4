package benchmarks

import (
	"context"
	"testing"

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

// Mock репозиторий для тестирования
type mockRepo struct {
	urls   map[string]*models.ShortURL
	clicks []*models.ClickAnalyticsEntry
}

func newMockRepo() *mockRepo {
	return &mockRepo{
		urls:   make(map[string]*models.ShortURL),
		clicks: make([]*models.ClickAnalyticsEntry, 0),
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

func BenchmarkCreateShortURL(b *testing.B) {
	repo := newMockRepo()
	svc := service.New(repo)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		url := "https://example.com/" + string(rune(i))
		_, _ = svc.CreateShortURL(ctx, url, "")
	}
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
