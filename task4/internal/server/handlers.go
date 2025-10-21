package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AllocationsHandler возвращает статистику аллокаций
func (serv *Server) AllocationsHandler(c *gin.Context) {
	stats := serv.analyzer.Collect()

	c.JSON(http.StatusOK, gin.H{
		"type":      "allocations",
		"data":      stats.Allocations,
		"timestamp": stats.Timestamp,
	})
}

// GCHandler возвращает статистику GC
func (serv *Server) GCHandler(c *gin.Context) {
	stats := serv.analyzer.Collect()

	c.JSON(http.StatusOK, gin.H{
		"type":      "gc",
		"data":      stats.GC,
		"timestamp": stats.Timestamp,
	})
}

// MemoryHandler возвращает статистику памяти
func (serv *Server) MemoryHandler(c *gin.Context) {
	stats := serv.analyzer.Collect()

	c.JSON(http.StatusOK, gin.H{
		"type":      "memory",
		"data":      stats.Memory,
		"timestamp": stats.Timestamp,
	})
}

// SystemHandler возвращает полную статистику
func (serv *Server) SystemHandler(c *gin.Context) {
	stats := serv.analyzer.Collect()

	c.JSON(http.StatusOK, gin.H{
		"type":      "system",
		"data":      stats,
		"timestamp": stats.Timestamp,
	})
}

// PrometheusMetricsHandler возвращает метрики в формате Prometheus
func (serv *Server) PrometheusMetricsHandler(c *gin.Context) {
	stats := serv.analyzer.Collect()

	c.String(http.StatusOK,
		"# HELP go_allocations_total Total number of allocations\n"+
			"# TYPE go_allocations_total counter\n"+
			"go_allocations_total %d\n\n"+

			"# HELP go_frees_total Total number of frees\n"+
			"# TYPE go_frees_total counter\n"+
			"go_frees_total %d\n\n"+

			"# HELP go_live_objects Current number of live objects\n"+
			"# TYPE go_live_objects gauge\n"+
			"go_live_objects %d\n\n"+

			"# HELP go_gc_cycles_total Total number of GC cycles\n"+
			"# TYPE go_gc_cycles_total counter\n"+
			"go_gc_cycles_total %d\n\n"+

			"# HELP go_forced_gc_cycles_total Total number of forced GC cycles\n"+
			"# TYPE go_forced_gc_cycles_total counter\n"+
			"go_forced_gc_cycles_total %d\n\n"+

			"# HELP go_memory_alloc_bytes Current memory allocated\n"+
			"# TYPE go_memory_alloc_bytes gauge\n"+
			"go_memory_alloc_bytes %d\n\n"+

			"# HELP go_memory_heap_alloc_bytes Heap memory allocated\n"+
			"# TYPE go_memory_heap_alloc_bytes gauge\n"+
			"go_memory_heap_alloc_bytes %d\n\n"+

			"# HELP go_memory_sys_bytes Memory obtained from OS\n"+
			"# TYPE go_memory_sys_bytes gauge\n"+
			"go_memory_sys_bytes %d\n\n"+

			"# HELP go_last_gc_time_seconds Last GC timestamp\n"+
			"# TYPE go_last_gc_time_seconds gauge\n"+
			"go_last_gc_time_seconds %d\n\n"+

			"# HELP go_gc_cpu_fraction Fraction of CPU used by GC\n"+
			"# TYPE go_gc_cpu_fraction gauge\n"+
			"go_gc_cpu_fraction %f\n\n"+

			"# HELP go_goroutines Number of goroutines\n"+
			"# TYPE go_goroutines gauge\n"+
			"go_goroutines %d\n",

		stats.Allocations.Mallocs,
		stats.Allocations.Frees,
		stats.Allocations.LiveObjects,
		stats.GC.NumGC,
		stats.GC.NumForcedGC,
		stats.Memory.Alloc,
		stats.Memory.HeapAlloc,
		stats.Memory.Sys,
		stats.GC.LastGC.Unix(),
		stats.GC.GCCPUFraction,
		stats.Goroutines,
	)
}

// MemoryAllocHandler для тестирования аллокаций
func (serv *Server) MemoryAllocHandler(c *gin.Context) {
	sizeStr := c.Query("size")
	size := 1048576 // 1MB по умолчанию

	if sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 {
			size = s
		}
	}

	// Создаем аллокацию
	data := make([]byte, size)
	for i := range data {
		data[i] = byte(i % 256)
	}

	c.JSON(http.StatusOK, gin.H{
		"action": "memory_allocated",
		"bytes":  size,
		"mb":     float64(size) / 1024 / 1024,
	})
}
