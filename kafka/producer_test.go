package kafka

import (
	"context"
	"testing"

	"github.com/segmentio/kafka-go"
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

// Тесты для проверки, что топик указывается в сообщении
// Для этих тестов мы создадим мок-объект, чтобы проверить, что топик передается в сообщении

// Мок-объект для Writer
type MockWriter struct {
	WriteMessagesFunc func(ctx context.Context, msgs ...kafka.Message) error
	CloseFunc         func() error
}

func (m *MockWriter) WriteMessages(ctx context.Context, msgs ...kafka.Message) error {
	if m.WriteMessagesFunc != nil {
		return m.WriteMessagesFunc(ctx, msgs...)
	}
	return nil
}

func (m *MockWriter) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

func (m *MockWriter) Address() string {
	return "mock"
}

func (m *MockWriter) Stats() kafka.WriterStats {
	return kafka.WriterStats{}
}

// Тест для проверки, что топик указывается в сообщении при использовании SendMessage
func TestProducerSendMessageTopicSpecification(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	// Создаем мок-объект для Writer
	mockWriter := &MockWriter{
		WriteMessagesFunc: func(ctx context.Context, msgs ...kafka.Message) error {
			// Проверяем, что в сообщении указан топик
			for _, msg := range msgs {
				if msg.Topic == "" {
					t.Errorf("Expected topic to be specified in message, got empty topic")
				}
			}
			return nil
		},
	}

	// Создаем продюсер с мок-объектом
	producer := &Producer{
		writer: mockWriter,
		logger: logger,
	}

	// Отправляем сообщение
	err := producer.SendMessage("test-topic", []byte("key"), []byte("value"))
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// Тест для проверки, что топик указывается в сообщении при использовании SendMessageWithHeaders
func TestProducerSendMessageWithHeadersTopicSpecification(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	// Создаем мок-объект для Writer
	mockWriter := &MockWriter{
		WriteMessagesFunc: func(ctx context.Context, msgs ...kafka.Message) error {
			// Проверяем, что в сообщении указан топик
			for _, msg := range msgs {
				if msg.Topic == "" {
					t.Errorf("Expected topic to be specified in message, got empty topic")
				}
			}
			return nil
		},
	}

	// Создаем продюсер с мок-объектом
	producer := &Producer{
		writer: mockWriter,
		logger: logger,
	}

	// Отправляем сообщение с заголовками
	headers := map[string]string{
		"header1": "value1",
	}
	err := producer.SendMessageWithHeaders("test-topic", []byte("key"), []byte("value"), headers)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
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
