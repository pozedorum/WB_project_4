all: run

run: clean_log
	go run cmd/main.go

clean:
	rm -rf ./log/calendar.log

lint:
	golangci-lint run ./...
	go vet ./...

# Заполнение календаря различными событиями
create_many: create_event create_event2 create_event3 create_event4

# Тестовые запросы с ошибками
test_errors: test_empty_id test_invalid_date test_not_found test_duplicate_id test_missing_user_id


create_event:
	curl -i -X POST http://localhost:8080/create_event \
  	-H "Content-Type: application/json" \
 	-d '{"id": "1", "user_id": "user1", "date": "2025-08-10T15:00:00Z", "text": "Тестовое событие"}'

create_event2:
	curl -i -X POST http://localhost:8080/create_event \
  	-H "Content-Type: application/json" \
 	-d '{"id": "2", "user_id": "user1", "date": "2025-08-11T15:00:00Z", "text": "Тестовое событие через день"}'

create_event3:
	curl -i -X POST http://localhost:8080/create_event \
  	-H "Content-Type: application/json" \
 	-d '{"id": "3", "user_id": "user1", "date": "2025-08-17T15:00:00Z", "text": "Тестовое событие через неделю"}'

create_event4:
	curl -i -X POST http://localhost:8080/create_event \
  	-H "Content-Type: application/json" \
 	-d '{"id": "4", "user_id": "user1", "date": "2025-09-10T15:00:00Z", "text": "Тестовое событие через месяц"}'

update_event:
	curl -i -X POST http://localhost:8080/update_event \
  	-H "Content-Type: application/json" \
  	-d '{"id": "1", "user_id": "user1", "date": "2025-08-11T16:00:00Z", "text": "Обновленное событие"}'

delete_event:
	curl -i -X POST http://localhost:8080/delete_event \
  	-H "Content-Type: application/json" \
  	-d '{"id": "1"}'

get_day_events:
	curl -i "http://localhost:8080/events_for_day?user_id=user1&date=2025-08-10"


get_week_events:
	curl -i "http://localhost:8080/events_for_week?user_id=user1&date=2025-08-10"

get_month_events:
	curl -i "http://localhost:8080/events_for_month?user_id=user1&date=2025-08-10"



# Пустой ID
test_empty_id:
	@echo "\n=== Тест: Пустой ID ==="
	curl -i -X POST http://localhost:8080/create_event \
		-H "Content-Type: application/x-www-form-urlencoded" \
		-d "id=&user_id=user1&date=2025-08-10T15:00:00Z&text=Событие с пустым ID"

# Неверный формат даты
test_invalid_date:
	@echo "\n=== Тест: Неверный формат даты ==="
	curl -i -X POST http://localhost:8080/create_event \
		-H "Content-Type: application/x-www-form-urlencoded" \
		-d "id=10&user_id=user1&date=2025-08-32T15:00:00Z&text=Событие с неверной датой"

# Попытка обновления несуществующего события
test_not_found:
	@echo "\n=== Тест: Обновление несуществующего события ==="
	curl -i -X POST http://localhost:8080/update_event \
		-H "Content-Type: application/x-www-form-urlencoded" \
		-d "id=999&user_id=user1&date=2025-08-10T15:00:00Z&text=Несуществующее событие"

# Дублирование ID
test_duplicate_id:
	@echo "\n=== Тест: Дублирование ID ==="
	curl -i -X POST http://localhost:8080/create_event \
		-H "Content-Type: application/x-www-form-urlencoded" \
		-d "id=1&user_id=user1&date=2025-08-10T15:00:00Z&text=Дубликат события"
	@echo "\nПовторная попытка:"
	curl -i -X POST http://localhost:8080/create_event \
		-H "Content-Type: application/x-www-form-urlencoded" \
		-d "id=1&user_id=user1&date=2025-08-10T15:00:00Z&text=Дубликат события"

# Отсутствие user_id
test_missing_user_id:
	@echo "\n=== Тест: Отсутствует user_id ==="
	curl -i -X POST http://localhost:8080/create_event \
		-H "Content-Type: application/x-www-form-urlencoded" \
		-d "id=11&date=2025-08-10T15:00:00Z&text=Событие без user_id"

# Получение событий без указания даты
test_missing_date:
	@echo "\n=== Тест: Отсутствует дата ==="
	curl -i "http://localhost:8080/events_for_day?user_id=user1"

# Удаление с неверным ID
test_delete_invalid:
	@echo "\n=== Тест: Удаление с неверным ID ==="
	curl -i -X POST http://localhost:8080/delete_event \
		-H "Content-Type: application/x-www-form-urlencoded" \
		-d "id=invalid_id"