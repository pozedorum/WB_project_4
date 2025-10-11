package concurrency

import (
	"fmt"
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
	resultMap    map[int]models.Result
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

	// Запускаем потоковое создание задач в отдельной горутине
	go m.createTasksStreaming(files, operation, pattern)

	// Ждем завершения сбора результатов
	<-m.done

	//log.Printf("Processing completed: %d tasks processed", m.totalTasks)
	return nil
}

// createTasksStreaming - потоково создает задачи и отправляет в канал
func (m *Master) createTasksStreaming(files []*os.File, operation, pattern string) {
	defer close(m.taskChan) // Гарантируем закрытие канала задач

	for _, file := range files {
		//log.Printf("Splitting file: %s", file.Name())

		// Разбиваем файл на чанки
		fileChunks, err := chunks.SplitFiles([]*os.File{file})
		if err != nil {
			//log.Printf("Error splitting file %s: %v", file.Name(), err)
			continue
		}

		// Отправляем чанки в канал задач
		for _, chunk := range fileChunks {
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

			//log.Printf("Task %d created for file %s", task.ID, file.Name())
		}
	}
	//log.Printf("All tasks created: %d total tasks", m.totalTasks)
}

// resultCollector собирает результаты из канала
func (m *Master) resultCollector() {
	//log.Println("Result collector started")

	// Ждем завершения всех воркеров в отдельной горутине
	go func() {
		m.wg.Wait() // Ждем завершения всех воркеров
		//fmt.Println("closing resultChan")
		close(m.resultChan) // Закрываем канал результатов
	}()

	for result := range m.resultChan {
		m.resultMap[result.ChunkID] = result
		// Неблокирующая отправка в канал прогресса
		select {
		case m.progressChan <- len(m.resultMap):
		default:
			// Пропускаем обновление прогресса если канал полный
		}

		//log.Printf("Result collected for task %d from worker %d (lines: %d)",
		//	result.TaskID, result.WorkerID, len(result.Lines))
	}
	// Все результаты собраны (канал закрыт)
	close(m.progressChan)
	m.done <- true
}

<<<<<<< Updated upstream
// progressDisplay отображает и обновляет полосу прогресса
func (m *Master) progressDisplay() {
	if m.totalTasks == 0 {
		fmt.Println("No tasks to process")
		return
	}

	lastPercent := -1
	for completed := range m.progressChan {
		percent := int(float64(completed) / float64(m.totalTasks) * 100)

		if percent != lastPercent {
			m.redrawProgressBar(completed, m.totalTasks, percent)
			lastPercent = percent
		}

		if completed >= m.totalTasks {
			m.redrawProgressBar(completed, m.totalTasks, 100)
			fmt.Println()
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

=======
>>>>>>> Stashed changes
// MergeResults объединяет все результаты
func (m *Master) MergeResults() []string {
	var allLines []string
	for i := 0; i < m.totalTasks; i++ {
		result := m.resultMap[i]
		if result.Error == nil {
			allLines = append(allLines, result.Lines...)
		} else {
			allLines = append(allLines, fmt.Sprintf("error with file %s: %v", result.FilePath, result.Error))
		}
		//fmt.Println(m.resultMap[i].ChunkID)
	}
	return allLines
}
