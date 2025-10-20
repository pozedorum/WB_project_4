package archiver

import (
	"context"
	"time"

	"github.com/pozedorum/WB_project_4/task3/internal/interfaces"
)

type ArchiveCleaner struct {
	repo     interfaces.Repository
	logger   interfaces.Logger
	interval time.Duration
}

func NewArchiveCleaner(repo interfaces.Repository, logger interfaces.Logger, interval time.Duration) interfaces.Archiver {
	return &ArchiveCleaner{
		repo:     repo,
		logger:   logger,
		interval: interval,
	}
}

// Start запускает фоновую очистку архива
func (c *ArchiveCleaner) Start(ctx context.Context) {
	c.logger.Info("ARCHIVER_START", "Starting archive cleaner",
		"interval_hours", c.interval.Hours())
	go func() {
		ticker := time.NewTicker(c.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				c.logger.Info("ARCHIVER_STOP", "Stopping archive cleaner")
				return
			case <-ticker.C:
				go c.checkOldEvents()
			}
		}
	}()
}

// checkOldEvents архивирует старые события
func (c *ArchiveCleaner) checkOldEvents() {
	start := time.Now()
	c.logger.Debug("ARCHIVER_RUN", "Starting archive cleanup")

	// Определяем временную границу для архивации
	archiveThreshold := time.Now().Add(-1 * time.Minute)

	deletedCount, err := c.repo.ArchiveOldEvents(archiveThreshold)
	if err != nil {
		c.logger.Error("ARCHIVER_CLEANUP", "Failed to cleanup archived events",
			"error", err,
			"threshold", archiveThreshold)
		return
	}

	duration := time.Since(start)
	if deletedCount > 0 {
		c.logger.Info("ARCHIVER_SUCCESS", "Archive cleanup completed",
			"deleted_count", deletedCount,
			"threshold", archiveThreshold,
			"duration_ms", duration.Milliseconds())
	} else {
		c.logger.Debug("ARCHIVER_NO_EVENTS", "No events to cleanup",
			"threshold", archiveThreshold,
			"duration_ms", duration.Milliseconds())
	}
}
