package concurrency

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/pozedorum/WB_project_4/task2/internal/chunks"
	"github.com/pozedorum/WB_project_4/task2/internal/models"
	"github.com/pozedorum/WB_project_4/task2/internal/options"
)

type Master struct {
	workers      []*Worker
	taskChan     chan models.Task
	resultChan   chan models.Result
	totalTasks   int
	totalFiles   int
	resultMap    map[int]models.Result
	resultMutex  sync.RWMutex
	done         chan bool
	progressChan chan int
	taskCounter  int
	wg           sync.WaitGroup // Добавляем WaitGroup для отслеживания воркеров
}

const (
	standardChanSize = 100
)

func NewMaster(workersCount int, flags *options.FlagStruct) (*Master, error) {
	workers := make([]*Worker, 0, workersCount)
	taskChan := make(chan models.Task, standardChanSize)
	resultChan := make(chan models.Result, standardChanSize)
	resultMap := make(map[int]models.Result, standardChanSize)

	master := &Master{
		workers:      workers,
		taskChan:     taskChan,
		resultChan:   resultChan,
		done:         make(chan bool),
		progressChan: make(chan int, standardChanSize), // Увеличиваем буфер
		taskCounter:  0,
		resultMap:    resultMap,
	}

	// Создаем и запускаем воркеры
	for id := 0; id < workersCount; id++ {
		master.wg.Add(1) // Увеличиваем счетчик для каждого воркера
		newWorker := newWorker(id, &master.wg, taskChan, resultChan, flags)
		master.workers = append(master.workers, newWorker)

	}

	// Запускаем сборщик результатов в отдельной горутине
	go master.resultCollector()

	return master, nil
}

// ProcessFilesStreaming - потоковая обработка файлов
func (m *Master) ProcessFilesStreaming(files []*os.File, operation, pattern string) error {
	m.totalFiles = len(files)
	// Запускаем потоковое создание задач в отдельной горутине
	go m.createTasksStreaming(files, operation, pattern)

	// Ждем завершения сбора результатов
	<-m.done

	// log.Printf("Processing completed: %d tasks processed", m.totalTasks)
	return nil
}

// createTasksStreaming - потоково создает задачи и отправляет в канал
func (m *Master) createTasksStreaming(files []*os.File, operation, pattern string) {
	defer close(m.taskChan) // Гарантируем закрытие канала задач
	lastChunkID := 0
	for _, file := range files {
		// log.Printf("Splitting file: %s", file.Name())

		// Разбиваем файл на чанки
		fileChunks, newLastChunkID, err := chunks.SplitFiles([]*os.File{file}, lastChunkID)
		if err != nil {
			log.Printf("Error splitting file %s: %v", file.Name(), err)
			continue
		}
		lastChunkID = newLastChunkID
		// Отправляем чанки в канал задач
		for _, chunk := range fileChunks {
			// fmt.Println("chunk id ", chunk.ChunkID)
			// fmt.Println("chunk start offset ", chunk.StartOffset)
			// fmt.Println("chunk end offset ", chunk.EndOffset)
			task := models.Task{
				ID:        m.taskCounter,
				Operation: operation,
				Pattern:   pattern,
				Chunk:     chunk,
			}

			m.taskChan <- task
			m.taskCounter++
			m.totalTasks++

			// log.Printf("Task %d created for file %s", task.ID, file.Name())
		}
	}
	// log.Printf("All tasks created: %d total tasks", m.totalTasks)
}

// resultCollector собирает результаты из канала
func (m *Master) resultCollector() {
	go func() {
		m.wg.Wait()
		close(m.resultChan)
	}()

	receivedResults := 0
	for result := range m.resultChan {
		m.resultMutex.Lock()
		m.resultMap[result.ChunkID] = result
		m.resultMutex.Unlock()

		receivedResults++

		select {
		case m.progressChan <- receivedResults:
		default:
		}

		if receivedResults >= m.totalTasks {
			break
		}
	}

	close(m.progressChan)
	m.done <- true
}

// MergeResults объединяет все результаты
func (m *Master) MergeResults() []string {
	var allLines []string

	m.resultMutex.RLock()
	defer m.resultMutex.RUnlock()

	for chunkID := 0; chunkID < m.totalTasks; chunkID++ {
		if result, exists := m.resultMap[chunkID]; exists {
			if result.Error == nil {
				allLines = append(allLines, m.AddPath(result.Lines, result.FilePath)...)
			} else {
				allLines = append(allLines, fmt.Sprintf("error with file %s: %v", result.FilePath, result.Error))
			}
		} else {
			allLines = append(allLines, fmt.Sprintf("error: missing result for chunk %d", chunkID))
		}
	}
	return allLines
}

func (m *Master) AddPath(lines []string, filepath string) []string {
	if m.totalFiles == 1 {
		return lines
	}
	for ind, line := range lines {
		lines[ind] = fmt.Sprintf("%s:%s", filepath, line)
	}
	return lines
}
