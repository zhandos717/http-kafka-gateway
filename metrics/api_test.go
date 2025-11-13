package metrics

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Устанавливаем Gin в режим тестирования
	gin.SetMode(gin.TestMode)
}

func TestGetMetrics(t *testing.T) {
	// Создаем HTTP recorder для записи ответа
	w := httptest.NewRecorder()

	// Создаем фейковый запрос
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/metrics", nil)

	// Вызываем функцию
	GetMetrics(c)

	// Проверяем статус
	assert.Equal(t, http.StatusOK, w.Code)

	// Проверяем, что ответ не пустой
	assert.NotEmpty(t, w.Body.String())
}

func TestGetTopics(t *testing.T) {
	// Создаем HTTP recorder для записи ответа
	w := httptest.NewRecorder()

	// Создаем фейковый запрос
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/topics", nil)

	// Вызываем функцию
	GetTopics(c)

	// Проверяем статус
	assert.Equal(t, http.StatusOK, w.Code)

	// Проверяем, что ответ не пустой
	assert.NotEmpty(t, w.Body.String())
}

func TestGetRecentMessages(t *testing.T) {
	// Создаем HTTP recorder для записи ответа
	w := httptest.NewRecorder()

	// Создаем фейковый запрос
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/messages", nil)

	// Вызываем функцию
	GetRecentMessages(c)

	// Проверяем статус
	assert.Equal(t, http.StatusOK, w.Code)

	// Проверяем, что ответ не пустой
	assert.NotEmpty(t, w.Body.String())
}

func TestUpdateDemoData(t *testing.T) {
	// Сохраняем начальные значения
	initialTopics := len(store.topicsList)

	// Вызываем функцию обновления данных
	UpdateDemoData()

	// Проверяем, что количество топиков не изменилось
	assert.Equal(t, initialTopics, len(store.topicsList))

	// Проверяем, что количество сообщений не превышает 10
	assert.True(t, len(store.recentMessages) <= 10)

	// Проверяем, что хотя бы одно сообщение есть
	assert.True(t, len(store.recentMessages) >= 1)
}

func TestMetricsStoreConcurrency(t *testing.T) {
	// Тест проверяет, что хранилище метрик безопасно для конкурентного доступа
	done := make(chan bool)

	// Запускаем несколько горутин для чтения
	for i := 0; i < 5; i++ {
		go func() {
			for j := 0; j < 10; j++ {
				// Создаем фейковый контекст для безопасного вызова
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request = httptest.NewRequest("GET", "/api/topics", nil)
				GetTopics(c)

				w2 := httptest.NewRecorder()
				c2, _ := gin.CreateTestContext(w2)
				c2.Request = httptest.NewRequest("GET", "/api/messages", nil)
				GetRecentMessages(c2)
			}
			done <- true
		}()
	}

	// Запускаем несколько горутин для записи
	for i := 0; i < 3; i++ {
		go func() {
			for j := 0; j < 10; j++ {
				UpdateDemoData()
				time.Sleep(1 * time.Millisecond) // Небольшая задержка для лучшего чередования
			}
			done <- true
		}()
	}

	// Ждем завершения всех горутин
	for i := 0; i < 8; i++ {
		<-done
	}
}

func TestMetricsRegistry(t *testing.T) {
	// Проверяем, что метрики существуют в пакете
	assert.NotNil(t, MessagesProcessed)
	assert.NotNil(t, RequestDuration)
	assert.NotNil(t, KafkaErrors)
	assert.NotNil(t, AuthAttempts)
	assert.NotNil(t, HTTPLatency)
}
