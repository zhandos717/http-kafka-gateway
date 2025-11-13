package utils

import (
	"testing"
)

func TestConvertInterfaceToBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected []byte
		hasError bool
	}{
		{
			name:     "string input",
			input:    "test string",
			expected: []byte("test string"),
			hasError: false,
		},
		{
			name:     "[]byte input",
			input:    []byte("test bytes"),
			expected: []byte("test bytes"),
			hasError: false,
		},
		{
			name:     "int input",
			input:    42,
			expected: []byte("42"),
			hasError: false,
		},
		{
			name:     "struct input",
			input:    struct{ Name string }{Name: "test"},
			expected: []byte(`{"Name":"test"}`),
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ConvertInterfaceToBytes(tt.input)
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if string(result) != string(tt.expected) {
				t.Errorf("Expected %s, got %s", string(tt.expected), string(result))
			}
		})
	}
}

func TestIsValidTopic(t *testing.T) {
	tests := []struct {
		name     string
		topic    string
		expected bool
	}{
		{
			name:     "valid topic",
			topic:    "valid-topic",
			expected: true,
		},
		{
			name:     "valid topic with underscores",
			topic:    "valid_topic",
			expected: true,
		},
		{
			name:     "valid topic with dots",
			topic:    "valid.topic",
			expected: true,
		},
		{
			name:     "topic starting with dot",
			topic:    ".invalid-topic",
			expected: false,
		},
		{
			name:     "topic starting with underscore",
			topic:    "_invalid-topic",
			expected: false,
		},
		{
			name:     "topic with invalid characters",
			topic:    "invalid@topic",
			expected: false,
		},
		{
			name:     "empty topic",
			topic:    "",
			expected: false,
		},
		{
			name:     "topic too long",
			topic:    "a23456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901......",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidTopic(tt.topic)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
