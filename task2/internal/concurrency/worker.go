package concurrency

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/pozedorum/WB_project_4/task2/internal/grep"
	"github.com/pozedorum/WB_project_4/task2/internal/models"
	"github.com/pozedorum/WB_project_4/task2/internal/options"
)

type Worker struct {
	id         int
	taskChan   <-chan models.Task
	resultChan chan<- models.Result
	flags      *options.FlagStruct
}

func newWorker(id int, taskChan <-chan models.Task, resultChan chan<- models.Result, flags *options.FlagStruct) *Worker {
	w := &Worker{
		id:         id,
		taskChan:   taskChan,
		resultChan: resultChan,
		flags:      flags,
	}
	go w.run()
	return w
}

func (w *Worker) run() {
	log.Printf("Worker %d started", w.id)

	for task := range w.taskChan { // Автоматически завершится при закрытии канала
		// Обрабатываем задачу
		result := w.processTask(task)
		w.resultChan <- result
	}

	log.Printf("Worker %d finished (task channel closed)", w.id)
}

// processTask обрабатывает одну задачу
func (w *Worker) processTask(task models.Task) models.Result {
	log.Printf("Worker %d processing task %d", w.id, task.ID)
	var reader io.Reader
	res := models.Result{
		TaskID:   task.ID,
		WorkerID: w.id,
		Lines:    nil,
		Error:    nil,
	}
	if task.Operation != "grep" {
		res.Error = fmt.Errorf("operation is not supported")
		return res
	}
	reader, res.Error = task.Chunk.GetChunkReader()
	if res.Error != nil {
		res.Error = fmt.Errorf("failed to get chunk reader: %v", res.Error)
		return res
	}
	defer func() {
		if closer, ok := reader.(io.Closer); ok {
			closer.Close()
		}
	}()
	res.Lines, res.Error = w.processChunkGrep(reader)

	return res
}

// processGrepChunk обрабатывает чанк с использованием пакета grep
func (w *Worker) processChunkGrep(reader io.Reader) ([]string, error) {
	var resultLines []string
	var outputBuffer strings.Builder

	// Вызываем функцию grep
	err := grep.Grep(reader, *w.flags, &outputBuffer)
	if err != nil {
		return nil, fmt.Errorf("grep error: %v", err)
	}

	// Разбиваем результат на строки
	if outputBuffer.Len() > 0 {
		resultLines = strings.Split(strings.TrimSuffix(outputBuffer.String(), "\n"), "\n")
	}

	return resultLines, nil
}
