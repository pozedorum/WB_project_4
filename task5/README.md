# URL Shortener Service

Сервис для сокращения URL-ссылок с сбором аналитики переходов.

## Реализовано

- Создание коротких ссылок из длинных URL
- Перенаправление по коротким ссылкам
- Детальная аналитика переходов:
  - Общее количество кликов и уникальных посетителей
  - Статистика по дням и месяцам
  - Аналитика по браузерам, ОС и устройствам (распаршенный user-agent)
- Веб-интерфейс для управления и просмотра статистики
- Возможность создания кастомных ссылок (буквы и цифры 1-6 символов)

## Quick Start

### 1. Клонирование и запуск

```bash
git clone <repository-url>
cd url-shortener
make build
make run
```

Сервис будет доступен по адресу: http://localhost:8080

### 2. Веб-интерфейс

После запуска откройте в браузере: http://localhost:8080

Интерфейс включает:
- Создание коротких ссылок
- Тестирование редиректов
- Просмотр детальной аналитики

## Тестирование производительности и профилирование

### Бенчмарк-тесты

```bash
# Запуск всех бенчмарков
make bench

# Результаты должны показать улучшение производительности:
# BenchmarkRedirectHandler: ~12605 ns/op (улучшение на 10.4%)
# BenchmarkShortenHandler: ~13759 ns/op
```

### Нагрузочное тестирование

```bash
# Локальное нагрузочное тестирование
make load-test-docker

# Ожидаемые результаты:
# - URL Shortening: ~1200 RPS
# - URL Redirects: ~1145 RPS (улучшение на 17.8%)
# - Latency: ~40ms (улучшение на 14.4%)
```

### Профилирование производительности

#### Локальное профилирование
```bash
# Полный цикл локального профилирования
make profile-local

# Анализ результатов
make analyze-local
```

#### Docker профилирование
```bash
# CPU профилирование в Docker
make profile-docker-cpu

# Memory профилирование в Docker  
make profile-docker-mem

# Trace профилирование в Docker
make profile-docker-trace

# Анализ всех Docker профилей
make analyze-docker
```

### Анализ профилей вручную

```bash
# Интерактивный анализ CPU профиля
go tool pprof pprof/cpu_docker.pprof

# Анализ памяти
go tool pprof pprof/heap_after_load_docker.pprof

# Анализ трассировки
go tool trace pprof/trace_docker.out
```

## Результаты оптимизации

### До оптимизации
- **Redirect RPS**: 972.42
- **Redirect Latency**: 47.76 ms
- **BenchmarkRedirectHandler**: 14077 ns/op

### После оптимизации  
- **Redirect RPS**: 1145.50 (+17.8%)
- **Redirect Latency**: 40.88 ms (-14.4%)
- **BenchmarkRedirectHandler**: 12605 ns/op (+10.4%)

### Выполненные оптимизации

1. **SQL оптимизация** - убраны транзакции, объединены запросы
2. **Индексы БД** - добавлены конкурентные индексы
3. **Пул соединений** - увеличен до 50 подключений
4. **Пакетная обработка** - внедрен ClickBatcher для аналитики

## API Endpoints

### Создание короткой ссылки
```bash
POST /shorten
Content-Type: application/json

{
    "url": "https://example.com/very/long/url/path"
}
```

Ответ:
```json
{
    "short_url": "/s/abc123",
    "original_url": "https://example.com/very/long/url/path"
}
```

### Переход по короткой ссылке
```bash
GET /s/{short_code}
```

### Получение аналитики
```bash
GET /analytics/{short_code}?period=7d&groupBy=browser
```

Параметры:
- `period`: "1d", "7d", "30d" (период фильтрации)
- `groupBy`: "day", "month", "browser", "os", "device", "user-agent"

### Health check
```bash
GET /health
```

## Мониторинг и отладка

### Pprof эндпоинты
Сервис предоставляет pprof эндпоинты для отладки:
- `http://localhost:8080/debug/pprof/` - общий pprof
- `http://localhost:8080/debug/pprof/profile` - CPU профиль
- `http://localhost:8080/debug/pprof/heap` - Memory профиль
- `http://localhost:8080/debug/pprof/trace` - Trace

### Логи
```bash
# Просмотр логов приложения
docker compose logs app

# Логи в реальном времени
docker compose logs -f app
```

## Структура проекта

```
url-shortener/
├── benchmarks/           # Бенчмарк-тесты
├── scripts/             # Скрипты для профилирования
│   ├── profile_local.sh
│   ├── profile_docker_cpu.sh
│   ├── profile_docker_memory.sh
│   ├── profile_docker_trace.sh
│   └── test_load.go
├── pprof/               # Собранные профили
├── migrations/          # SQL миграции
├── internal/
│   ├── server/          # HTTP обработчики
│   ├── service/         # Бизнес-логика + ClickBatcher
│   ├── postgres/        # Репозиторий PostgreSQL
│   └── models/          # Модели данных
├── web/                 # Веб-интерфейс (HTML/JS)
├── docker-compose.yml
├── Dockerfile
└── Makefile
```

## Проверка оптимизаций

Для проверки результатов оптимизации:

1. **Запустите бенчмарки**: `make bench`
2. **Проведите нагрузочное тестирование**: `make load-test-docker`  
3. **Соберите профили**: `make profile-docker-cpu`
4. **Проанализируйте результаты**: `make analyze-docker`

Ожидаемые улучшения:
- Увеличение RPS на 15-20%
- Снижение задержек на 10-15%
- Стабильное потребление памяти

## Технологии

- **Backend**: Go 1.24
- **База данных**: PostgreSQL 15
- **HTTP фреймворк**: Gin
- **Профилирование**: pprof, go tool trace
- **Контейнеризация**: Docker + Docker Compose
- **Логирование**: Zerolog

## Разработка

### Полный цикл тестирования
```bash
make clean
make bench
make profile-docker
make analyze-docker
```

### Локальная разработка
```bash
make build
make up

# Или напрямую
go run ./cmd/main.go
```
