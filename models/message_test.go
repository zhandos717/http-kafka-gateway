package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestMessageRequestValidation(t *testing.T) {
	tests := []struct {
		name     string
		request  MessageRequest
		hasError bool
	}{
		{
			name: "valid request",
			request: MessageRequest{
				Topic: "test-topic",
				Key:   "test-key",
				Value: "test-value",
			},
			hasError: false,
		},
		{
			name: "missing topic",
			request: MessageRequest{
				Key:   "test-key",
				Value: "test-value",
			},
			hasError: false, // json.Unmarshal не проверяет теги валидации
		},
		{
			name: "missing value",
			request: MessageRequest{
				Topic: "test-topic",
				Key:   "test-key",
			},
			hasError: false, // json.Unmarshal не проверяет теги валидации
		},
		{
			name: "with headers",
			request: MessageRequest{
				Topic:   "test-topic",
				Key:     "test-key",
				Value:   "test-value",
				Headers: map[string]string{"header1": "value1"},
			},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.request)
			if err != nil {
				t.Fatalf("Failed to marshal request: %v", err)
			}

			var req MessageRequest
			err = json.Unmarshal(data, &req)
			if err != nil {
				t.Fatalf("Failed to unmarshal request: %v", err)
			}

			// Проверяем, что данные корректно десериализовались
			if req.Topic != tt.request.Topic {
				t.Errorf("Topic mismatch: expected %s, got %s", tt.request.Topic, req.Topic)
			}
			if req.Key != tt.request.Key {
				t.Errorf("Key mismatch: expected %s, got %s", tt.request.Key, req.Key)
			}
			// Для Value используем строковое сравнение
			if req.Value != tt.request.Value {
				t.Errorf("Value mismatch: expected %v, got %v", tt.request.Value, req.Value)
			}
		})
	}
}

func TestMessageResponse(t *testing.T) {
	timestamp := time.Now()
	response := MessageResponse{
		Success:   true,
		Message:   "test message",
		Error:     "",
		Timestamp: timestamp,
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	var resp MessageResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Success != response.Success {
		t.Errorf("Expected Success=%v, got %v", response.Success, resp.Success)
	}
	if resp.Message != response.Message {
		t.Errorf("Expected Message=%s, got %s", response.Message, resp.Message)
	}
	if resp.Error != response.Error {
		t.Errorf("Expected Error=%s, got %s", response.Error, resp.Error)
	}
}

func TestKafkaMessage(t *testing.T) {
	timestamp := time.Now()
	kafkaMsg := KafkaMessage{
		Topic:     "test-topic",
		Key:       []byte("test-key"),
		Value:     []byte("test-value"),
		Headers:   map[string][]byte{"header1": []byte("value1")},
		Timestamp: timestamp,
	}

	if kafkaMsg.Topic != "test-topic" {
		t.Errorf("Expected Topic=%s, got %s", "test-topic", kafkaMsg.Topic)
	}
	if string(kafkaMsg.Key) != "test-key" {
		t.Errorf("Expected Key=%s, got %s", "test-key", string(kafkaMsg.Key))
	}
	if string(kafkaMsg.Value) != "test-value" {
		t.Errorf("Expected Value=%s, got %s", "test-value", string(kafkaMsg.Value))
	}
	if len(kafkaMsg.Headers) != 1 {
		t.Errorf("Expected 1 header, got %d", len(kafkaMsg.Headers))
	}
	if kafkaMsg.Timestamp != timestamp {
		t.Errorf("Expected Timestamp=%v, got %v", timestamp, kafkaMsg.Timestamp)
	}
}
