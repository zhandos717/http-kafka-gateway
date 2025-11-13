package metrics

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

// MetricsResponse структура для ответа API метрик
type MetricsResponse struct {
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

// В реальном приложении здесь будет глобальная переменная или структура для хранения данных
// Но для демонстрационных целей мы создадим фиктивные данные
var (
	// topicsList список топиков
	topicsList = []TopicInfo{
		{Name: "default-topic", MessageCount: 100},
		{Name: "user-events", MessageCount: 50},
		{Name: "system-logs", MessageCount: 75},
		{Name: "notifications", MessageCount: 25},
		{Name: "analytics", MessageCount: 120},
	}

	// recentMessages список недавних сообщений
	recentMessages = []MessageInfo{
		{Topic: "default-topic", Timestamp: time.Now().Add(-1 * time.Minute), Status: "success"},
		{Topic: "user-events", Timestamp: time.Now().Add(-2 * time.Minute), Status: "success"},
		{Topic: "system-logs", Timestamp: time.Now().Add(-3 * time.Minute), Status: "error"},
		{Topic: "default-topic", Timestamp: time.Now().Add(-4 * time.Minute), Status: "success"},
		{Topic: "notifications", Timestamp: time.Now().Add(-5 * time.Minute), Status: "success"},
	}
)

// UpdateDemoData обновляет демонстрационные данные для симуляции активности
func UpdateDemoData() {
	// Обновляем счётчики сообщений в топиках
	for i := range topicsList {
		// Случайное изменение счётчика
		topicsList[i].MessageCount += int64(1 + time.Now().Nanosecond()%5)
	}

	// Добавляем новое сообщение и удаляем самое старое, чтобы поддерживать размер
	newMessage := MessageInfo{
		Topic:     topicsList[time.Now().Nanosecond()%len(topicsList)].Name,
		Timestamp: time.Now(),
		Status:    []string{"success", "error"}[time.Now().Nanosecond()%2],
	}

	// Добавляем новое сообщение в начало
	recentMessages = append([]MessageInfo{newMessage}, recentMessages...)

	// Ограничиваем размер списка до 10 элементов
	if len(recentMessages) > 10 {
		recentMessages = recentMessages[:10]
	}
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

	c.JSON(http.StatusOK, MetricsResponse{
		Timestamp: time.Now().Unix(),
		Metrics:   result,
	})
}

// GetTopics возвращает список активных топиков
func GetTopics(c *gin.Context) {
	// В реальном приложении здесь будет логика получения списка топиков из Kafka
	// или из внутреннего состояния приложения
	// Для демонстрации используем фиктивные данные

	c.JSON(http.StatusOK, gin.H{
		"topics":    topicsList,
		"timestamp": time.Now().Unix(),
	})
}

// GetRecentMessages возвращает список недавних сообщений
func GetRecentMessages(c *gin.Context) {
	// В реальном приложении здесь будет логика получения недавних сообщений
	// из внутреннего лога или другого источника
	// Для демонстрации используем фиктивные данные

	c.JSON(http.StatusOK, gin.H{
		"messages":  recentMessages,
		"timestamp": time.Now().Unix(),
	})
}
