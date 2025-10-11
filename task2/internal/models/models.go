package models

import "github.com/pozedorum/WB_project_4/task2/internal/chunks"

type Task struct {
	ID        int
	FilePath  string
	Chunk     chunks.Chunk // для больших файлов
	Operation string       // "grep", "cut", "sort"
	Pattern   string
}

type Result struct {
	TaskID   int
	WorkerID int
	Lines    []string
	Error    error
	FilePath string // важно для сборки обратно
	ChunkID  int    // для сборки чанков
}

// ChunkMetadata - метаинформация для сборки результатов
type ChunkMetadata struct {
	ChunkID   int
	FilePath  string
	WorkerID  int
	LineCount int   // Количество обработанных строк
	ByteCount int64 // Количество обработанных байт
}

const (
	OperationGrep = "grep"
	// OperationCut  = "cut"
	// OperationSort = "sort"
)
