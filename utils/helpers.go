package utils

import (
	"encoding/json"
	"time"
)

// ConvertInterfaceToBytes конвертирует interface{} в []byte
func ConvertInterfaceToBytes(data interface{}) ([]byte, error) {
	switch v := data.(type) {
	case string:
		return []byte(v), nil
	case []byte:
		return v, nil
	default:
		return json.Marshal(data)
	}
}

// GetCurrentTime возвращает текущее время
func GetCurrentTime() time.Time {
	return time.Now()
}

// IsValidTopic проверяет валидность имени топика Kafka
func IsValidTopic(topic string) bool {
	if len(topic) == 0 || len(topic) > 249 {
		return false
	}

	// Проверяем, что топик соответствует требованиям Kafka
	// - Только буквы, цифры, точки, подчеркивания, тире
	// - Не начинается с точки или подчеркивания
	for i, r := range topic {
		if i == 0 && (r == '.' || r == '_') {
			return false
		}
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == '.' || r == '_' || r == '-') {
			return false
		}
	}

	return true
}
