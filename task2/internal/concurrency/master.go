package concurrency

import (
	"fmt"
	"log"
	"os"

	"github.com/pozedorum/WB_project_4/task2/internal/chunks"
	"github.com/pozedorum/WB_project_4/task2/internal/models"
	"github.com/pozedorum/WB_project_4/task2/internal/options"
)

type Master struct {
	workers      []*Worker
	taskChan     chan models.Task
	resultChan   chan models.Result
	quorum       int
	totalTasks   int
	resultList   []*models.Result
	done         chan bool
	progressChan chan int
	taskCounter  int // Счётчик задач для ID
}

const (
	standardChanSize = 100
)

func NewMaster(workersCount int, flags *options.FlagStruct) (*Master, error) {
	workers := make([]*Worker, 0, workersCount)
	taskChan := make(chan models.Task, standardChanSize)
	resultChan := make(chan models.Result, standardChanSize)

	for id := 0; id < workersCount; id++ {
		newWorker := newWorker(id, taskChan, resultChan, flags)
		workers = append(workers, newWorker)
	}

	master := &Master{
		workers:      workers,
		taskChan:     taskChan,
		resultChan:   resultChan,
		done:         make(chan bool),
		progressChan: make(chan int, 10),
		taskCounter:  0,
	}

	// Запускаем сборщик результатов в отдельной горутине
	go master.resultCollector()

	return master, nil
}

// ProcessFilesStreaming - потоковая обработка файлов
func (m *Master) ProcessFilesStreaming(files []*os.File, operation, pattern string) error {
	// Запускаем отображение прогресса
	go m.progressDisplay()

	// Запускаем потоковое создание задач в отдельной горутине
	go m.createTasksStreaming(files, operation, pattern)

	// Ждем завершения сбора результатов
	<-m.done

	log.Printf("Processing completed: %d tasks processed", m.totalTasks)
	return nil
}

// createTasksStreaming - потоково создает задачи и отправляет в канал
func (m *Master) createTasksStreaming(files []*os.File, operation, pattern string) {
	defer close(m.taskChan) // Закрываем канал когда все задачи отправлены

	for _, file := range files {
		log.Printf("Splitting file: %s", file.Name())

		// Разбиваем файл на чанки
		fileChunks, err := chunks.SplitFiles([]*os.File{file})
		if err != nil {
			log.Printf("Error splitting file %s: %v", file.Name(), err)
			continue
		}

		// Отправляем чанки в канал задач
		for _, chunk := range fileChunks {
			task := models.Task{
				ID:        m.taskCounter,
				Operation: operation,
				Pattern:   pattern,
				Chunk:     chunk,
			}

			// Отправляем задачу в канал (блокируется если канал полный)
			m.taskChan <- task
			m.taskCounter++
			m.totalTasks++

			log.Printf("Task %d created for file %s", task.ID, file.Name())
		}
	}

	log.Printf("All tasks created: %d total tasks", m.totalTasks)
}

// resultCollector собирает результаты из канала
func (m *Master) resultCollector() {
	log.Println("Result collector started")
	defer close(m.progressChan) // Важно закрыть канал прогресса

	for result := range m.resultChan {
		m.resultList = append(m.resultList, &result)
		m.progressChan <- len(m.resultList)

		log.Printf("Result collected for task %d from worker %d (lines: %d)",
			result.TaskID, result.WorkerID, len(result.Lines))
	}

	// Все результаты собраны (канал закрыт)
	m.done <- true
}

// progressDisplay отображает и обновляет полосу прогресса
func (m *Master) progressDisplay() {
	lastPercent := -1

	for completed := range m.progressChan {
		if m.totalTasks == 0 {
			continue
		}

		percent := int(float64(completed) / float64(m.totalTasks) * 100)

		// Обновляем только если процент изменился
		if percent != lastPercent {
			m.redrawProgressBar(completed, m.totalTasks, percent)
			lastPercent = percent
		}

		// Если все задачи завершены, выходим
		if completed >= m.totalTasks {
			m.redrawProgressBar(completed, m.totalTasks, 100)
			fmt.Println() // Переход на новую строку после завершения
			return
		}
	}
}

// redrawProgressBar перерисовывает полосу прогресса
func (m *Master) redrawProgressBar(completed, total, percent int) {
	// Очищаем строку и возвращаем каретку в начало
	fmt.Print("\r")

	// Рисуем прогресс-бар
	barWidth := 50
	filled := barWidth * percent / 100
	empty := barWidth - filled

	fmt.Printf("Progress: [")
	for i := 0; i < filled; i++ {
		fmt.Print("=")
	}
	for i := 0; i < empty; i++ {
		fmt.Print(" ")
	}
	fmt.Printf("] %d/%d (%d%%)", completed, total, percent)
}

// GetResults возвращает результаты
func (m *Master) GetResults() []*models.Result {
	return m.resultList
}

// MergeResults объединяет все результаты
func (m *Master) MergeResults() []string {
	var allLines []string
	for _, result := range m.resultList {
		if result.Error == nil {
			allLines = append(allLines, result.Lines...)
		} else {
			allLines = append(allLines, fmt.Sprintf("error with file %s: %v\n", result.FilePath, result.Error))
		}
	}
	return allLines
}

// Close останавливает мастер
func (m *Master) Close() {
	close(m.resultChan)
}
