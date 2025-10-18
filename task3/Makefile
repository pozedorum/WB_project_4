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
	rm -f warehousecontrol

sql_connect:
	docker exec -it task6-postgres-1 psql -U postgres -d WarehouseControl

# Заполнение календаря различными событиями
create_many: create_event create_event2 create_event3 create_event4

# Тестовые запросы с ошибками
test_errors: test_empty_username test_empty_text test_empty_datetime test_invalid_event_id test_not_found test_past_date

# Успешные операции
create_event:
	@echo "\n=== Создание события 1 ==="
	curl -i -X POST http://localhost:8080/create_event \
  	-H "Content-Type: application/json" \
 	-d '{"username": "user1", "datetime": "2025-11-19T15:00:00Z", "text": "Тестовое событие", "title": "Встреча"}'

create_event2:
	@echo "\n=== Создание события 2 ==="
	curl -i -X POST http://localhost:8080/create_event \
  	-H "Content-Type: application/json" \
 	-d '{"username": "user1", "datetime": "2025-11-20T15:00:00Z", "text": "Тестовое событие через день", "title": "Совещание", "remind_before": 3600000000000}'

create_event3:
	@echo "\n=== Создание события 3 ==="
	curl -i -X POST http://localhost:8080/create_event \
  	-H "Content-Type: application/json" \
 	-d '{"username": "user1", "datetime": "2025-11-25T15:00:00Z", "text": "Тестовое событие через неделю", "title": "Презентация", "remind_before": 7200000000000}'

create_event4:
	@echo "\n=== Создание события 4 ==="
	curl -i -X POST http://localhost:8080/create_event \
  	-H "Content-Type: application/json" \
 	-d '{"username": "user2", "datetime": "2025-12-10T15:00:00Z", "text": "Тестовое событие через месяц", "title": "Отчет"}'

update_event:
	@echo "\n=== Обновление события ==="
	curl -i -X POST http://localhost:8080/update_event \
  	-H "Content-Type: application/json" \
  	-d '{"event_id": 1, "datetime": "2025-11-19T16:00:00Z", "text": "Обновленное событие", "title": "Обновленная встреча", "remind_before": 1800000000000}'

delete_event:
	@echo "\n=== Удаление события ==="
	curl -i -X POST http://localhost:8080/delete_event \
  	-H "Content-Type: application/json" \
  	-d '{"event_id": 2}'

get_day_events:
	@echo "\n=== События на день ==="
	curl -i "http://localhost:8080/events_for_day?username=user1&date=2025-11-19"

get_week_events:
	@echo "\n=== События на неделю ==="
	curl -i "http://localhost:8080/events_for_week?username=user1&date=2025-11-19"

get_month_events:
	@echo "\n=== События на месяц ==="
	curl -i "http://localhost:8080/events_for_month?username=user1&date=2025-11-01"

# Тесты ошибок согласно вашей новой структуре ошибок

# Пустой username
test_empty_username:
	@echo "\n=== Тест: Пустой username ==="
	curl -i -X POST http://localhost:8080/create_event \
		-H "Content-Type: application/json" \
		-d '{"username": "", "datetime": "2025-11-20T15:00:00Z", "text": "Событие с пустым username", "title": "Тест"}'

# Пустой текст (должен работать, т.к. text не required)
test_empty_text:
	@echo "\n=== Тест: Пустой текст ==="
	curl -i -X POST http://localhost:8080/create_event \
		-H "Content-Type: application/json" \
		-d '{"username": "user1", "datetime": "2025-11-20T15:00:00Z", "text": "", "title": "Тест"}'

# Пустая дата
test_empty_datetime:
	@echo "\n=== Тест: Пустая дата ==="
	curl -i -X POST http://localhost:8080/create_event \
		-H "Content-Type: application/json" \
		-d '{"username": "user1", "datetime": "", "text": "Событие без даты", "title": "Тест"}'

# Неверный формат event_id при удалении
test_invalid_event_id:
	@echo "\n=== Тест: Неверный формат event_id ==="
	curl -i -X POST http://localhost:8080/delete_event \
		-H "Content-Type: application/json" \
		-d '{"event_id": "invalid_id"}'

# Попытка обновления несуществующего события
test_not_found:
	@echo "\n=== Тест: Обновление несуществующего события ==="
	curl -i -X POST http://localhost:8080/update_event \
		-H "Content-Type: application/json" \
		-d '{"event_id": 999, "datetime": "2025-11-20T15:00:00Z", "text": "Несуществующее событие", "title": "Тест"}'

# Событие в прошлом
test_past_date:
	@echo "\n=== Тест: Событие в прошлом ==="
	curl -i -X POST http://localhost:8080/create_event \
		-H "Content-Type: application/json" \
		-d '{"username": "user1", "datetime": "2020-01-01T15:00:00Z", "text": "Событие в прошлом", "title": "Тест"}'

# Получение событий без указания username
test_missing_username:
	@echo "\n=== Тест: Отсутствует username в запросе событий ==="
	curl -i "http://localhost:8080/events_for_day?date=2025-11-19"

# Получение событий без указания даты
test_missing_date:
	@echo "\n=== Тест: Отсутствует дата в запросе событий ==="
	curl -i "http://localhost:8080/events_for_day?username=user1"

# Дополнительные тесты для разных пользователей
get_other_user_events:
	@echo "\n=== События другого пользователя ==="
	curl -i "http://localhost:8080/events_for_day?username=user2&date=2025-12-10"

# Тест с form-data (альтернативный формат)
test_form_data:
	@echo "\n=== Тест: Form-data формат ==="
	curl -i -X POST http://localhost:8080/create_event \
		-H "Content-Type: application/x-www-form-urlencoded" \
		-d "username=user1&datetime=2025-11-21T15:00:00Z&text=Событие через form-data&title=FormData Test"

# Комплексный тест - полный цикл
full_test:
	@echo "\n=== ПОЛНЫЙ ТЕСТОВЫЙ ЦИКЛ ==="
	@echo "1. Создание события..."
	curl -s -X POST http://localhost:8080/create_event \
  	-H "Content-Type: application/json" \
 	-d '{"username": "test_user", "datetime": "2025-11-22T15:00:00Z", "text": "Тест полного цикла", "title": "Полный цикл"}' | grep -E "(result|error)"
	
	@echo "2. Получение событий на день..."
	curl -s "http://localhost:8080/events_for_day?username=test_user&date=2025-11-22" | grep -E "(result|error)"
	
	@echo "3. Обновление события..."
	curl -s -X POST http://localhost:8080/update_event \
  	-H "Content-Type: application/json" \
  	-d '{"event_id": 5, "datetime": "2025-11-22T16:00:00Z", "text": "Обновленное событие полного цикла", "title": "Обновленный цикл"}' | grep -E "(result|error)"
	
	@echo "4. Удаление события..."
	curl -s -X POST http://localhost:8080/delete_event \
  	-H "Content-Type: application/json" \
  	-d '{"event_id": 5}' | grep -E "(result|error)"

# Тесты с напоминаниями (remind_before в наносекундах)
test_reminders:
	@echo "\n=== Тест с напоминаниями ==="
	@echo "Напоминание за 1 час (3600000000000 наносекунд)"
	curl -i -X POST http://localhost:8080/create_event \
  	-H "Content-Type: application/json" \
 	-d '{"username": "user1", "datetime": "2025-11-23T15:00:00Z", "text": "С напоминанием за 1 час", "title": "С напоминанием", "remind_before": 3600000000000}'
	
	@echo "Напоминание за 30 минут (1800000000000 наносекунд)"
	curl -i -X POST http://localhost:8080/create_event \
  	-H "Content-Type: application/json" \
 	-d '{"username": "user1", "datetime": "2025-11-23T16:00:00Z", "text": "С напоминанием за 30 мин", "title": "Короткое напоминание", "remind_before": 1800000000000}'
# Запуск всех основных тестов
test_all: create_many get_day_events get_week_events get_month_events test_errors

.PHONY: create_many test_errors create_event create_event2 create_event3 create_event4 update_event delete_event get_day_events get_week_events get_month_events test_empty_username test_empty_text test_empty_datetime test_invalid_event_id test_not_found test_past_date test_missing_username test_missing_date get_other_user_events test_form_data full_test test_all