package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Сохраняем оригинальные значения
	originalKafkaBrokers := os.Getenv("KAFKA_BROKERS")
	originalKafkaLogLevel := os.Getenv("KAFKA_LOG_LEVEL")
	originalServerPort := os.Getenv("SERVER_PORT")
	originalAPIKeys := os.Getenv("API_KEYS")

	// Устанавливаем тестовые значения
	os.Setenv("KAFKA_BROKERS", "test-broker:9092")
	os.Setenv("KAFKA_LOG_LEVEL", "3")
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("API_KEYS", "test-api-key")

	// Загружаем конфигурацию
	config := LoadConfig()

	// Проверяем значения
	if config.KafkaBrokers != "test-broker:9092" {
		t.Errorf("Expected KafkaBrokers=test-broker:9092, got %s", config.KafkaBrokers)
	}

	if config.KafkaLogLevel != 3 {
		t.Errorf("Expected KafkaLogLevel=3, got %d", config.KafkaLogLevel)
	}

	if config.ServerPort != "9090" {
		t.Errorf("Expected ServerPort=9090, got %s", config.ServerPort)
	}

	if len(config.APIKeys) != 1 || config.APIKeys[0] != "test-api-key" {
		t.Errorf("Expected APIKeys=[test-api-key], got %v", config.APIKeys)
	}

	if config.Logger == nil {
		t.Errorf("Expected logger to be initialized")
	}

	// Восстанавливаем оригинальные значения
	os.Setenv("KAFKA_BROKERS", originalKafkaBrokers)
	os.Setenv("KAFKA_LOG_LEVEL", originalKafkaLogLevel)
	os.Setenv("SERVER_PORT", originalServerPort)
	os.Setenv("API_KEYS", originalAPIKeys)
}

func TestLoadConfigWithDefaults(t *testing.T) {
	// Удаляем переменные окружения для тестирования значений по умолчанию
	os.Unsetenv("KAFKA_BROKERS")
	os.Unsetenv("KAFKA_LOG_LEVEL")
	os.Unsetenv("SERVER_PORT")
	os.Unsetenv("API_KEYS")

	config := LoadConfig()

	// Проверяем значения по умолчанию
	if config.KafkaBrokers != "localhost:9092" {
		t.Errorf("Expected default KafkaBrokers=localhost:9092, got %s", config.KafkaBrokers)
	}

	if config.KafkaLogLevel != 1 {
		t.Errorf("Expected default KafkaLogLevel=1, got %d", config.KafkaLogLevel)
	}

	if config.ServerPort != "8080" {
		t.Errorf("Expected default ServerPort=8080, got %s", config.ServerPort)
	}

	if len(config.APIKeys) != 1 || config.APIKeys[0] != "default-api-key" {
		t.Errorf("Expected default APIKeys=[default-api-key], got %v", config.APIKeys)
	}
}
