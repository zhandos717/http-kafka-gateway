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

## Использование Makefile

Для удобства управления проектом доступен Makefile с основными командами:

```bash
# Сборка проекта
make build

# Запуск проекта
make run

# Установка зависимостей
make deps

# Запуск тестов
make test

# Запуск тестов с покрытием
make test-coverage

# Очистка артефактов сборки
make clean

# Форматирование кода
make fmt

# Проверка кода с помощью go vet
make vet

# Сборка для разных платформ
make build-linux    # Сборка для Linux
make build-windows  # Сборка для Windows
make build-macos    # Сборка для macOS
make build-all      # Сборка для всех платформ

# Просмотр всех доступных команд
make help
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

## Использование с PHP приложениями

Kafka Gateway позволяет PHP приложениям легко интегрироваться с Apache Kafka через простой HTTP API. Это упрощает отправку сообщений в Kafka без необходимости устанавливать и настраивать Kafka клиенты в PHP коде.

### Пример интеграции с PHP

```php
<?php
function sendToKafka($topic, $message, $apiKey) {
    $data = [
        'topic' => $topic,
        'key' => 'message-key',
        'value' => $message,
        'headers' => [
            'source' => 'php-app',
            'timestamp' => time()
        ]
    ];

    $options = [
        'http' => [
            'header' => [
                'Content-Type: application/json',
                'Authorization: Bearer ' . $apiKey
            ],
            'method' => 'POST',
            'content' => json_encode($data)
        ]
    ];

    $context = stream_context_create($options);
    $result = file_get_contents('http://localhost:8080/message', false, $context);
    
    return json_decode($result, true);
}

// Использование
$apiKey = 'your-api-key-here';
$result = sendToKafka('user-events', ['userId' => 123, 'action' => 'login'], $apiKey);
if ($result['success']) {
    echo "Сообщение успешно отправлено в Kafka";
} else {
    echo "Ошибка при отправке: " . $result['error'];
}
?>
```

Преимущества использования Kafka Gateway с PHP:
- Упрощенная интеграция без установки дополнительных библиотек
- Централизованное управление доступом через API-ключи
- Надежная отправка сообщений с обработкой ошибок
- Мониторинг и логирование всех операций

## Покрытие тестами

| Пакет | Покрытие |
|-------|----------|
| config | 93.3% |
| handlers | 79.6% |
| kafka | 80.0% |
| logger | 100.0% |
| middleware | 100.0% |
| utils | 92.3% |
| cmd/kafkaGateway | 0.0% |
| metrics | 0.0% |
| models | - |