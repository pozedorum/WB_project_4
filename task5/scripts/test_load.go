package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func main() {
	baseURL := "http://localhost:8080"

	// Тест создания коротких ссылок
	fmt.Println("=== Load Test: URL Shortening ===")
	testShorten(baseURL, 100, 1000) // 100 concurrent, 1000 requests each

	// Тест редиректов
	fmt.Println("\n=== Load Test: URL Redirects ===")
	testRedirect(baseURL, 50, 2000)
}

func testShorten(baseURL string, concurrency, totalRequests int) {
	var successCount int64
	var errorCount int64
	var totalLatency int64

	urls := generateURLs(totalRequests)

	start := time.Now()
	wg := sync.WaitGroup{}

	requestsPerWorker := totalRequests / concurrency

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			client := &http.Client{Timeout: 10 * time.Second}

			for j := 0; j < requestsPerWorker; j++ {
				idx := workerID*requestsPerWorker + j
				if idx >= len(urls) {
					break
				}

				reqBody := ShortenRequest{URL: urls[idx]}
				jsonData, _ := json.Marshal(reqBody)

				startTime := time.Now()
				resp, err := client.Post(baseURL+"/shorten", "application/json", bytes.NewBuffer(jsonData))
				latency := time.Since(startTime).Microseconds()

				atomic.AddInt64(&totalLatency, latency)

				if err != nil || resp.StatusCode != http.StatusAccepted {
					atomic.AddInt64(&errorCount, 1)
					continue
				}

				atomic.AddInt64(&successCount, 1)
				resp.Body.Close()
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(start)

	printResults(successCount, errorCount, totalLatency, duration)
}

func testRedirect(baseURL string, concurrency, totalRequests int) {
	// Сначала создаем тестовые ссылки
	shortURLs := createTestURLs(baseURL, 100)
	if len(shortURLs) == 0 {
		log.Println("No URLs created for redirect test")
		return
	}

	var successCount int64
	var errorCount int64
	var totalLatency int64

	start := time.Now()
	wg := sync.WaitGroup{}

	requestsPerWorker := totalRequests / concurrency

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			client := &http.Client{
				Timeout: 10 * time.Second,
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}

			for j := 0; j < requestsPerWorker; j++ {
				urlIdx := (workerID*requestsPerWorker + j) % len(shortURLs)

				startTime := time.Now()
				resp, err := client.Get(baseURL + shortURLs[urlIdx])
				latency := time.Since(startTime).Microseconds()

				atomic.AddInt64(&totalLatency, latency)

				if err != nil || (resp.StatusCode != http.StatusFound && resp.StatusCode != http.StatusMovedPermanently) {
					atomic.AddInt64(&errorCount, 1)
					continue
				}

				atomic.AddInt64(&successCount, 1)
				resp.Body.Close()
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(start)

	printResults(successCount, errorCount, totalLatency, duration)
}

func generateURLs(count int) []string {
	urls := make([]string, count)
	for i := 0; i < count; i++ {
		urls[i] = fmt.Sprintf("https://example.com/page/%d?param=value%d", i, i)
	}
	return urls
}

func createTestURLs(baseURL string, count int) []string {
	var urls []string
	client := &http.Client{Timeout: 5 * time.Second}

	for i := 0; i < count; i++ {
		reqBody := ShortenRequest{URL: fmt.Sprintf("https://test.com/page%d", i)}
		jsonData, _ := json.Marshal(reqBody)

		resp, err := client.Post(baseURL+"/shorten", "application/json", bytes.NewBuffer(jsonData))
		if err != nil || resp.StatusCode != http.StatusAccepted {
			continue
		}

		var result ShortenResponse
		json.NewDecoder(resp.Body).Decode(&result)
		urls = append(urls, result.ShortURL)
		resp.Body.Close()
	}

	return urls
}

func printResults(success, errors, totalLatency int64, duration time.Duration) {
	fmt.Printf("Duration: %v\n", duration)
	fmt.Printf("Requests: %d\n", success+errors)
	fmt.Printf("Success: %d\n", success)
	fmt.Printf("Errors: %d\n", errors)
	fmt.Printf("RPS: %.2f\n", float64(success)/duration.Seconds())
	if success > 0 {
		fmt.Printf("Avg Latency: %.2f ms\n", float64(totalLatency)/float64(success)/1000.0)
	}
}
