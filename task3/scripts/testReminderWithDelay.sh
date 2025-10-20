#!/bin/bash

CHAT_ID="1105031510"

# Время события: через 3 минуты
EVENT_TIME_UTC=$(date -u -d "+2 minutes" +"%Y-%m-%dT%H:%M:%SZ")

echo "=== ПРОСТОЙ ТЕСТ REMIND_BEFORE ==="
echo "Событие: через 2 минуты ($EVENT_TIME_UTC)"
echo "Напоминание: через 1 минуту (remind_before: 1 минуты)"

# Основной тест
curl -X POST http://localhost:8080/create_event \
  -H "Content-Type: application/json" \
  -d '{
    "usertoken": "simple_remind_test",
    "datetime": "'$EVENT_TIME_UTC'",
    "text": "Событие через 2 минуты. Напоминание придет через 1 минуту (за 2 минуты до события).", 
    "title": "Тест remind_before",
    "remind_before": 1,
    "telegram_id": '$CHAT_ID'
  }' | jq .

echo -e "\nТест создан! Ожидайте уведомление через ~1 минуту"
