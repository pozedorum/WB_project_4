// Package or предоставляет утилиту для объединения нескольких done-каналов в один.
//
// Основная функция Or позволяет отслеживать закрытие любого из переданных каналов
// через единый результирующий канал. Это полезно в конкурентных сценариях,
// когда необходимо ожидать завершения одной из нескольких операций.
//
// Пример использования:
//
//	sig := func(after time.Duration) <-chan any {
//	    c := make(chan any)
//	    go func() {
//	        defer close(c)
//	        time.Sleep(after)
//	    }()
//	    return c
//	}
//
//	start := time.Now()
//	<-or.Or(
//	    sig(2*time.Hour),
//	    sig(5*time.Minute),
//	    sig(1*time.Second),
//	)
//	fmt.Printf("done after %v", time.Since(start))
//	// Output: done after 1s
package or

// Or объединяет несколько done-каналов в один результирующий канал.
//
// Функция принимает переменное количество каналов типа <-chan any и возвращает
// единый канал <-chan any, который закрывается, когда закрывается ЛЮБОЙ из
// переданных входных каналов.
//
// Особенности:
//   - Если не передано ни одного канала, возвращается закрытый канал
//   - Если передан только один канал, то он и возвращается
//   - Для двух и более каналов создается новый канал, который отслеживает
//     закрытие любого из входных каналов
//
// Использование:
//
//	ch1 := make(chan any)
//	ch2 := make(chan any)
//	orChan := Or(ch1, ch2)
//
//	// В другой горутине
//	go func() {
//	    time.Sleep(time.Second)
//	    close(ch1)
//	}()
//
//	// Будет получено значение, когда закроется ch1
//	<-orChan
//	fmt.Println("One of the channels closed!")
//
// Параметры:
//   - channels: переменное количество каналов для отслеживания
//
// Возвращает:
//   - <-chan any: канал, который закрывается при закрытии любого из входных каналов
func Or(channels ...<-chan any) <-chan any {
	switch len(channels) {
	case 0:
		ch := make(chan any)
		close(ch)
		return ch
	case 1:
		return channels[0]
	}

	orDone := make(chan any)

	go func() {
		defer close(orDone)

		switch len(channels) {
		case 2:
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default:
			select {
			case <-channels[0]:
			case <-channels[1]:
			case <-channels[2]:
			case <-Or(append(channels[3:], orDone)...):
			}
		}
	}()

	return orDone
}
