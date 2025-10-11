package concurrency

import (
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/pozedorum/WB_project_4/task2/internal/grep"
	"github.com/pozedorum/WB_project_4/task2/internal/models"
	"github.com/pozedorum/WB_project_4/task2/internal/options"
)

type Worker struct {
	id         int
	taskChan   <-chan models.Task
	resultChan chan<- models.Result
	flags      *options.FlagStruct
	wg         *sync.WaitGroup
}

func newWorker(id int, wg *sync.WaitGroup, taskChan <-chan models.Task, resultChan chan<- models.Result, flags *options.FlagStruct) *Worker {
	w := &Worker{
		id:         id,
		taskChan:   taskChan,
		resultChan: resultChan,
		flags:      flags,
		wg:         wg,
	}
	go w.run()
	return w
}

func (w *Worker) run() {
	// log.Printf("Worker %d started", w.id)

	for task := range w.taskChan {
		result := w.processTask(task)
		// fmt.Println("result uploaded", result.ChunkID)
		w.resultChan <- result
	}

	// log.Printf("Worker %d finished (task channel closed)", w.id)

	// Уведомляем master о завершении работы
	w.wg.Done()
}

// processTask обрабатывает одну задачу
func (w *Worker) processTask(task models.Task) models.Result {
	// log.Printf("Worker %d processing task %d", w.id, task.ID)

	var reader io.Reader
	res := models.Result{
		TaskID:   task.ID,
		ChunkID:  task.Chunk.ChunkID,
		WorkerID: w.id,
		Lines:    nil,
		Error:    nil,
		FilePath: task.Chunk.FilePath, // Добавляем информацию о файле
	}

	if task.Operation != models.OperationGrep {
		res.Error = fmt.Errorf("operation is not supported")
		return res
	}
	// fmt.Println("worker offsets: ", task.Chunk.StartOffset, task.Chunk.EndOffset)
	reader, res.Error = task.Chunk.GetChunkReader()
	if res.Error != nil {
		res.Error = fmt.Errorf("failed to get chunk reader: %v", res.Error)
		return res
	}

	// Обрабатываем данные
	res.Lines, res.Error = w.processChunkGrep(reader)

	// Закрываем reader если он реализует интерфейс Closer
	if closer, ok := reader.(io.Closer); ok {
		if err := closer.Close(); err != nil {
			fmt.Printf("grep internal error: %v", err)
		}
	}
	// fmt.Println("worker end with: ", res.Lines)
	return res
}

// processChunkGrep обрабатывает чанк с использованием пакета grep
func (w *Worker) processChunkGrep(reader io.Reader) ([]string, error) {
	var outputBuffer strings.Builder

	// Вызываем функцию grep
	err := grep.Grep(reader, *w.flags, &outputBuffer)
	if err != nil {
		return nil, fmt.Errorf("grep error: %v", err)
	}

	if outputBuffer.Len() == 0 {
		return []string{}, nil
	}

	output := strings.TrimSuffix(outputBuffer.String(), "\n")
	if output == "" {
		return []string{}, nil
	}

	return strings.Split(output, "\n"), nil
}
