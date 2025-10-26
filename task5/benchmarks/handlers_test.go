// handlers_test.go
package benchmarks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pozedorum/WB_project_4/task5/internal/models"
	"github.com/pozedorum/WB_project_4/task5/internal/server"
	"github.com/pozedorum/WB_project_4/task5/internal/service"
	"github.com/pozedorum/wbf/ginext"
)

func BenchmarkShortenHandler(b *testing.B) {
	gin.SetMode(gin.ReleaseMode)
	repo := newMockRepo()
	svc := service.New(repo)
	srv := server.New(svc)

	router := ginext.New()
	srv.SetupRoutes(router.Group("/"))

	requestBody := map[string]string{
		"url": "https://example.com/benchmark",
	}
	body, _ := json.Marshal(requestBody)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/shorten", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)
	}

	// Останавливаем батчер после теста
	svc.Stop()
}

func BenchmarkRedirectHandler(b *testing.B) {
	gin.SetMode(gin.ReleaseMode)

	repo := newMockRepo()
	svc := service.New(repo)
	srv := server.New(svc)

	router := ginext.New()
	srv.SetupRoutes(router.Group("/"))

	// Сначала создаем ссылку
	createReqBase := map[string]string{"url": "https://example.com"}
	createBody, _ := json.Marshal(createReqBase)
	createReq := httptest.NewRequest("POST", "/shorten", bytes.NewReader(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createResp := httptest.NewRecorder()
	router.ServeHTTP(createResp, createReq)

	var result map[string]string
	if err := json.Unmarshal(createResp.Body.Bytes(), &result); err != nil {
		b.Fatalf("Failed to parse response: %v", err)
	}

	shortURL := result["short_url"]
	if shortURL == "" {
		b.Fatal("Failed to create short URL")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", shortURL, nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)
	}

	// Даем время батчеру обработать оставшиеся клики
	time.Sleep(100 * time.Millisecond)
	svc.Stop()
}

func BenchmarkRedirectHandlerWithBatcher(b *testing.B) {
	gin.SetMode(gin.ReleaseMode)

	repo := newMockRepo()
	svc := service.New(repo)
	srv := server.New(svc)

	router := ginext.New()
	srv.SetupRoutes(router.Group("/"))

	// Создаем несколько ссылок для тестирования
	shortURLs := make([]string, 10)
	for i := 0; i < 10; i++ {
		createReqBase := map[string]string{"url": fmt.Sprintf("https://example.com/%d", i)}
		createBody, _ := json.Marshal(createReqBase)
		createReq := httptest.NewRequest("POST", "/shorten", bytes.NewReader(createBody))
		createReq.Header.Set("Content-Type", "application/json")
		createResp := httptest.NewRecorder()
		router.ServeHTTP(createResp, createReq)

		var result map[string]string
		if err := json.Unmarshal(createResp.Body.Bytes(), &result); err != nil {
			b.Fatalf("Failed to parse response for URL %d: %v", i, err)
		}

		shortURL := result["short_url"]
		if shortURL == "" {
			b.Fatalf("Failed to create short URL for %d", i)
		}
		shortURLs[i] = shortURL
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		url := shortURLs[i%len(shortURLs)]
		req := httptest.NewRequest("GET", url, nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)
	}

	// Даем время батчеру обработать оставшиеся клики
	time.Sleep(100 * time.Millisecond)
	svc.Stop()
}

func BenchmarkClickBatcher(b *testing.B) {
	repo := newMockRepo()
	batcher := service.NewClickBatcher(repo, 100, 100*time.Millisecond)
	batcher.Start()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		click := models.ClickTask{
			ShortCode: "test123",
			UserAgent: "Mozilla/5.0",
			IPAddress: "192.168.1.1",
			CreatedAt: time.Now(),
		}
		batcher.AddClick(click)
	}

	// Ждем обработки всех батчей
	time.Sleep(200 * time.Millisecond)
	batcher.Stop()
}

func BenchmarkBatchRegisterClicks(b *testing.B) {
	repo := newMockRepo()
	ctx := context.Background()

	// Тестируем только маленькие батчи
	b.Run("BatchSize-5", func(b *testing.B) {
		clicks := make([]models.ClickTask, 5)
		for i := 0; i < 5; i++ {
			clicks[i] = models.ClickTask{
				ShortCode: fmt.Sprintf("test%d", i),
				UserAgent: "Mozilla/5.0",
				IPAddress: fmt.Sprintf("192.168.1.%d", i+1),
				CreatedAt: time.Now(),
			}
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			repo.BatchRegisterClicks(ctx, clicks)
		}
	})
}
