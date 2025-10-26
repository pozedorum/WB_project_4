package service

import (
	"context"
	"sync"
	"time"

	"github.com/pozedorum/WB_project_4/task5/internal/models"
	"github.com/pozedorum/wbf/zlog"
	"golang.org/x/sync/semaphore"
)

type ClickBatcher struct {
	repo       Repository
	clicks     chan models.ClickTask
	batch      []models.ClickTask
	batchSize  int
	timeout    time.Duration
	mu         sync.Mutex
	done       chan struct{}
	wg         sync.WaitGroup
	processSem *semaphore.Weighted // Ограничиваем параллельную обработку
	maxWorkers int
}

func NewClickBatcher(repo Repository, batchSize int, timeout time.Duration) *ClickBatcher {
	return &ClickBatcher{
		repo:       repo,
		clicks:     make(chan models.ClickTask, 10000),
		batchSize:  batchSize,
		timeout:    timeout,
		done:       make(chan struct{}),
		processSem: semaphore.NewWeighted(10), // Максимум 10 параллельных обработчиков
		maxWorkers: 10,
	}
}

func (b *ClickBatcher) Start() {
	b.wg.Add(1)
	go b.processBatches()
}

func (b *ClickBatcher) Stop() {
	close(b.done)
	b.wg.Wait()
}

func (b *ClickBatcher) AddClick(click models.ClickTask) bool {
	select {
	case b.clicks <- click:
		return true
	case <-b.done:
		return false
	default:
		zlog.Logger.Warn().Str("short_code", click.ShortCode).Msg("Click queue overflow, skipping analytics")
		return false
	}
}

func (b *ClickBatcher) processBatches() {
	defer b.wg.Done()

	ticker := time.NewTicker(b.timeout)
	defer ticker.Stop()

	for {
		select {
		case click := <-b.clicks:
			b.mu.Lock()
			b.batch = append(b.batch, click)
			shouldFlush := len(b.batch) >= b.batchSize
			b.mu.Unlock()

			if shouldFlush {
				b.flushBatch()
			}

		case <-ticker.C:
			b.flushBatch()

		case <-b.done:
			b.flushBatch()
			return
		}
	}
}

func (b *ClickBatcher) flushBatch() {
	b.mu.Lock()
	if len(b.batch) == 0 {
		b.mu.Unlock()
		return
	}

	batch := make([]models.ClickTask, len(b.batch))
	copy(batch, b.batch)
	b.batch = b.batch[:0]
	b.mu.Unlock()

	// Ограничиваем параллельную обработку
	if err := b.processSem.Acquire(context.Background(), 1); err != nil {
		zlog.Logger.Error().Err(err).Msg("Failed to acquire semaphore")
		return
	}

	go func(batch []models.ClickTask) {
		defer b.processSem.Release(1)
		b.processBatch(batch)
	}(batch)
}

func (b *ClickBatcher) processBatch(clicks []models.ClickTask) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := b.repo.BatchRegisterClicks(ctx, clicks); err != nil {
		zlog.Logger.Error().Err(err).Int("batch_size", len(clicks)).Msg("Failed to batch register clicks")
	} else {
		zlog.Logger.Debug().Int("batch_size", len(clicks)).Msg("Successfully processed click batch")
	}
}
