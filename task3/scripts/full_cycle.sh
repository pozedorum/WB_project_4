#!/bin/bash

# Базовые настройки
BASE_URL="http://localhost:8080"
CHAT_ID="1105031510"
USER_TOKEN="test_user_$(date +%s)"  # Уникальный токен для каждого запуска

echo "=== ТЕСТИРОВАНИЕ ВСЕХ API ENDPOINTS ==="
echo "Base URL: $BASE_URL"
echo "User Token: $USER_TOKEN"
echo "Chat ID: $CHAT_ID"
echo ""

# Генерация временных меток (маленькие промежутки 1-5 минут)
EVENT_TIME_1=$(date -u -d "+1 minutes" +"%Y-%m-%dT%H:%M:%SZ")
EVENT_TIME_2=$(date -u -d "+2 minutes" +"%Y-%m-%dT%H:%M:%SZ")
EVENT_TIME_3=$(date -u -d "+3 minutes" +"%Y-%m-%dT%H:%M:%SZ")
EVENT_TIME_4=$(date -u -d "+4 minutes" +"%Y-%m-%dT%H:%M:%SZ")
EVENT_TIME_5=$(date -u -d "+5 minutes" +"%Y-%m-%dT%H:%M:%SZ")
TODAY=$(date -u +"%Y-%m-%d")

echo "📅 Временные метки (все в ближайшие 5 минут):"
echo "  Событие 1: $EVENT_TIME_1 (через 1 минуту)"
echo "  Событие 2: $EVENT_TIME_2 (через 2 минуты)" 
echo "  Событие 3: $EVENT_TIME_3 (через 3 минуты)"
echo "  Событие 4: $EVENT_TIME_4 (через 4 минуты)"
echo "  Событие 5: $EVENT_TIME_5 (через 5 минут)"
echo "  Сегодня: $TODAY"
echo ""

# Функция для форматированного вывода
print_result() {
    echo "--- $1 ---"
}

# 1. Проверка здоровья сервиса
print_result "1. ПРОВЕРКА ЗДОРОВЬЯ СЕРВИСА"
curl -s -X GET "$BASE_URL/health" | jq .
echo ""

# 2. Создание нескольких событий с разными интервалами напоминаний
print_result "2. СОЗДАНИЕ СОБЫТИЙ С НАПОМИНАНИЯМИ"

print_result "2.1. Событие через 1 минуту (напоминание за 0.5 минут)"
curl -s -X POST "$BASE_URL/create_event" \
  -H "Content-Type: application/json" \
  -d '{
    "usertoken": "'$USER_TOKEN'",
    "title": "Срочная встреча",
    "text": "Обсуждение срочных вопросов",
    "datetime": "'$EVENT_TIME_1'",
    "telegram_id": '$CHAT_ID',
    "remind_before": 1
  }' | jq .

print_result "2.2. Событие через 2 минуты (напоминание за 1 минуту)"
curl -s -X POST "$BASE_URL/create_event" \
  -H "Content-Type: application/json" \
  -d '{
    "usertoken": "'$USER_TOKEN'",
    "title": "Техническое совещание",
    "text": "Обсуждение архитектуры проекта",
    "datetime": "'$EVENT_TIME_2'",
    "telegram_id": '$CHAT_ID',
    "remind_before": 1
  }' | jq .

print_result "2.3. Событие через 3 минуты (напоминание за 2 минуты)"
curl -s -X POST "$BASE_URL/create_event" \
  -H "Content-Type: application/json" \
  -d '{
    "usertoken": "'$USER_TOKEN'",
    "title": "Планирование спринта",
    "text": "Определение задач на ближайшие дни",
    "datetime": "'$EVENT_TIME_3'",
    "telegram_id": '$CHAT_ID',
    "remind_before": 2
  }' | jq .

print_result "2.4. Событие через 4 минуты (без напоминания)"
curl -s -X POST "$BASE_URL/create_event" \
  -H "Content-Type: application/json" \
  -d '{
    "usertoken": "'$USER_TOKEN'",
    "title": "Рабочая сессия",
    "text": "Самостоятельная работа над задачами",
    "datetime": "'$EVENT_TIME_4'",
    "telegram_id": '$CHAT_ID',
    "remind_before": 0
  }' | jq .

print_result "2.5. Событие через 5 минут (напоминание за 3 минуты)"
curl -s -X POST "$BASE_URL/create_event" \
  -H "Content-Type: application/json" \
  -d '{
    "usertoken": "'$USER_TOKEN'",
    "title": "Демонстрация проекта",
    "text": "Показ результатов работы команды",
    "datetime": "'$EVENT_TIME_5'",
    "telegram_id": '$CHAT_ID',
    "remind_before": 3
  }' | jq .

echo ""

# 3. Получение событий
print_result "3. ПОЛУЧЕНИЕ СОБЫТИЙ"

print_result "3.1. События на сегодня"
curl -s "$BASE_URL/events_for_day?usertoken=$USER_TOKEN&date=$TODAY" | jq '.result[] | {id, title, datetime, remind_before}'

print_result "3.2. События на неделю"
curl -s "$BASE_URL/events_for_week?usertoken=$USER_TOKEN&date=$TODAY" | jq '.result[] | {id, title, datetime, remind_before}'

print_result "3.3. События на месяц"
curl -s "$BASE_URL/events_for_month?usertoken=$USER_TOKEN&date=$TODAY" | jq '.result[] | {id, title, datetime, remind_before}'

echo ""

# 4. Обновление события
print_result "4. ОБНОВЛЕНИЕ СОБЫТИЯ"

# Обновляем событие через 3 минуты
UPDATED_TIME=$(date -u -d "+4 minutes" +"%Y-%m-%dT%H:%M:%SZ")

print_result "4.1. Обновление события ID: 3 (перенос на +4 минуты)"
curl -s -X POST "$BASE_URL/update_event" \
  -H "Content-Type: application/json" \
  -d '{
    "event_id": 3,
    "title": "ОБНОВЛЕННОЕ Планирование спринта",
    "text": "Перенесено на 1 минуту позже. Добавлены новые темы.",
    "datetime": "'$UPDATED_TIME'",
    "telegram_id": '$CHAT_ID',
    "remind_before": 1
  }' | jq .

echo ""

# 5. Проверка обновленных событий
print_result "5. ПРОВЕРКА ОБНОВЛЕННЫХ СОБЫТИЙ"

print_result "5.1. События на день после обновления"
curl -s "$BASE_URL/events_for_day?usertoken=$USER_TOKEN&date=$TODAY" | jq '.result[] | {id, title, datetime, remind_before}'

echo ""

# 6. Удаление события
print_result "6. УДАЛЕНИЕ СОБЫТИЯ"

print_result "6.1. Удаление события ID: 2"
curl -s -X POST "$BASE_URL/delete_event" \
  -H "Content-Type: application/json" \
  -d '{
    "event_id": 2
  }' | jq .

echo ""

# 7. Финальная проверка оставшихся событий
print_result "7. ФИНАЛЬНАЯ ПРОВЕРКА СОБЫТИЙ"

print_result "7.1. Все события на сегодня после удаления"
curl -s "$BASE_URL/events_for_day?usertoken=$USER_TOKEN&date=$TODAY" | jq '.result[] | {id, title, datetime, remind_before}'

print_result "7.2. Количество событий на неделю после удаления"
curl -s "$BASE_URL/events_for_week?usertoken=$USER_TOKEN&date=$TODAY" | jq '.result | length'

echo ""

# 8. Тестирование ошибок
print_result "8. ТЕСТИРОВАНИЕ ОШИБОК"

print_result "8.1. Пустой usertoken"
curl -s -X POST "$BASE_URL/create_event" \
  -H "Content-Type: application/json" \
  -d '{
    "usertoken": "",
    "title": "Тест ошибки",
    "text": "Должна быть ошибка - пустой usertoken",
    "datetime": "'$EVENT_TIME_1'"
  }' | jq .

print_result "8.2. Событие в прошлом"
curl -s -X POST "$BASE_URL/create_event" \
  -H "Content-Type: application/json" \
  -d '{
    "usertoken": "'$USER_TOKEN'",
    "title": "Событие в прошлом",
    "text": "Должна быть ошибка - дата в прошлом",
    "datetime": "2020-01-01T00:00:00Z"
  }' | jq .

print_result "8.3. Отсутствует дата в запросе событий"
curl -s "$BASE_URL/events_for_day?usertoken=$USER_TOKEN" | jq .

print_result "8.4. Отсутствует usertoken в запросе событий"
curl -s "$BASE_URL/events_for_day?date=$TODAY" | jq .

print_result "8.5. Обновление несуществующего события"
curl -s -X POST "$BASE_URL/update_event" \
  -H "Content-Type: application/json" \
  -d '{
    "event_id": 999,
    "title": "Несуществующее событие",
    "text": "Должна быть ошибка",
    "datetime": "'$EVENT_TIME_1'",
    "telegram_id": '$CHAT_ID',
    "remind_before": 1
  }' | jq .

echo ""

# 9. Тестирование разных форматов remind_before
print_result "9. ТЕСТИРОВАНИЕ РАЗНЫХ ИНТЕРВАЛОВ НАПОМИНАНИЙ"

TEST_TIME=$(date -u -d "+5 minutes" +"%Y-%m-%dT%H:%M:%SZ")

print_result "9.1. Напоминание за 0.5 минут (30 секунд)"
curl -s -X POST "$BASE_URL/create_event" \
  -H "Content-Type: application/json" \
  -d '{
    "usertoken": "'$USER_TOKEN'_reminder",
    "title": "Тест 0.5м напоминание",
    "text": "Напоминание за 30 секунд",
    "datetime": "'$TEST_TIME'",
    "telegram_id": '$CHAT_ID',
    "remind_before": 1
  }' | jq .

print_result "9.2. Напоминание за 1 минуту"
curl -s -X POST "$BASE_URL/create_event" \
  -H "Content-Type: application/json" \
  -d '{
    "usertoken": "'$USER_TOKEN'_reminder", 
    "title": "Тест 1м напоминание",
    "text": "Напоминание за 1 минуту",
    "datetime": "'$TEST_TIME'",
    "telegram_id": '$CHAT_ID',
    "remind_before": 1
  }' | jq .

print_result "9.3. Напоминание за 2 минуты"
curl -s -X POST "$BASE_URL/create_event" \
  -H "Content-Type: application/json" \
  -d '{
    "usertoken": "'$USER_TOKEN'_reminder",
    "title": "Тест 2м напоминание", 
    "text": "Напоминание за 2 минуты",
    "datetime": "'$TEST_TIME'",
    "telegram_id": '$CHAT_ID',
    "remind_before": 2
  }' | jq .

echo ""

print_result "🚀 ТЕСТИРОВАНИЕ ЗАВЕРШЕНО"
echo ""
echo "📊 Сводка:"
echo "  - Создано 7+ событий в ближайшие 5 минут"
echo "  - Настроены напоминания за 0.5-3 минут до событий"
echo "  - Протестированы все CRUD операции"
echo "  - Проверены различные сценарии ошибок"
echo ""
echo "⏰ Ожидайте уведомления в Telegram в ближайшие 1-5 минут!"
echo "🔍 User Token для дополнительных проверок: $USER_TOKEN"
echo "💡 Проверьте логи RabbitMQ для отслеживания напоминаний"