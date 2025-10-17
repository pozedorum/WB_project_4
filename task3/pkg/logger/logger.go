package logger

import (
	"io"
	"os"
	"sync"

	"github.com/rs/zerolog"
)

/*
На случай использвоания логгера в других проектах оставлю здесь интерфейс
type Logger interface {
	Debug(operation, message string, keyvals ...interface{})
	Info(operation, message string, keyvals ...interface{})
	Warn(operation, message string, keyvals ...interface{})
	Error(operation, message string, keyvals ...interface{})
	Shutdown()
}
*/

type Logger struct {
	logChan  chan func()
	done     chan struct{}
	wg       sync.WaitGroup
	zerolog  zerolog.Logger
	isClosed bool
}

func NewLogger(serviceName, logFilePath string) (*Logger, error) {
	if err := os.MkdirAll("./logs", 0o755); err != nil {
		return nil, err
	}

	var output io.Writer
	if logFilePath != "" {
		file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return nil, err
		}
		output = file
	} else {
		output = os.Stdout
	}

	zerologLogger := zerolog.New(output).
		With().
		Timestamp().
		Str("service", serviceName).
		Logger()

	logger := &Logger{
		logChan:  make(chan func(), 1000),
		done:     make(chan struct{}),
		zerolog:  zerologLogger,
		isClosed: false,
	}

	logger.wg.Add(1)
	go logger.processLogs()

	return logger, nil
}

func (l *Logger) processLogs() {
	defer l.wg.Done()

	for {
		select {
		case logFunc := <-l.logChan:
			logFunc()
		case <-l.done:
			// Дописываем оставшиеся логи
			for {
				select {
				case logFunc := <-l.logChan:
					logFunc()
				default:
					return
				}
			}
		}
	}
}

// Методы для логирования
func (l *Logger) Debug(operation, message string, keyvals ...interface{}) {
	l.log(l.zerolog.Debug(), operation, message, keyvals...)
}

func (l *Logger) Info(operation, message string, keyvals ...interface{}) {
	l.log(l.zerolog.Info(), operation, message, keyvals...)
}

func (l *Logger) Warn(operation, message string, keyvals ...interface{}) {
	l.log(l.zerolog.Warn(), operation, message, keyvals...)
}

func (l *Logger) Error(operation, message string, keyvals ...interface{}) {
	l.log(l.zerolog.Error(), operation, message, keyvals...)
}

func (l *Logger) log(event *zerolog.Event, operation, message string, keyvals ...interface{}) {
	if l.isClosed {
		return
	}

	// Создаем замыкание с уже подготовленными данными
	logFunc := func() {
		event.Str("operation", operation)

		// Обрабатываем key-value пары
		for i := 0; i < len(keyvals); i += 2 {
			if i+1 < len(keyvals) {
				key, ok := keyvals[i].(string)
				if !ok {
					continue
				}
				event.Interface(key, keyvals[i+1])
			}
		}

		event.Msg(message)
	}

	// Асинхронная отправка
	select {
	case l.logChan <- logFunc:
	default:
		// Fallback синхронное логирование
		logFunc()
	}
}

func (l *Logger) Shutdown() {
	if l.isClosed {
		return
	}

	l.isClosed = true
	close(l.done)
	l.wg.Wait()
}

// var _ interfaces.Logger = (*Logger)(nil)
