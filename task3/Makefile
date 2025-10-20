.PHONY: all build run lint rebuild stop logs clean create_many test_errors test_all full_test test_reminder test_reminders_multiple

all: run

build:
	docker compose build

run:
	docker compose up

lint:
	golangci-lint run ./...
	go vet ./...

rebuild:
	docker compose down -v
	docker compose up --build

stop:
	docker compose down

logs:
	docker compose logs -f app

clean:
	docker compose down -v
	rm -f event-calendar

# Тестовые запросы с использованием актуальных моделей
create_event:
	@echo "\n=== Создание события 1 ==="
	curl -s -X POST http://localhost:8080/create_event \
  	-H "Content-Type: application/json" \
 	-d '{"usertoken": "user1", "datetime": "2025-11-19T15:00:00Z", "text": "Тестовое событие", "title": "Встреча"}' | jq .

create_event2:
	@echo "\n=== Создание события 2 (с напоминанием) ==="
	curl -s -X POST http://localhost:8080/create_event \
  	-H "Content-Type: application/json" \
 	-d '{"usertoken": "user1", "datetime": "2025-11-20T15:00:00Z", "text": "Событие с напоминанием", "title": "Совещание", "remind_before": 5, "telegram_id": 1105031510}' | jq .

create_event3:
	@echo "\n=== Создание события 3 (с напоминанием за 10 мин) ==="
	curl -s -X POST http://localhost:8080/create_event \
  	-H "Content-Type: application/json" \
 	-d '{"usertoken": "user1", "datetime": "2025-11-25T15:00:00Z", "text": "Событие с напоминанием за 10 минут", "title": "Презентация", "remind_before": 10, "telegram_id": 1105031510}' | jq .

create_event4:
	@echo "\n=== Создание события для другого пользователя ==="
	curl -s -X POST http://localhost:8080/create_event \
  	-H "Content-Type: application/json" \
 	-d '{"usertoken": "user2", "datetime": "2025-12-10T15:00:00Z", "text": "Событие другого пользователя", "title": "Отчет"}' | jq .

# Заполнение календаря различными событиями
create_many: create_event create_event2 create_event3 create_event4

# Операции чтения
get_day_events:
	@echo "\n=== События на день ==="
	curl -s "http://localhost:8080/events_for_day?usertoken=user1&date=2025-11-19" | jq .

get_week_events:
	@echo "\n=== События на неделю ==="
	curl -s "http://localhost:8080/events_for_week?usertoken=user1&date=2025-11-19" | jq .

get_month_events:
	@echo "\n=== События на месяц ==="
	curl -s "http://localhost:8080/events_for_month?usertoken=user1&date=2025-11-01" | jq .

# Операции обновления и удаления
update_event:
	@echo "\n=== Обновление события ==="
	curl -s -X POST http://localhost:8080/update_event \
  	-H "Content-Type: application/json" \
  	-d '{"event_id": 1, "datetime": "2025-11-19T16:00:00Z", "text": "Обновленное событие", "title": "Обновленная встреча", "remind_before": 15, "telegram_id": 1105031510}' | jq .

delete_event:
	@echo "\n=== Удаление события ==="
	curl -s -X POST http://localhost:8080/delete_event \
  	-H "Content-Type: application/json" \
  	-d '{"event_id": 2}' | jq .

# Тесты ошибок
test_empty_usertoken:
	@echo "\n=== Тест: Пустой usertoken ==="
	curl -s -X POST http://localhost:8080/create_event \
		-H "Content-Type: application/json" \
		-d '{"usertoken": "", "datetime": "2025-11-20T15:00:00Z", "text": "Событие с пустым usertoken", "title": "Тест"}' | jq .

test_empty_datetime:
	@echo "\n=== Тест: Пустая дата ==="
	curl -s -X POST http://localhost:8080/create_event \
		-H "Content-Type: application/json" \
		-d '{"usertoken": "user1", "datetime": "", "text": "Событие без даты", "title": "Тест"}' | jq .

test_not_found:
	@echo "\n=== Тест: Обновление несуществующего события ==="
	curl -s -X POST http://localhost:8080/update_event \
		-H "Content-Type: application/json" \
		-d '{"event_id": 999, "datetime": "2025-11-20T15:00:00Z", "text": "Несуществующее событие", "title": "Тест"}' | jq .

test_past_date:
	@echo "\n=== Тест: Событие в прошлом ==="
	curl -s -X POST http://localhost:8080/create_event \
		-H "Content-Type: application/json" \
		-d '{"usertoken": "user1", "datetime": "2020-01-01T15:00:00Z", "text": "Событие в прошлом", "title": "Тест"}' | jq .

# Тесты RabbitMQ напоминаний
test_reminder:
	@echo "\n=== Тест отложенного напоминания ==="
	./testReminderWithDelay.sh

test_reminders_multiple:
	@echo "\n=== Тест нескольких напоминаний ==="
	./json_remind_before.sh

# Комплексные тесты
test_errors: test_empty_usertoken test_empty_datetime test_not_found test_past_date
