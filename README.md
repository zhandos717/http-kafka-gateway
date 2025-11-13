# Kafka Gateway

Kafka Gateway - это HTTP-сервер, который принимает сообщения через HTTP API и отправляет их в Apache Kafka. Сервер включает в себя аутентификацию по API-ключам, логирование и мониторинг.

## Функциональность

- HTTP-сервер для приема сообщений
- Аутентификация через API-ключи
- Отправка сообщений в Kafka
- Логирование операций
- Мониторинг с метриками Prometheus
- Валидация топиков и сообщений

## Установка и запуск

1. Установите зависимости:
```bash
go mod tidy
```

2. Настройте переменные окружения в файле `.env`:
```env
KAFKA_BROKERS=localhost:9092
KAFKA_LOG_LEVEL=3
SERVER_PORT=8080
API_KEYS=your-api-key-here
```

3. Запустите сервер:
```bash
go run main.go
```

## Использование

### Отправка сообщения в Kafka

```bash
curl -X POST http://localhost:8080/message \
  -H "Authorization: Bearer your-api-key-here" \
  -H "Content-Type: application/json" \
  -d '{
    "topic": "your-topic-name",
    "key": "message-key",
    "value": "Hello, Kafka!",
    "headers": {
      "custom-header": "header-value"
    }
  }'
```

### Проверка состояния

```bash
curl http://localhost:8080/health
```

### Получение метрик

```bash
curl http://localhost:8080/metrics
```

## API

### POST /message

Отправляет сообщение в указанный топик Kafka.

**Headers:**
- `Authorization: Bearer <api-key>` - обязательный заголовок с API-ключом

**Body:**
```json
{
  "topic": "string",      // Название топика (обязательно)
  "key": "string",        // Ключ сообщения (опционально)
  "value": "any",         // Значение сообщения (обязательно)
  "headers": {            // Заголовки сообщения (опционально)
    "header-name": "header-value"
  }
}
```

**Ответы:**
- `200 OK` - сообщение успешно отправлено
- `400 Bad Request` - неверный формат запроса
- `401 Unauthorized` - неверный или отсутствующий API-ключ
- `500 Internal Server Error` - ошибка при отправке в Kafka

### GET /health

Проверяет состояние сервера.

### GET /metrics

Возвращает метрики в формате Prometheus.

## Архитектура

Проект состоит из следующих модулей:

- `config` - загрузка и управление конфигурацией
- `handlers` - обработчики HTTP-запросов
- `middleware` - промежуточное ПО (аутентификация)
- `kafka` - взаимодействие с Kafka
- `logger` - система логирования
- `metrics` - система метрик
- `models` - модели данных
- `utils` - вспомогательные функции

## Метрики

Сервер экспортирует следующие метрики Prometheus:

- `kafka_gateway_messages_processed_total` - количество обработанных сообщений
- `kafka_gateway_request_duration_seconds` - время обработки запросов
- `kafka_gateway_kafka_errors_total` - количество ошибок при отправке в Kafka
- `kafka_gateway_auth_attempts_total` - количество попыток аутентификации
- `kafka_gateway_http_response_time_seconds` - время отклика HTTP-эндпоинтов