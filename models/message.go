package models

import "time"

type MessageRequest struct {
	Topic   string            `json:"topic" binding:"required"`
	Key     string            `json:"key,omitempty"`
	Value   interface{}       `json:"value" binding:"required"`
	Headers map[string]string `json:"headers,omitempty"`
}

type MessageResponse struct {
	Success   bool      `json:"success"`
	Message   string    `json:"message,omitempty"`
	Error     string    `json:"error,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

type KafkaMessage struct {
	Topic     string
	Key       []byte
	Value     []byte
	Headers   map[string][]byte
	Timestamp time.Time
}
