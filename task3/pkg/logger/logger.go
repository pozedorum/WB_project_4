package logger

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Logger struct {
	fileLogger *log.Logger // Для записи в файл
	consoleLog bool        // Дублировать в консоль

}

func New(logFile string, consoleLog bool) (*Logger, error) {

	dir := filepath.Dir(logFile)
	if dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}
	}
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &Logger{
		fileLogger: log.New(file, "", log.LstdFlags),
		consoleLog: consoleLog,
	}, nil
}

// Middleware для HTTP-запросов
func (l *Logger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		l.log("[REQUEST]", r.Method, r.URL.Path, r.RemoteAddr, duration.String())

	})
}

func (l *Logger) LogError(r *http.Request, err error) {
	if err == nil {
		return
	}
	l.log("[ERROR]", r.Method, r.URL.Path, err.Error())

}

func (l *Logger) LogInfo(info string) {
	l.log("[INFO]", info)
}

func (l *Logger) log(prefix string, v ...string) {
	parts := make([]string, 0, len(v)+1)
	parts = append(parts, prefix)
	parts = append(parts, v...)

	msg := time.Now().Format("2006/01/02 15:04:05 ") + strings.Join(parts, " ")
	l.fileLogger.Println(msg)

	if l.consoleLog {
		log.Println(msg)
	}
}
