#!/bin/bash

echo "=== ТЕСТИРОВАНИЕ EVENT SERVICE ==="

# Проверка здоровья
echo -e "\n1. Проверка здоровья сервиса:"
curl -s "http://localhost:8080/health" | jq .

# Создание событий
echo -e "\n2. Создание тестовых событий:"

echo "Событие 1 (без напоминания):"
curl -s -X POST http://localhost:8080/create_event \
  -H "Content-Type: application/json" \
  -d '{"usertoken": "test_user", "datetime": "2025-11-19T15:00:00Z", "text": "Первое тестовое событие", "title": "Тест 1"}' | jq .

echo "Событие 2 (с напоминанием):"
curl -s -X POST http://localhost:8080/create_event \
  -H "Content-Type: application/json" \
  -d '{"usertoken": "test_user", "datetime": "2025-11-20T15:00:00Z", "text": "Событие с напоминанием", "title": "Тест 2", "remind_before": 3600000000000, "telegram_id": 123456789}' | jq .

# Получение событий
echo -e "\n3. Получение событий на день:"
curl -s "http://localhost:8080/events_for_day?usertoken=test_user&date=2025-11-19" | jq .

echo -e "\n4. Получение событий на неделю:"
curl -s "http://localhost:8080/events_for_week?usertoken=test_user&date=2025-11-19" | jq .

# Обновление события
echo -e "\n5. Обновление события:"
curl -s -X POST http://localhost:8080/update_event \
  -H "Content-Type: application/json" \
  -d '{"event_id": 1, "datetime": "2025-11-19T16:00:00Z", "text": "Обновленное событие", "title": "Обновленный тест", "remind_before": 1800000000000, "telegram_id": 123456789}' | jq .

# Тест ошибок
echo -e "\n6. Тест ошибок:"

echo "Пустой usertoken:"
curl -s -X POST http://localhost:8080/create_event \
  -H "Content-Type: application/json" \
  -d '{"usertoken": "", "datetime": "2025-11-20T15:00:00Z", "text": "Ошибка", "title": "Тест"}' | jq .

echo "Событие в прошлом:"
curl -s -X POST http://localhost:8080/create_event \
  -H "Content-Type: application/json" \
  -d '{"usertoken": "test_user", "datetime": "2020-01-01T15:00:00Z", "text": "Прошлое событие", "title": "Тест"}' | jq .

echo -e "\n=== ТЕСТИРОВАНИЕ ЗАВЕРШЕНО ==="