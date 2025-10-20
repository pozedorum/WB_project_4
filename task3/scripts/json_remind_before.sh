#!/bin/bash

CHAT_ID="1105031510"
EVENT_TIME=$(date -u -d "+6 minutes" +"%Y-%m-%dT%H:%M:%SZ")

echo "=== ТЕСТ REMIND_BEFORE В JSON ==="
echo "Event time: '$EVENT_TIME'"
echo "1. 🔔 Напоминание за 1 минуту:"
curl -s -X POST http://localhost:8080/create_event \
  -H "Content-Type: application/json" \
  -d '{
    "usertoken": "duration_test_1",
    "datetime": "'$EVENT_TIME'",
    "text": "Напоминание за 1 минуту до события",
    "title": "⏰ 1m",
    "remind_before": 1,
    "telegram_id": '$CHAT_ID'
  }' | jq '.result // .error'

echo -e "\n2. 🕑 Напоминание за 2 минуты:"
curl -s -X POST http://localhost:8080/create_event \
  -H "Content-Type: application/json" \
  -d '{
    "usertoken": "duration_test_2", 
    "datetime": "'$EVENT_TIME'",
    "text": "Напоминание за 2 минуты до события",
    "title": "🕑 2m",
    "remind_before": 2,
    "telegram_id": '$CHAT_ID'
  }' | jq '.result // .error'

echo -e "\n3. 🕔 Напоминание за 5 минут:"
curl -s -X POST http://localhost:8080/create_event \
  -H "Content-Type: application/json" \
  -d '{
    "usertoken": "duration_test_3",
    "datetime": "'$EVENT_TIME'",
    "text": "Напоминание за 5 минут до события",
    "title": "🕔 5m", 
    "remind_before": 5,
    "telegram_id": '$CHAT_ID'
  }' | jq '.result // .error'

echo -e "\n✅ Тесты завершены!"