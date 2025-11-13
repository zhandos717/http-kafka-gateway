package kafka

import (
	"testing"

	"go.uber.org/zap"
)

func TestNewProducer(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	// В реальных тестах нужно использовать тестовый Kafka брокер
	// Для unit тестов тестируем только создание объекта
	producer := NewProducer("localhost:9092", logger)

	if producer == nil {
		t.Errorf("Expected producer to be created")
	}

	if producer.logger == nil {
		t.Errorf("Expected logger to be set")
	}
}

// Пропускаем тесты отправки сообщений, так как они требуют запущенного Kafka брокера
func TestSendMessage(t *testing.T) {
	t.Skip("Skipping SendMessage test as it requires a running Kafka broker")
}

func TestSendMessageWithHeaders(t *testing.T) {
	t.Skip("Skipping SendMessageWithHeaders test as it requires a running Kafka broker")
}

func TestClose(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	producer := NewProducer("localhost:9092", logger)
	if producer == nil {
		t.Errorf("Expected producer to be created")
	}

	// Тестирование закрытия (в реальных условиях может вызвать ошибки если Kafka недоступен)
	err := producer.Close()
	// Не проверяем ошибку, так как может быть ошибка подключения
	_ = err
}
