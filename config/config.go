package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Config struct {
	KafkaBrokers  string
	KafkaLogLevel int
	ServerPort    string
	APIKeys       []string
	Logger        *zap.Logger
}

func LoadConfig() *Config {
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Получаем переменные окружения
	kafkaBrokers := getEnv("KAFKA_BROKERS", "localhost:9092")
	kafkaLogLevelStr := getEnv("KAFKA_LOG_LEVEL", "1")
	serverPort := getEnv("SERVER_PORT", "8080")

	// Парсим уровень логирования
	var kafkaLogLevel int
	fmt.Sscanf(kafkaLogLevelStr, "%d", &kafkaLogLevel)

	// Получаем API ключи (в реальном приложении можно загружать из безопасного хранилища)
	apiKeys := getEnv("API_KEYS", "default-api-key")

	// Создаем logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}

	return &Config{
		KafkaBrokers:  kafkaBrokers,
		KafkaLogLevel: kafkaLogLevel,
		ServerPort:    serverPort,
		APIKeys:       []string{apiKeys}, // В реальном приложении можно разделить по запятой
		Logger:        logger,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
