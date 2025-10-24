package benchmarks

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
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
	json.Unmarshal(createResp.Body.Bytes(), &result)
	shortURL := result["short_url"]

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", shortURL, nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)
	}
}
