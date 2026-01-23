package middleware

import (
	"net/http"
	"strings"

	"InteractiveScraper/internal/api/models"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("super-secret-key-change-in-production")

func AuthMiddleware() gin.HandlerFunc {
	secretStr := os.Getenv("JWT_SECRET")
	if secretStr == "" {
		secretStr = "super-secret-key-change-in-production"
	}
	jwtSecret := []byte(secretStr)

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.StandardResponse{Success: false, Error: &models.ErrorInfo{Code: "UNAUTHORIZED", Message: "Token gerekli"}})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.StandardResponse{Success: false, Error: &models.ErrorInfo{Code: "UNAUTHORIZED", Message: "Geçersiz token formatı"}})
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.StandardResponse{Success: false, Error: &models.ErrorInfo{Code: "UNAUTHORIZED", Message: "Geçersiz veya süresi dolmuş token"}})
			return
		}

		c.Next()
	}
}
