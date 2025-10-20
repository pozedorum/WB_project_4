#!/bin/bash

CHAT_ID="1105031510"

# Для Москвы (UTC+3) - корректное время
# Текущее время в UTC
TIME_NOW_UTC=$(date -u -d "+2 minutes" +"%Y-%m-%dT%H:%M:%SZ")
TIME_SOON_UTC=$(date -u -d "+4 minutes" +"%Y-%m-%dT%H:%M:%SZ") 

# Альтернативно - с явным указанием смещения
TIME_NOW_MSK=$(date -d "+2 minutes" +"%Y-%m-%dT%H:%M:%S+03:00")
TIME_SOON_MSK=$(date -d "+3 minutes" +"%Y-%m-%dT%H:%M:%S+03:00")

echo "=== ТЕСТ ДЛЯ МОСКВЫ (UTC+3) ==="
echo "UTC время:"
echo "- Сейчас+2мин UTC: $TIME_NOW_UTC"
echo "- Сейчас+4мин UTC: $TIME_SOON_UTC"
echo ""
echo "Московское время:"
echo "- Сейчас+2мин MSK: $TIME_NOW_MSK" 
echo "- Сейчас+4мин MSK: $TIME_SOON_MSK"

# Тест с UTC временем
echo -e "\n1. Тест с UTC временем:"
curl -s -X POST http://localhost:8080/create_event \
  -H "Content-Type: application/json" \
  -d '{
    "usertoken": "test_moscow",
    "datetime": "'$TIME_SOON_UTC'",
    "text": "Тест с UTC временем. Должно прийти в '$TIME_SOON_UTC' UTC", 
    "title": "UTC Тест",
    "telegram_id": '$CHAT_ID'
  }' | jq '.result // .error'

# Тест с московским временем
echo -e "\n2. 🇷🇺 Тест с московским временем:"
curl -s -X POST http://localhost:8080/create_event \
  -H "Content-Type: application/json" \
  -d '{
    "usertoken": "test_moscow", 
    "datetime": "'$TIME_SOON_MSK'",
    "text": "Тест с московским временем. Должно прийти в '$TIME_SOON_MSK'",
    "title": "🇷🇺 MSK Тест",
    "telegram_id": '$CHAT_ID'
  }' | jq '.result // .error'

echo -e "\nТесты созданы!"
echo "💡 Время в логах будет показываться в UTC"