# Отчёт по производительности URL Shortener (Исходная версия)

## 📊 Текущие показатели производительности

### Результаты бенчмарков

```bash
make bench
```

| Бенчмарк | Результат | Память | Аллокации |
|----------|-----------|--------|-----------|
| **BenchmarkShortenHandler** | ~15000 ns/op | ~8000 B/op | ~40 allocs/op |
| **BenchmarkRedirectHandler** | ~14000 ns/op | ~9000 B/op | ~30 allocs/op |
| **BenchmarkCreateShortURL** | ~3200 ns/op | ~500 B/op | ~10 allocs/op |
| **BenchmarkGenerateShortURL** | ~600 ns/op | ~150 B/op | ~5 allocs/op |

### Нагрузочное тестирование

```bash
# Запуск нагрузочного теста редиректов
make load-test-redirect
```

**Ожидаемые показатели:**
- **RPS (Requests Per Second)**: ~970 запросов/сек
- **Latency**: ~48 ms средняя задержка
- **P95 Latency**: ~65 ms

### Профилирование производительности

#### Профилирование CPU
```bash
make profile-docker-cpu
make analyze-docker
```

**Ключевые функции (ожидаемые):**
- `GetOriginalURLIfExists`: ~29% CPU времени
- `RegisterClick`: ~21% CPU времени  
- `Redirect.func1`: ~22% CPU времени
- `Service.Redirect`: ~21% CPU времени

#### Профилирование памяти
```bash
make profile-docker-mem
make analyze-docker
```

**Потребление памяти (ожидаемое):**
- `GetOriginalURLIfExists`: ~134 kB
- `Redirect.func1`: ~27 kB
- `RegisterClick`: ~21 kB

## 🛠 Команды для сбора метрик

### Базовое профилирование
```bash
# Полный цикл профилирования
make profile-docker
make analyze-docker

# Или локальное профилирование
make profile-local
make analyze-local
```

### Специфичные профили
```bash
# Только CPU профиль
make profile-docker-cpu

# Только memory профиль  
make profile-docker-mem

# Только trace
make profile-docker-trace
```

### Бенчмарки
```bash
# Запуск всех бенчмарков
make bench

# Детальный бенчмарк с флагами
go test -bench=. -benchmem -benchtime=5s ./benchmarks/...
```

### Нагрузочное тестирование
```bash
# Тест редиректов
make load-test-redirect

# Тест сокращения URL
make load-test-shorten
```

## 📈 Метрики для сравнения

Для последующего сравнения с оптимизированной версией сохраните:

1. **Бенчмарки** - наносекунды на операцию, потребление памяти
2. **Load Test Results** - RPS и latency показатели  
3. **CPU Profile** - распределение времени выполнения по функциям
4. **Memory Profile** - потребление памяти по функциям
5. **Trace Analysis** - анализ параллелизма и блокировок

## 📝 Примечания

- Все тесты выполняются на одинаковом hardware
- База данных предварительно заполняется тестовыми данными
- Нагрузочные тесты выполняются с одинаковыми параметрами
- Профилирование проводится при стабильной нагрузке ~80% CPU

После сбора этих метрик можно приступать к оптимизациям и сравнивать результаты.