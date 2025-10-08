package concurrency

type Task struct {
	ID        string
	FilePath  string
	Chunk     Chunk  // для больших файлов
	Operation string // "grep", "cut", "sort"
	Pattern   string
}

type Result struct {
	TaskID   string
	WorkerID string
	Lines    []string
	Error    error
	FilePath string // важно для сборки обратно
	ChunkID  int    // для сборки чанков
}

// Chunk - описывает часть файла для обработки
type Chunk struct {
	FilePath    string // Путь к исходному файлу
	StartOffset int64  // Начальное смещение в байтах
	EndOffset   int64  // Конечное смещение в байтах
	ChunkID     int    // Уникальный ID чанка
	TotalChunks int    // Общее количество чанков
	FileSize    int64  // Размер всего файла (для валидации)

	// Для текстовых файлов - границы по строкам
	StartLineOffset int64 // Смещение до начала первой полной строки
	EndLineOffset   int64 // Смещение до конца последней полной строки
}

// ChunkMetadata - метаинформация для сборки результатов
type ChunkMetadata struct {
	ChunkID   int
	FilePath  string
	WorkerID  string
	LineCount int   // Количество обработанных строк
	ByteCount int64 // Количество обработанных байт
}
