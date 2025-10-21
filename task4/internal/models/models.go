package models

import "time"

// SystemStats представляет полный набор метрик
type SystemStats struct {
	Memory      *MemStats
	GC          *GCStats
	Allocations *AllocationStats
	Goroutines  int
	Timestamp   time.Time
}

// MemStats представляет ключевые метрики памяти
type MemStats struct {
	// Alloc - текущие байты в использовании (живые объекты)
	Alloc uint64
	// TotalAlloc - всего байт выделено за время работы
	TotalAlloc uint64
	// Sys - всего байт получено от ОС (HeapAlloc + память под служебные данные)
	Sys uint64
	// HeapAlloc - байт в куче (основная область памяти)
	HeapAlloc uint64
	// HeapInuse - байт в используемой куче
	HeapInuse uint64
	// HeapObjects - количество объектов в куче
	HeapObjects uint64
}

// GCStats представляет ключевые метрики сборщика мусора
type GCStats struct {
	// NumGC - всего циклов сборки мусора
	NumGC uint32
	// NumForcedGC - принудительных сборок
	NumForcedGC uint32
	// LastGC - время последней сборки
	LastGC time.Time
	// LastGCPause - длительность последней паузы GC
	LastGCPause time.Duration
	// PauseTotal - общее время в паузах GC
	PauseTotal time.Duration
	// GCCPUFraction - доля CPU на GC
	GCCPUFraction float64
}

// AllocationStats представляет ключевые метрики аллокаций
type AllocationStats struct {
	// Mallocs - всего аллокаций (вызовы malloc)
	Mallocs uint64
	// Frees - всего освобождений (вызовы free)
	Frees uint64
	// LiveObjects - живые объекты (Mallocs - Frees)
	LiveObjects uint64
}
