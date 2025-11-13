package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// MessagesProcessed Количество обработанных сообщений
	MessagesProcessed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kafka_gateway_messages_processed_total",
			Help: "Total number of messages processed by the gateway",
		},
		[]string{"topic", "status"},
	)

	// RequestDuration Время обработки запросов
	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "kafka_gateway_request_duration_seconds",
			Help:    "Duration of requests to the gateway",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"method", "endpoint"},
	)

	// KafkaErrors Количество ошибок при отправке в Kafka
	KafkaErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kafka_gateway_kafka_errors_total",
			Help: "Total number of errors when sending messages to Kafka",
		},
		[]string{"topic", "error_type"},
	)

	// AuthAttempts Количество успешных/неуспешных аутентификаций
	AuthAttempts = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kafka_gateway_auth_attempts_total",
			Help: "Total number of authentication attempts",
		},
		[]string{"status"},
	)

	// HTTPLatency Время отклика HTTP-эндпоинтов
	HTTPLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "kafka_gateway_http_response_time_seconds",
			Help:    "HTTP response time in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"endpoint", "method", "status_code"},
	)
)
