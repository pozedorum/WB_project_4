# ТЗ: Расширение функциональности сервиса календаря

## Текущее состояние
Реализован HTTP-сервер для управления событиями календаря с CRUD операциями, middleware логирования и хранением данных в памяти.

## Новые требования

### 1. Миграция на PostgreSQL
- Заменить хранение в памяти на PostgreSQL
- Создать структуру базы данных:
  - Таблица `events` с полями:
    - `id` (UUID, primary key)
    - `user_id` (integer)
    - `title` (text)
    - `description` (text)
    - `date` (date)
    - `time` (time, nullable)
    - `reminder_time` (timestamp, nullable) - время напоминания
    - `is_archived` (boolean, default false)
    - `created_at` (timestamp)
    - `updated_at` (timestamp)

### 2. Система напоминаний через Kafka
- Реализовать фоновый воркер для обработки напоминаний
- При создании/обновлении события с `reminder_time`:
  - Создавать сообщение в Kafka с данными о напоминании
  - Воркер отслеживает время и отправляет уведомления

### 3. Архивация старых событий
- Реализовать горутину для периодической архивации:
  - Переносить события старше 30 дней в архив (`is_archived = true`)
  - Интервал очистки: каждые 10 минут (конфигурируемо)

### 4. Асинхронная система логирования
- HTTP-обработчики не должны писать логи напрямую
- Реализовать канал для логов
- Отдельная горутина для обработки и записи логов

### 5. Конфигурация
- Поддержка конфигурации через environment variables:
  - `PORT` - порт сервера
  - `DB_URL` - строка подключения к PostgreSQL
  - `KAFKA_BROKERS` - адреса брокеров Kafka
  - `ARCHIVE_INTERVAL` - интервал архивации (минуты)
  - `LOG_LEVEL` - уровень логирования

## API Endpoints

### Существующие endpoints (обновлены):
```
POST    /events          - создание события
PUT     /events/{id}     - обновление события  
DELETE  /events/{id}     - удаление события
GET     /events/day      - события на день
GET     /events/week     - события на неделю
GET     /events/month    - события на месяц
```

### Новые endpoints:
```
GET     /events/archived - получить архивные события
POST    /events/{id}/remind - установить напоминание
```

## Форматы данных

### Создание события (POST /events):
```json
{
  "user_id": 123,
  "title": "Meeting",
  "description": "Team meeting",
  "date": "2024-01-15",
  "time": "14:30:00",
  "reminder_time": "2024-01-15T14:00:00Z"
}
```

### Напоминание (Kafka message):
```json
{
  "event_id": "uuid",
  "user_id": 123,
  "title": "Meeting",
  "reminder_time": "2024-01-15T14:00:00Z",
  "notification_type": "reminder"
}
```

## Структура проекта
```
calendar-service/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   ├── handlers/
│   ├── repository/
│   │   └── postgres/
│   ├── service/
│   ├── worker/
│   │   ├── reminder.go
│   │   ├── archiver.go
│   │   └── logger.go
│   └── kafka/
├── migrations/
├── docker-compose.yml
└── README.md
```

## Требования к реализации

### Бизнес-логика:
1. **Service Layer** - отделить бизнес-логику от HTTP-обработчиков
2. **Repository Pattern** - абстракция для работы с БД
3. **Dependency Injection** - для тестируемости

### Обработка ошибок:
- Глобальный обработчик ошибок
- Соответствующие HTTP статусы
- Детализация ошибок в логах

### Тестирование:
- Unit-тесты для service layer
- Интеграционные тесты с testcontainers
- Mock-репозитории для изоляции

### Безопасность:
- Валидация входных данных
- SQL injection protection
- Обработка крайних случаев

## Запуск и развертывание

### Требования:
- Go 1.21+
- PostgreSQL 14+
- Kafka 3.0+
- Docker & Docker Compose

### Локальный запуск:
```bash
# Запуск инфраструктуры
docker-compose up -d

# Применение миграций
go run cmd/migrate/main.go

# Запуск сервера
go run cmd/server/main.go
```

## Мониторинг и логи
- Structured logging с JSON форматом
- Метрики для мониторинга производительности
- Health checks для БД и Kafka

Это ТЗ обеспечивает масштабируемость, поддерживаемость и тестируемость кода while adding the required new functionality.