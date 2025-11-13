package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthMiddleware struct {
	APIKeys []string
	Logger  *zap.Logger
}

func NewAuthMiddleware(apiKeys []string, logger *zap.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		APIKeys: apiKeys,
		Logger:  logger,
	}
}

func (am *AuthMiddleware) AuthRequired(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		am.Logger.Info("Missing authorization header")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization header"})
		c.Abort()
		return
	}

	// Ожидаем формат "Bearer <api-key>" или просто "<api-key>"
	var apiKey string
	if strings.HasPrefix(authHeader, "Bearer ") {
		apiKey = strings.TrimPrefix(authHeader, "Bearer ")
	} else {
		apiKey = authHeader
	}

	// Проверяем, есть ли API ключ в списке разрешенных
	isValid := false
	for _, key := range am.APIKeys {
		if apiKey == key {
			isValid = true
			break
		}
	}

	if !isValid {
		am.Logger.Info("Invalid API key provided", zap.String("api_key", apiKey))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
		c.Abort()
		return
	}

	// Добавляем информацию о ключе в контекст
	c.Set("api_key", apiKey)
	c.Next()
}
