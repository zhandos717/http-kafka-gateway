package logger

import (
	"testing"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name     string
		level    int
		hasError bool
	}{
		{
			name:     "debug level",
			level:    4,
			hasError: false,
		},
		{
			name:     "info level",
			level:    3,
			hasError: false,
		},
		{
			name:     "warn level",
			level:    2,
			hasError: false,
		},
		{
			name:     "error level",
			level:    1,
			hasError: false,
		},
		{
			name:     "fatal level",
			level:    0,
			hasError: false,
		},
		{
			name:     "invalid level",
			level:    10,
			hasError: false, // Should default to info level
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := NewLogger(tt.level)
			if tt.hasError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.hasError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if logger == nil && !tt.hasError {
				t.Errorf("Expected logger but got nil")
			}
			if logger != nil {
				// Test that logger can be used without panic
				logger.Info("Test message")
				logger.Sync()
			}
		})
	}
}
