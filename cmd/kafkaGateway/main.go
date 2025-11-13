package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"kafkaGateway/config"
	"kafkaGateway/handlers"
	"kafkaGateway/kafka"
	"kafkaGateway/metrics"
	"kafkaGateway/middleware"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.LoadConfig()
	defer cfg.Logger.Sync()

	// Создаем Kafka Producer
	kafkaProducer := kafka.NewProducer(cfg.KafkaBrokers, cfg.Logger)
	defer kafkaProducer.Close()

	// Создаем обработчик сообщений
	messageHandler := handlers.NewMessageHandler(kafkaProducer, cfg.Logger)

	// Создаем middleware для аутентификации
	authMiddleware := middleware.NewAuthMiddleware(cfg.APIKeys, cfg.Logger)

	// Создаем Gin роутер
	router := gin.New()

	// Добавляем CORS middleware
	configCORS := cors.DefaultConfig()
	configCORS.AllowAllOrigins = true
	configCORS.AllowCredentials = true
	configCORS.AllowHeaders = append(configCORS.AllowHeaders, "Authorization", "Content-Type")
	router.Use(cors.New(configCORS))

	// Добавляем логирование запросов
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Маршрут для проверки состояния
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now(),
		})
	})

	// Маршрут для метрик Prometheus
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Маршрут для UI
	router.Static("/static", "./static")

	// Главная страница UI
	router.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	// API для получения метрик в формате JSON
	router.GET("/api/metrics", metrics.GetMetrics)

	// API для получения списка топиков
	router.GET("/api/topics", metrics.GetTopics)

	// API для получения недавних сообщений
	router.GET("/api/messages", metrics.GetRecentMessages)

	// Защищенные маршруты
	protected := router.Group("/")
	protected.Use(authMiddleware.AuthRequired)
	{
		protected.POST("/message", messageHandler.SendMessage)
		// Добавим новый маршрут для получения статуса
		protected.GET("/api/status", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":    "running",
				"timestamp": time.Now(),
			})
		})
	}

	// Создаем HTTP сервер
	server := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		log.Printf("Kafka Gateway starting on port %s", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Запускаем горутину для обновления демонстрационных данных
	go func() {
		ticker := time.NewTicker(5 * time.Second) // Обновляем каждые 5 секунд
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				metrics.UpdateDemoData()
			case <-context.Background().Done():
				return
			}
		}
	}()

	// Ждем сигнал остановки
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Плавное завершение
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Server exited")
}
