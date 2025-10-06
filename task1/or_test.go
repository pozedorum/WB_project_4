package or

import (
	"testing"
	"time"
)

func testFunc(ch chan any, sleepTime time.Duration) {
	time.Sleep(sleepTime)
	close(ch)
}

// TestOr тестирует функцию Or на различных сценариях:
// - Один канал
// - Несколько каналов с разным временем закрытия
// - Пустой список каналов
func TestOr(t *testing.T) {
	// Тест с одним каналом
	t.Run("single channel", func(t *testing.T) {
		ch := make(chan interface{})

		go testFunc(ch, 50*time.Millisecond)
		start := time.Now()

		<-Or(ch)
		duration := time.Since(start)

		if duration < 50*time.Millisecond || duration > 100*time.Millisecond {
			t.Errorf("expected ~50ms, got %v", duration)
		}
	})

	// Тест с несколькими каналами
	t.Run("multiple channels", func(t *testing.T) {
		ch1 := make(chan interface{})
		ch2 := make(chan interface{})
		ch3 := make(chan interface{})

		go testFunc(ch1, 100*time.Millisecond)
		go testFunc(ch2, 200*time.Millisecond)
		go testFunc(ch3, 30*time.Millisecond)

		start := time.Now()
		<-Or(ch1, ch2, ch3)
		duration := time.Since(start)

		// Должен завершиться когда закрывается самый быстрый канал (~30ms)
		if duration < 25*time.Millisecond || duration > 80*time.Millisecond {
			t.Errorf("expected ~30ms, got %v", duration)
		}
	})

	// Тест без каналов
	t.Run("no channels", func(t *testing.T) {
		start := time.Now()
		<-Or()
		duration := time.Since(start)

		// Без каналов должен сразу завершиться
		if duration > 5*time.Millisecond {
			t.Errorf("expected immediate return, got %v", duration)
		}
	})
}
