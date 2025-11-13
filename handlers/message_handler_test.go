package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"kafkaGateway/kafka"
	"kafkaGateway/models"
)

// MockProducer - имитация Kafka Producer для тестирования
type MockProducer struct {
	SendMessageFunc            func(topic string, key, value []byte) error
	SendMessageWithHeadersFunc func(topic string, key, value []byte, headers map[string]string) error
	CloseFunc                  func() error
}

func (m *MockProducer) SendMessage(topic string, key, value []byte) error {
	if m.SendMessageFunc != nil {
		return m.SendMessageFunc(topic, key, value)
	}
	return nil
}

func (m *MockProducer) SendMessageWithHeaders(topic string, key, value []byte, headers map[string]string) error {
	if m.SendMessageWithHeadersFunc != nil {
		return m.SendMessageWithHeadersFunc(topic, key, value, headers)
	}
	return nil
}

func (m *MockProducer) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

// Создаем мокнуть, который соответствует интерфейсу Producer
type ProducerMock struct {
	*kafka.Producer
	MockSendMessage            func(topic string, key, value []byte) error
	MockSendMessageWithHeaders func(topic string, key, value []byte, headers map[string]string) error
}

func (p *ProducerMock) SendMessage(topic string, key, value []byte) error {
	if p.MockSendMessage != nil {
		return p.MockSendMessage(topic, key, value)
	}
	return nil
}

func (p *ProducerMock) SendMessageWithHeaders(topic string, key, value []byte, headers map[string]string) error {
	if p.MockSendMessageWithHeaders != nil {
		return p.MockSendMessageWithHeaders(topic, key, value, headers)
	}
	return nil
}

func TestNewMessageHandler(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	mockProducer := &ProducerMock{}
	handler := NewMessageHandler(mockProducer, logger)

	if handler.producer == nil {
		t.Errorf("Expected producer to be set")
	}

	if handler.logger == nil {
		t.Errorf("Expected logger to be set")
	}
}

func TestMessageHandler_SendMessage(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		request        models.MessageRequest
		expectedStatus int
		sendError      error
	}{
		{
			name: "valid message",
			request: models.MessageRequest{
				Topic: "test-topic",
				Key:   "test-key",
				Value: "test-value",
			},
			expectedStatus: http.StatusOK,
			sendError:      nil,
		},
		{
			name: "invalid topic name",
			request: models.MessageRequest{
				Topic: ".invalid-topic",
				Key:   "test-key",
				Value: "test-value",
			},
			expectedStatus: http.StatusBadRequest,
			sendError:      nil,
		},
		{
			name: "missing topic",
			request: models.MessageRequest{
				Key:   "test-key",
				Value: "test-value",
			},
			expectedStatus: http.StatusBadRequest,
			sendError:      nil,
		},
		{
			name: "missing value",
			request: models.MessageRequest{
				Topic: "test-topic",
				Key:   "test-key",
			},
			expectedStatus: http.StatusBadRequest,
			sendError:      nil,
		},
		{
			name: "send error",
			request: models.MessageRequest{
				Topic: "test-topic",
				Key:   "test-key",
				Value: "test-value",
			},
			expectedStatus: http.StatusInternalServerError,
			sendError:      &KafkaError{"failed to send"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем мок продюсера
			mockProducer := &ProducerMock{
				MockSendMessage: func(topic string, key, value []byte) error {
					return tt.sendError
				},
				MockSendMessageWithHeaders: func(topic string, key, value []byte, headers map[string]string) error {
					return tt.sendError
				},
			}

			handler := NewMessageHandler(mockProducer, logger)

			// Подготовка запроса
			jsonData, _ := json.Marshal(tt.request)
			req, _ := http.NewRequest("POST", "/message", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			// Создание контекста Gin
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Вызов обработчика
			handler.SendMessage(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Response body: %s", tt.expectedStatus, w.Code, w.Body.String())
			}
		})
	}
}

// KafkaError - имитация ошибки Kafka для тестирования
type KafkaError struct {
	msg string
}

func (e *KafkaError) Error() string {
	return e.msg
}

// Тест для проверки, что топик указывается правильно при отправке сообщения
func TestMessageHandler_SendMessageTopicSpecification(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	gin.SetMode(gin.TestMode)

	// Переменная для хранения полученного топика
	var receivedTopic string

	// Создаем мок продюсера, который сохраняет полученный топик
	mockProducer := &ProducerMock{
		MockSendMessage: func(topic string, key, value []byte) error {
			receivedTopic = topic
			return nil
		},
		MockSendMessageWithHeaders: func(topic string, key, value []byte, headers map[string]string) error {
			receivedTopic = topic
			return nil
		},
	}

	handler := NewMessageHandler(mockProducer, logger)

	// Подготовка запроса с топиком
	request := models.MessageRequest{
		Topic: "test-topic-for-verification",
		Key:   "test-key",
		Value: "test-value",
	}
	jsonData, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/message", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Создание контекста Gin
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Вызов обработчика
	handler.SendMessage(c)

	// Проверяем, что топик был передан правильно
	if receivedTopic != "test-topic-for-verification" {
		t.Errorf("Expected topic 'test-topic-for-verification', got '%s'", receivedTopic)
	}

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Response body: %s", http.StatusOK, w.Code, w.Body.String())
	}
}
