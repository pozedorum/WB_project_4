package analyzer

import (
	"fmt"
	"runtime"
	"time"

	"github.com/pozedorum/WB_project_4/task4/internal/models"
)

type Analyzer struct {
	startTime time.Time
}

func NewAnalyzer() *Analyzer {
	return &Analyzer{startTime: time.Now()}
}

func (a *Analyzer) PrintMemStats() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	fmt.Printf("Allocated Memory: %v bytes\n", memStats.Alloc)
	fmt.Printf("Total Allocated Memory: %v bytes\n", memStats.TotalAlloc)
	fmt.Printf("Heap Memory: %v bytes\n", memStats.HeapAlloc)
	fmt.Printf("Heap System Memory: %v bytes\n", memStats.HeapSys)
	fmt.Printf("Garbage Collector Memory: %v bytes\n", memStats.GCSys)
}

// Collect собирает все метрики системы
func (a *Analyzer) Collect() *models.SystemStats {
	var memStats *runtime.MemStats
	runtime.ReadMemStats(memStats)

	return &models.SystemStats{
		Memory:      a.collectMemoryStats(memStats),
		GC:          a.collectGCStats(memStats),
		Allocations: a.collectAllocationStats(memStats),
		Goroutines:  runtime.NumGoroutine(),
		Timestamp:   time.Now(),
	}
}

func (a *Analyzer) collectMemoryStats(stat *runtime.MemStats) *models.MemStats {
	return &models.MemStats{
		Alloc:       stat.Alloc,
		TotalAlloc:  stat.TotalAlloc,
		Sys:         stat.Sys,
		HeapAlloc:   stat.HeapAlloc,
		HeapInuse:   stat.HeapInuse,
		HeapObjects: stat.HeapObjects,
	}
}

func (a *Analyzer) collectGCStats(stat *runtime.MemStats) *models.GCStats {
	res := &models.GCStats{
		NumGC:         stat.NumGC,
		NumForcedGC:   stat.NumForcedGC,
		LastGC:        time.Unix(0, int64(stat.LastGC)),
		PauseTotal:    time.Duration(stat.PauseTotalNs),
		GCCPUFraction: stat.GCCPUFraction,
	}
	// Длительность последней паузы GC
	if stat.NumGC > 0 {
		lastIndex := (stat.NumGC - 1) % uint32(len(stat.PauseNs))
		res.LastGCPause = time.Duration(stat.PauseNs[lastIndex])
	}
	return res
}

func (a *Analyzer) collectAllocationStats(stat *runtime.MemStats) *models.AllocationStats {
	return &models.AllocationStats{
		Mallocs:     stat.Mallocs,
		Frees:       stat.Frees,
		LiveObjects: stat.Mallocs - stat.Frees,
	}
}
