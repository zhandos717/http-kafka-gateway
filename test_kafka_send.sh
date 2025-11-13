#!/bin/bash

# Скрипт для тестирования отправки сообщений в Kafka топик через API
# Использование: ./test_kafka_send.sh

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Параметры по умолчанию
API_URL="http://localhost:8080"
TOPIC="test-topic"
API_KEY="default-api-key"

# Функция для отображения справки
show_help() {
    echo "Скрипт для тестирования отправки сообщений в Kafka топик через API"
    echo ""
    echo "Использование: $0 [ОПЦИИ]"
    echo ""
    echo "Опции:"
    echo "  -u, --url URL       URL API (по умолчанию: $API_URL)"
    echo "  -t, --topic TOPIC   Название топика (по умолчанию: $TOPIC)"
    echo "  -k, --api-key KEY   API ключ для аутентификации"
    echo "  -h, --help          Показать эту справку"
    echo ""
    echo "Примеры:"
    echo "  $0                                  # Отправить сообщение в test-topic"
    echo "  $0 -t my-topic                      # Отправить сообщение в my-topic"
    echo "  $0 -u http://localhost:8080 -k abc123  # С указанием URL и API ключа"
}

# Парсинг аргументов
while [[ $# -gt 0 ]]; do
    case $1 in
        -u|--url)
            API_URL="$2"
            shift 2
            ;;
        -t|--topic)
            TOPIC="$2"
            shift 2
            ;;
        -k|--api-key)
            API_KEY="$2"
            shift 2
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            echo "Неизвестный параметр: $1"
            show_help
            exit 1
            ;;
    esac
done

echo -e "${YELLOW}=== Тестирование отправки сообщения в Kafka топик ===${NC}"
echo "API URL: $API_URL"
echo "Топик: $TOPIC"
echo ""

# Проверка доступности API
echo -e "${YELLOW}Проверка доступности API...${NC}"
if curl -s --connect-timeout 5 "$API_URL/health" > /dev/null; then
    echo -e "${GREEN}✓ API доступен${NC}"
else
    echo -e "${RED}✗ API недоступен по адресу $API_URL${NC}"
    echo -e "${RED}Пожалуйста, убедитесь, что Kafka Gateway запущен${NC}"
    exit 1
fi

echo ""

# Подготовка заголовков
HEADERS=("-H" "Content-Type: application/json")
if [ -n "$API_KEY" ]; then
    HEADERS+=("-H" "Authorization: Bearer $API_KEY")
else
    echo -e "${YELLOW}Предупреждение: API ключ не указан. Запрос может завершиться с ошибкой 401.${NC}"
fi

# Тест 1: Отправка простого сообщения
echo -e "${YELLOW}Тест 1: Отправка простого сообщения${NC}"
MESSAGE_DATA=$(cat <<EOF
{
    "topic": "$TOPIC",
    "key": "test-key-$(date +%s)",
    "value": {
        "message": "Hello from curl test!",
        "timestamp": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
        "test_id": "$(uuidgen 2>/dev/null || echo 'test-$(date +%s)-$(random)')"
    }
}
EOF
)

echo "Отправляем сообщение:"
echo "$MESSAGE_DATA" | jq . 2>/dev/null || echo "$MESSAGE_DATA"
echo ""

RESPONSE=$(curl -s -w "\n%{http_code}" "${HEADERS[@]}" \
    -X POST "$API_URL/message" \
    -d "$MESSAGE_DATA")

HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
RESPONSE_BODY=$(echo "$RESPONSE" | sed '$d')

if [ "$HTTP_CODE" -eq 200 ]; then
    echo -e "${GREEN}✓ Сообщение успешно отправлено${NC}"
    echo "Ответ сервера:"
    echo "$RESPONSE_BODY" | jq . 2>/dev/null || echo "$RESPONSE_BODY"
else
    echo -e "${RED}✗ Ошибка при отправке сообщения (HTTP $HTTP_CODE)${NC}"
    echo "Ответ сервера:"
    echo "$RESPONSE_BODY" | jq . 2>/dev/null || echo "$RESPONSE_BODY"
fi

echo ""

# Тест 2: Отправка сообщения с заголовками
echo -e "${YELLOW}Тест 2: Отправка сообщения с заголовками${NC}"
MESSAGE_WITH_HEADERS=$(cat <<EOF
{
    "topic": "$TOPIC",
    "key": "test-key-with-headers-$(date +%s)",
    "value": {
        "message": "Hello with headers!",
        "timestamp": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
        "type": "test"
    },
    "headers": {
        "Content-Type": "application/json",
        "Source": "curl-test",
        "Test-ID": "$(uuidgen 2>/dev/null || echo 'test-$(date +%s)-$(random)')"
    }
}
EOF
)

echo "Отправляем сообщение с заголовками:"
echo "$MESSAGE_WITH_HEADERS" | jq . 2>/dev/null || echo "$MESSAGE_WITH_HEADERS"
echo ""

RESPONSE=$(curl -s -w "\n%{http_code}" "${HEADERS[@]}" \
    -X POST "$API_URL/message" \
    -d "$MESSAGE_WITH_HEADERS")

HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
RESPONSE_BODY=$(echo "$RESPONSE" | sed '$d')

if [ "$HTTP_CODE" -eq 200 ]; then
    echo -e "${GREEN}✓ Сообщение с заголовками успешно отправлено${NC}"
    echo "Ответ сервера:"
    echo "$RESPONSE_BODY" | jq . 2>/dev/null || echo "$RESPONSE_BODY"
else
    echo -e "${RED}✗ Ошибка при отправке сообщения с заголовками (HTTP $HTTP_CODE)${NC}"
    echo "Ответ сервера:"
    echo "$RESPONSE_BODY" | jq . 2>/dev/null || echo "$RESPONSE_BODY"
fi

echo ""

# Тест 3: Попытка отправки сообщения с неверным топиком
echo -e "${YELLOW}Тест 3: Попытка отправки в неверный топик${NC}"
INVALID_TOPIC_MESSAGE=$(cat <<EOF
{
    "topic": ".invalid-topic-starting-with-dot",
    "key": "test-key-invalid",
    "value": "This should fail"
}
EOF
)

echo "Отправляем сообщение в неверный топик:"
echo "$INVALID_TOPIC_MESSAGE" | jq . 2>/dev/null || echo "$INVALID_TOPIC_MESSAGE"
echo ""

RESPONSE=$(curl -s -w "\n%{http_code}" "${HEADERS[@]}" \
    -X POST "$API_URL/message" \
    -d "$INVALID_TOPIC_MESSAGE")

HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
RESPONSE_BODY=$(echo "$RESPONSE" | sed '$d')

if [ "$HTTP_CODE" -eq 400 ]; then
    echo -e "${GREEN}✓ Корректно отклонено сообщение с неверным топиком${NC}"
    echo "Ответ сервера:"
    echo "$RESPONSE_BODY" | jq . 2>/dev/null || echo "$RESPONSE_BODY"
else
    echo -e "${YELLOW}! Неожиданный статус при отправке в неверный топик (HTTP $HTTP_CODE)${NC}"
    echo "Ответ сервера:"
    echo "$RESPONSE_BODY" | jq . 2>/dev/null || echo "$RESPONSE_BODY"
fi

echo ""
echo -e "${GREEN}=== Тестирование завершено ===${NC}"

# Сводка
echo ""
echo "Сводка:"
echo "- Проверена доступность API"
echo "- Отправлено простое сообщение"
echo "- Отправлено сообщение с заголовками" 
echo "- Проверена валидация топика"