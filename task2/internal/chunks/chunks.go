package chunks

import (
	"bufio"
	"io"
	"os"
)

type Chunk struct {
	FilePath    string // Путь к исходному файлу
	StartOffset int64  // Начальное смещение в байтах
	EndOffset   int64  // Конечное смещение в байтах
	ChunkID     int    // Уникальный ID чанка
	TotalChunks int    // Общее количество чанков
	FileSize    int64  // Размер всего файла (для валидации)
}

const (
	MaxChunkSize = 10 * 1024 * 1024 // 10MB
)

func SplitFiles(files []*os.File) ([]Chunk, error) {
	result := make([]Chunk, 0, len(files))
	lastChunkID := 0

	for _, file := range files {
		fileInfo, err := file.Stat()
		if err != nil {
			return nil, err
		}

		fileSize := fileInfo.Size()

		if fileSize > MaxChunkSize {
			// Большой файл - разбиваем на части по MaxChunkSize
			fileChunks, chunksCount, err := SplitBigFile(file, lastChunkID, fileSize)
			if err != nil {
				return nil, err
			}
			lastChunkID += chunksCount
			result = append(result, fileChunks...)
		} else {
			// Маленький файл - один чанк
			fileChunk := MakeChunkFromFile(file, lastChunkID, fileSize)
			lastChunkID++
			result = append(result, fileChunk)
		}
	}

	return result, nil
}

func MakeChunkFromFile(file *os.File, chunkID int, fileSize int64) Chunk {
	return Chunk{
		FilePath:    file.Name(),
		StartOffset: 0,
		EndOffset:   fileSize,
		ChunkID:     chunkID,
		TotalChunks: 1,
		FileSize:    fileSize,
	}
}

func SplitBigFile(file *os.File, startChunkID int, fileSize int64) ([]Chunk, int, error) {
	numChunks := int(fileSize / MaxChunkSize)
	if fileSize%MaxChunkSize != 0 {
		numChunks++
	}

	chunks := make([]Chunk, 0, numChunks)
	currentOffset := int64(0)

	for i := 0; i < numChunks; i++ {
		startOffset := currentOffset
		endOffset := startOffset + MaxChunkSize

		if endOffset > fileSize {
			endOffset = fileSize
		}

		// Для всех чанков кроме первого корректируем начало до границы строки
		if i > 0 {
			adjustedStart := adjustToLineStart(file, startOffset)
			// Если корректировка сдвинула начало за пределы конца файла, пропускаем чанк
			if adjustedStart >= fileSize {
				break
			}
			startOffset = adjustedStart
		}

		// Для всех чанков кроме последнего корректируем конец до границы строки
		if i < numChunks-1 && endOffset < fileSize {
			adjustedEnd := adjustToLineEnd(file, endOffset, fileSize)
			// Если корректировка вернула ту же позицию, увеличиваем немного
			if adjustedEnd <= endOffset {
				adjustedEnd = adjustToLineEnd(file, endOffset+1, fileSize)
			}
			endOffset = adjustedEnd
		}

		// Если чанк пустой, пропускаем
		if startOffset >= endOffset {
			break
		}

		chunk := Chunk{
			FilePath:    file.Name(),
			StartOffset: startOffset,
			EndOffset:   endOffset,
			ChunkID:     startChunkID + i,
			TotalChunks: numChunks,
			FileSize:    fileSize,
		}

		chunks = append(chunks, chunk)
		currentOffset = endOffset // Следующий чанк начинается с конца текущего

		// fmt.Printf("Created chunk %d: %d-%d (size: %d)\n",
		// 	chunk.ChunkID, startOffset, endOffset, endOffset-startOffset)
	}

	return chunks, len(chunks), nil
}

// adjustToLineStart - находит начало первой полной строки
func adjustToLineStart(file *os.File, offset int64) int64 {
	if offset == 0 {
		return 0
	}

	// Сохраняем текущую позицию файла
	originalPos, _ := file.Seek(0, io.SeekCurrent)
	defer file.Seek(originalPos, io.SeekStart) // Восстанавливаем позицию

	file.Seek(offset, io.SeekStart)
	reader := bufio.NewReader(file)

	// Читаем до следующей новой строки
	_, err := reader.ReadBytes('\n')
	if err == io.EOF {
		return offset // Достигли конца файла
	}

	// Новая позиция - начало следующей строки
	newPos, _ := file.Seek(0, io.SeekCurrent)
	return newPos
}

// adjustToLineEnd - находит конец последней полной строки
func adjustToLineEnd(file *os.File, offset int64, fileSize int64) int64 {
	if offset >= fileSize {
		return fileSize
	}

	// Сохраняем текущую позицию файла
	originalPos, _ := file.Seek(0, io.SeekCurrent)
	defer file.Seek(originalPos, io.SeekStart) // Восстанавливаем позицию

	file.Seek(offset, io.SeekStart)
	reader := bufio.NewReader(file)

	// Ищем конец текущей строки
	line, err := reader.ReadBytes('\n')
	if err == io.EOF {
		return fileSize // Достигли конца файла
	}

	// Новая позиция - после символа новой строки
	newPos := offset + int64(len(line))
	return newPos
}

// GetChunkReader - создает reader для чтения чанка
func (c *Chunk) GetChunkReader() (io.Reader, error) {
	file, err := os.Open(c.FilePath)
	if err != nil {
		return nil, err
	}

	// Перемещаемся к началу чанка
	file.Seek(c.StartOffset, io.SeekStart)

	// Ограничиваем чтение размером чанка
	return io.LimitReader(file, c.EndOffset-c.StartOffset), nil
}

// GetChunkSize - возвращает размер чанка в байтах
func (c *Chunk) GetChunkSize() int64 {
	return c.EndOffset - c.StartOffset
}
