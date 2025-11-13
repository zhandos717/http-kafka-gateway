package metrics

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

// Response структура для ответа API метрик
type Response struct {
	Timestamp int64       `json:"timestamp"`
	Metrics   interface{} `json:"metrics"`
}

// TopicInfo информация о топике
type TopicInfo struct {
	Name         string `json:"name"`
	MessageCount int64  `json:"messageCount"`
}

// MessageInfo информация о сообщении
type MessageInfo struct {
	Topic     string    `json:"topic"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`
}

// MetricsStore структура для хранения информации о топиках и сообщениях
type MetricsStore struct {
	mu sync.RWMutex
}

// UpdateDemoData обновляет демонстрационные данные для симуляции активности
func UpdateDemoData() {
	// В реальном приложении здесь будет логика обновления данных из Prometheus метрик
	// Пока оставляем пустой, так как данные будут получаться напрямую из Prometheus
}

// GetMetrics возвращает метрики в формате JSON для UI
func GetMetrics(c *gin.Context) {
	// Получаем все метрики из реестра
	metrics, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to gather metrics",
		})
		return
	}

	// Формируем структуру данных для UI
	result := make(map[string]interface{})

	for _, metricFamily := range metrics {
		metricName := metricFamily.GetName()

		// Обрабатываем только нужные метрики
		switch metricName {
		case "kafka_gateway_messages_processed_total":
			successCount := 0.0
			errorCount := 0.0
			for _, metric := range metricFamily.GetMetric() {
				labels := metric.GetLabel()
				for _, label := range labels {
					if label.GetName() == "status" {
						if label.GetValue() == "success" {
							successCount = metric.GetCounter().GetValue()
						} else if label.GetValue() == "error" {
							errorCount = metric.GetCounter().GetValue()
						}
						break
					}
				}
			}
			result["messages_processed"] = map[string]float64{
				"success": successCount,
				"error":   errorCount,
			}

		case "kafka_gateway_kafka_errors_total":
			errorCount := 0.0
			for _, metric := range metricFamily.GetMetric() {
				errorCount += metric.GetCounter().GetValue()
			}
			result["kafka_errors"] = errorCount

		case "kafka_gateway_auth_attempts_total":
			successCount := 0.0
			errorCount := 0.0
			for _, metric := range metricFamily.GetMetric() {
				labels := metric.GetLabel()
				for _, label := range labels {
					if label.GetName() == "status" {
						if label.GetValue() == "success" {
							successCount = metric.GetCounter().GetValue()
						} else if label.GetValue() == "failed" {
							errorCount = metric.GetCounter().GetValue()
						}
						break
					}
				}
			}
			result["auth_attempts"] = map[string]float64{
				"success": successCount,
				"error":   errorCount,
			}

		case "kafka_gateway_request_duration_seconds_sum",
			"kafka_gateway_http_response_time_seconds_sum":
			// Здесь можно обработать суммарное время запросов
			totalTime := 0.0
			for _, metric := range metricFamily.GetMetric() {
				totalTime += metric.GetCounter().GetValue()
			}
			result["request_duration_sum"] = totalTime
		}
	}

	// Добавляем дополнительные вычисляемые метрики
	if messagesProcessed, ok := result["messages_processed"].(map[string]float64); ok {
		totalMessages := messagesProcessed["success"] + messagesProcessed["error"]
		result["total_messages"] = totalMessages
	}

	c.JSON(http.StatusOK, Response{
		Timestamp: time.Now().Unix(),
		Metrics:   result,
	})
}

// GetTopics возвращает список активных топиков
func GetTopics(c *gin.Context) {
	// Получаем метрики из Prometheus и извлекаем информацию о топиках
	metrics, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to gather metrics",
		})
		return
	}

	// Собираем информацию о топиках из метрик
	topicsMap := make(map[string]int64)
	for _, metricFamily := range metrics {
		metricName := metricFamily.GetName()
		if metricName == "kafka_gateway_messages_processed_total" {
			for _, metric := range metricFamily.GetMetric() {
				labels := metric.GetLabel()
				topic := ""
				for _, label := range labels {
					if label.GetName() == "topic" {
						topic = label.GetValue()
						break
					}
				}
				if topic != "" {
					count := int64(metric.GetCounter().GetValue())
					if count > topicsMap[topic] {
						topicsMap[topic] = count
					}
				}
			}
		}
	}

	// Преобразуем map в slice
	var topicsList []TopicInfo
	for topic, count := range topicsMap {
		topicsList = append(topicsList, TopicInfo{
			Name:         topic,
			MessageCount: count,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"topics":    topicsList,
		"timestamp": time.Now().Unix(),
	})
}

// GetRecentMessages возвращает список недавних сообщений
func GetRecentMessages(c *gin.Context) {
	// В реальном приложении здесь будет логика получения недавних сообщений
	// из внутреннего лога или другого источника
	// Пока возвращаем пустой список, так как нам нужно реализовать логику хранения недавних сообщений
	// или получения их из другого источника

	// Для демонстрации возвращаем последние успешные/ошибочные сообщения из метрик
	metrics, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to gather metrics",
		})
		return
	}

	var recentMessages []MessageInfo
	for _, metricFamily := range metrics {
		metricName := metricFamily.GetName()
		if metricName == "kafka_gateway_messages_processed_total" {
			for _, metric := range metricFamily.GetMetric() {
				labels := metric.GetLabel()
				topic := ""
				status := ""
				for _, label := range labels {
					if label.GetName() == "topic" {
						topic = label.GetValue()
					} else if label.GetName() == "status" {
						status = label.GetValue()
					}
				}
				if topic != "" && status != "" {
					count := metric.GetCounter().GetValue()
					if count > 0 { // Только если есть сообщения с таким статусом
						recentMessages = append(recentMessages, MessageInfo{
							Topic:     topic,
							Timestamp: time.Now(), // В реальности нужно использовать реальное время
							Status:    status,
						})
					}
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"messages":  recentMessages,
		"timestamp": time.Now().Unix(),
	})
}
