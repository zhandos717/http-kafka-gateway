package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// WriterInterface определяет интерфейс для Kafka Writer
type WriterInterface interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

type Producer struct {
	writer WriterInterface
	logger *zap.Logger
}

func NewProducer(brokers string, logger *zap.Logger) *Producer {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(brokers),
		Balancer:               &kafka.Hash{},
		WriteTimeout:           10 * time.Second,
		ReadTimeout:            10 * time.Second,
		RequiredAcks:           kafka.RequireAll,
		MaxAttempts:            3,
		AllowAutoTopicCreation: true,
		// Указываем топик как пустую строку, так как будем указывать его в каждом сообщении
	}

	return &Producer{
		writer: writer,
		logger: logger,
	}
}

func (p *Producer) SendMessage(topic string, key, value []byte) error {
	message := kafka.Message{
		Topic: topic, // Указываем топик в сообщении
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}

	err := p.writer.WriteMessages(context.Background(), message)
	if err != nil {
		p.logger.Error("Failed to send message to Kafka",
			zap.String("topic", topic),
			zap.Error(err))
		return err
	}

	p.logger.Info("Message sent to Kafka",
		zap.String("topic", topic),
		zap.Int("value_length", len(value)))

	return nil
}

func (p *Producer) SendMessageWithHeaders(topic string, key, value []byte, headers map[string]string) error {
	// Преобразуем map[string]string в []kafka.Header
	kafkaHeaders := make([]kafka.Header, 0, len(headers))
	for k, v := range headers {
		kafkaHeaders = append(kafkaHeaders, kafka.Header{
			Key:   k,
			Value: []byte(v),
		})
	}

	message := kafka.Message{
		Topic:   topic, // Указываем топик в сообщении
		Key:     key,
		Value:   value,
		Headers: kafkaHeaders,
		Time:    time.Now(),
	}

	err := p.writer.WriteMessages(context.Background(), message)
	if err != nil {
		p.logger.Error("Failed to send message with headers to Kafka",
			zap.String("topic", topic),
			zap.Error(err))
		return err
	}

	p.logger.Info("Message with headers sent to Kafka",
		zap.String("topic", topic),
		zap.Int("value_length", len(value)))

	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
