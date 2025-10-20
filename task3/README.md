# Сервис Календаря Событий

Микросервис для управления событиями с отложенными напоминаниями через Telegram.

## Основные возможности

- **CRUD операции** для событий (создание, обновление, удаление, получение)
- **Отложенные напоминания** через RabbitMQ DLX
- **Telegram уведомления** с точным временем доставки
- **Фильтрация событий** по дням, неделям и месяцам
- **Асинхронная архитектура** с фоновыми воркерами
- **Поддержка напоминаний** за N минут до события
- **Логирование** с структурным форматом

## Архитектура

- **HTTP Server** - Gin framework с REST API
- **PostgreSQL** - хранение событий
- **RabbitMQ** - отложенная доставка напоминаний через DLX
- **Telegram Bot** - отправка уведомлений
- **Docker** - контейнеризация всех сервисов

## Быстрый старт

### 1. Запуск сервиса

```bash
# Сборка и запуск
make build
make run

# Или напрямую
docker-compose up -d
```

Сервис будет доступен по адресу: http://localhost:8080

### 2. Настройка Telegram Bot

1. Создайте бота через @BotFather
2. Получите токен бота
3. Добавьте токен в `docker-compose.yml`:

```yaml
environment:
  - TELEGRAM_BOT_TOKEN=your_telegram_bot_token_here
```
4. Напишите сообщение своему боту, а затем запустите скрипт `bash ./scripts/GetTelegramChatID.sh `
5. Получите ваш Chat ID и вставьте его в :


### 3. Проверка здоровья

```bash
curl http://localhost:8080/health
```

## API Endpoints

### Создание события
```bash
POST /create_event
Content-Type: application/json

{
    "usertoken": "user123",
    "title": "Встреча с командой",
    "text": "Обсуждение нового проекта",
    "datetime": "2025-12-22T15:00:00Z",
    "telegram_id": 1105031510,
    "remind_before": 5
}
```

### Обновление события
```bash
POST /update_event
Content-Type: application/json

{
    "event_id": 1,
    "title": "Обновленная встреча",
    "text": "Перенесено на другой зал",
    "datetime": "2025-12-22T16:00:00Z",
    "telegram_id": 1105031510,
    "remind_before": 10
}
```

### Удаление события
```bash
POST /delete_event
Content-Type: application/json

{
    "event_id": 1
}
```

### Получение событий

```bash
# События на день
GET /events_for_day?usertoken=user123&date=2025-12-22

# События на неделю  
GET /events_for_week?usertoken=user123&date=2025-12-22

# События на месяц
GET /events_for_month?usertoken=user123&date=2025-12-01
```

## Напоминания (Reminders)

### Параметр `remind_before`
- **Тип**: целое число (минуты)
- **Описание**: За сколько минут до события отправить напоминание
- **Пример**: `5` = напоминание за 5 минут до события

### Примеры использования

```bash
# Напоминание за 5 минут до события
{
    "remind_before": 5
}

# Напоминание за 1 час до события  
{
    "remind_before": 60
}

# Без напоминания (только запись в календаре)
{
    "remind_before": 0
}
```

## Тестирование

### Комплексное тестирование
```bash
# Все основные тесты
bash ./scripts/full_cycle.sh

# Тест отложенных напоминаний
./scripts/testReminderWithDelay.sh


### Использование Makefile

```bash
# Создание тестовых событий
make create_many

# Тестирование операций чтения
make get_day_events
make get_week_events  
make get_month_events

# Тестирование ошибок
make test_errors

```

### Тест отложенных напоминаний

```bash
# Создание события с напоминанием через 1 минуту
bash ./scripts/testReminderWithDelay.sh
```

## Конфигурация

### Переменные окружения

```env
# Server
SERVER_PORT=8080

#Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=eventbooker
DB_SSLMODE=disable

# RabbitMQ
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
RABBITMQ_QUEUE=event-reminders

# Telegram
TELEGRAM_BOT_TOKEN=your_bot_token_here
```

### Docker сервисы

- **app**: Основное приложение (Go)
- **postgres**: База данных PostgreSQL
- **rabbitmq**: Очередь сообщений с DLX
- **redis**: Кэширование (опционально)

## Мониторинг

### Логи приложения
```bash
# Просмотр логов в реальном времени
make logs

# Все логи
docker-compose logs

# Логи конкретного сервиса
docker-compose logs app
docker-compose logs rabbitmq
docker-compose logs postgres
```

### Статус очередей RabbitMQ
```bash
docker exec -it rabbitmq-1 rabbitmqctl list_queues name messages_ready messages_unacknowledged
```


## Workflow напоминаний

1. **Создание события** → Сохранение в PostgreSQL
2. **Планирование напоминания** → Отправка в RabbitMQ DLX очередь
3. **Ожидание** → Сообщение ждет в DLX до истечения TTL
4. **Доставка** → RabbitMQ перемещает сообщение в основную очередь
5. **Обработка** → Воркер отправляет уведомление в Telegram
6. **Уведомление** → Пользователь получает сообщение

### Частые проблемы

**Telegram бот не отправляет сообщения:**
- Проверьте корректность токена бота
- Убедитесь, что бот добавлен в чат
- Проверьте Chat ID получателя

**Напоминания не работают:**
- Проверьте статус RabbitMQ: `docker-compose logs rabbitmq`
- Убедитесь, что очередь создана: `docker exec -it rabbitmq-1 rabbitmqctl list_queues`
- Проверьте логи приложения: `make logs`

**Ошибки базы данных:**
- Проверьте подключение к PostgreSQL
- Убедитесь, что миграции выполнены

### Режим отладки

Для изменения количества логов меняйте переменные в пакете pkg/logger.go
По умолчанию стоят:
var (
	LevelInfo  = true
	LevelDebug = false
	LevelWarn  = true
	LevelError = true
)
