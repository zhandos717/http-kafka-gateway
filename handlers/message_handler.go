package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"kafkaGateway/metrics"
	"kafkaGateway/models"
	"kafkaGateway/utils"
)

// Интерфейс для Producer, чтобы можно было использовать мок
type ProducerInterface interface {
	SendMessage(topic string, key, value []byte) error
	SendMessageWithHeaders(topic string, key, value []byte, headers map[string]string) error
	Close() error
}

type MessageHandler struct {
	producer ProducerInterface
	logger   *zap.Logger
}

func NewMessageHandler(producer ProducerInterface, logger *zap.Logger) *MessageHandler {
	return &MessageHandler{
		producer: producer,
		logger:   logger,
	}

}

func (mh *MessageHandler) SendMessage(c *gin.Context) {
	startTime := time.Now()

	var req models.MessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mh.logger.Error("Invalid request format", zap.Error(err))
		metrics.RequestDuration.WithLabelValues("POST", "/message").Observe(time.Since(startTime).Seconds())
		metrics.HTTPLatency.WithLabelValues("/message", "POST", "400").Observe(time.Since(startTime).Seconds())

		c.JSON(http.StatusBadRequest, models.MessageResponse{
			Success:   false,
			Error:     "Invalid request format: " + err.Error(),
			Timestamp: time.Now(),
		})
		metrics.AuthAttempts.WithLabelValues("failed").Inc()
		return
	}

	// Проверяем валидность топика
	if !utils.IsValidTopic(req.Topic) {
		mh.logger.Error("Invalid topic name", zap.String("topic", req.Topic))
		metrics.RequestDuration.WithLabelValues("POST", "/message").Observe(time.Since(startTime).Seconds())
		metrics.HTTPLatency.WithLabelValues("/message", "POST", "400").Observe(time.Since(startTime).Seconds())

		c.JSON(http.StatusBadRequest, models.MessageResponse{
			Success:   false,
			Error:     "Invalid topic name",
			Timestamp: time.Now(),
		})
		metrics.AuthAttempts.WithLabelValues("failed").Inc()
		return
	}

	// Конвертируем значение в байты
	valueBytes, err := utils.ConvertInterfaceToBytes(req.Value)
	if err != nil {
		mh.logger.Error("Failed to convert message value to bytes", zap.Error(err))
		metrics.RequestDuration.WithLabelValues("POST", "/message").Observe(time.Since(startTime).Seconds())
		metrics.HTTPLatency.WithLabelValues("/message", "POST", "500").Observe(time.Since(startTime).Seconds())

		c.JSON(http.StatusInternalServerError, models.MessageResponse{
			Success:   false,
			Error:     "Failed to convert message value: " + err.Error(),
			Timestamp: time.Now(),
		})
		metrics.AuthAttempts.WithLabelValues("failed").Inc()
		return
	}

	// Конвертируем ключ в байты, если он есть
	var keyBytes []byte
	if req.Key != "" {
		keyBytes = []byte(req.Key)
	}

	// Отправляем сообщение в Kafka
	var sendErr error
	if len(req.Headers) > 0 {
		// Преобразуем строки заголовков в байты
		headerBytes := make(map[string]string)
		for k, v := range req.Headers {
			headerBytes[k] = v
		}

		sendErr = mh.producer.SendMessageWithHeaders(req.Topic, keyBytes, valueBytes, headerBytes)
	} else {
		sendErr = mh.producer.SendMessage(req.Topic, keyBytes, valueBytes)
	}

	if sendErr != nil {
		mh.logger.Error("Failed to send message to Kafka",
			zap.String("topic", req.Topic),
			zap.Error(sendErr))

		metrics.KafkaErrors.WithLabelValues(req.Topic, "send_error").Inc()
		metrics.RequestDuration.WithLabelValues("POST", "/message").Observe(time.Since(startTime).Seconds())
		metrics.HTTPLatency.WithLabelValues("/message", "POST", "500").Observe(time.Since(startTime).Seconds())

		c.JSON(http.StatusInternalServerError, models.MessageResponse{
			Success:   false,
			Error:     "Failed to send message to Kafka: " + sendErr.Error(),
			Timestamp: time.Now(),
		})
		metrics.AuthAttempts.WithLabelValues("failed").Inc()
		return
	}

	// Успешная отправка
	mh.logger.Info("Message sent to Kafka successfully",
		zap.String("topic", req.Topic),
		zap.String("key", req.Key),
		zap.Int("value_length", len(valueBytes)))

	metrics.MessagesProcessed.WithLabelValues(req.Topic, "success").Inc()
	metrics.RequestDuration.WithLabelValues("POST", "/message").Observe(time.Since(startTime).Seconds())
	metrics.HTTPLatency.WithLabelValues("/message", "POST", "200").Observe(time.Since(startTime).Seconds())

	c.JSON(http.StatusOK, models.MessageResponse{
		Success:   true,
		Message:   "Message sent to Kafka successfully",
		Timestamp: time.Now(),
	})
	metrics.AuthAttempts.WithLabelValues("success").Inc()
}
