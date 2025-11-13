package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func TestNewAuthMiddleware(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	apiKeys := []string{"test-key-1", "test-key-2"}
	authMiddleware := NewAuthMiddleware(apiKeys, logger)

	if len(authMiddleware.APIKeys) != 2 {
		t.Errorf("Expected 2 API keys, got %d", len(authMiddleware.APIKeys))
	}

	if authMiddleware.Logger == nil {
		t.Errorf("Expected logger to be set")
	}
}

func TestAuthMiddleware_AuthRequired(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	apiKeys := []string{"valid-key"}
	authMiddleware := NewAuthMiddleware(apiKeys, logger)

	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
	}{
		{
			name:           "valid Bearer token",
			authHeader:     "Bearer valid-key",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "valid token without Bearer prefix",
			authHeader:     "valid-key",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid token",
			authHeader:     "Bearer invalid-key",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "missing header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid Bearer token",
			authHeader:     "Bearer another-invalid-key",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			if tt.authHeader != "" {
				c.Request, _ = http.NewRequest("POST", "/test", nil)
				c.Request.Header.Set("Authorization", tt.authHeader)
			} else {
				c.Request, _ = http.NewRequest("POST", "/test", nil)
			}

			authMiddleware.AuthRequired(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// For successful auth, check if api_key is set in context
			if tt.expectedStatus == http.StatusOK {
				apiKey, exists := c.Get("api_key")
				if !exists {
					t.Errorf("Expected api_key to be set in context")
				}
				if apiKey != "valid-key" {
					t.Errorf("Expected api_key to be 'valid-key', got '%s'", apiKey)
				}
			}
		})
	}
}
