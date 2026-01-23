package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"InteractiveScraper/internal/api/models"
	"InteractiveScraper/internal/storage"
)

type AuthHandler struct {
	storage   storage.Storage
	logger    *slog.Logger
	jwtSecret []byte
}

func NewAuthHandler(storage storage.Storage, logger *slog.Logger, jwtSecret string) *AuthHandler {

	if jwtSecret == "" {
		jwtSecret = "default-secret-key-change-in-production"
		logger.Warn("JWT Secret boş geldi, varsayılan anahtar kullanılıyor!")
	}

	return &AuthHandler{
		storage:   storage,
		logger:    logger,
		jwtSecret: []byte(jwtSecret),
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.StandardResponse{
			Success: false,
			Error:   &models.ErrorInfo{Code: "VALIDATION_ERROR", Message: "Eksik bilgi"},
		})
		return
	}

	user, err := h.storage.GetUserByUsername(c.Request.Context(), req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.StandardResponse{
			Success: false,
			Error:   &models.ErrorInfo{Code: "AUTH_FAILED", Message: "Kullanıcı adı veya şifre hatalı"},
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, models.StandardResponse{
			Success: false,
			Error:   &models.ErrorInfo{Code: "AUTH_FAILED", Message: "Kullanıcı adı veya şifre hatalı"},
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // 24 Saat geçerli
	})

	tokenString, err := token.SignedString(h.jwtSecret)
	if err != nil {
		h.logger.Error("Token oluşturma hatası", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, models.StandardResponse{
			Success: false,
			Error:   &models.ErrorInfo{Code: "TOKEN_ERROR", Message: "Token oluşturulamadı"},
		})
		return
	}

	c.JSON(http.StatusOK, models.StandardResponse{
		Success: true,
		Data: models.LoginResponse{
			Token: tokenString,
			User:  *user,
		},
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.StandardResponse{
			Success: false,
			Error:   &models.ErrorInfo{Code: "VALIDATION_ERROR", Message: "Eksik bilgi"},
		})
		return
	}

	existingUser, _ := h.storage.GetUserByUsername(c.Request.Context(), req.Username)
	if existingUser != nil {
		c.JSON(http.StatusConflict, models.StandardResponse{
			Success: false,
			Error:   &models.ErrorInfo{Code: "USER_EXISTS", Message: "Bu kullanıcı adı zaten alınmış"},
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.StandardResponse{
			Success: false,
			Error:   &models.ErrorInfo{Code: "SERVER_ERROR", Message: "Şifre işlenirken hata oluştu"},
		})
		return
	}

	newUser := &models.User{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		Role:         "analyst",
	}

	if err := h.storage.CreateUser(c.Request.Context(), newUser); err != nil {
		h.logger.Error("Kullanıcı oluşturulamadı", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, models.StandardResponse{
			Success: false,
			Error:   &models.ErrorInfo{Code: "DB_ERROR", Message: "Kayıt işlemi başarısız"},
		})
		return
	}

	h.logger.Info("Yeni kullanıcı kayıt oldu", slog.String("username", newUser.Username))

	c.JSON(http.StatusCreated, models.StandardResponse{
		Success: true,
		Data:    map[string]string{"message": "Kayıt başarılı! Giriş yapabilirsiniz."},
	})
}

func (h *AuthHandler) SeedAdmin() {
	ctx := context.Background()
	_, err := h.storage.GetUserByUsername(ctx, "admin")

	if err != nil {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)

		adminUser := &models.User{
			Username:     "admin",
			PasswordHash: string(hashedPassword),
			Role:         "admin",
		}

		if err := h.storage.CreateUser(ctx, adminUser); err != nil {
			h.logger.Error("Default admin oluşturulamadı", slog.String("error", err.Error()))
		} else {
			h.logger.Info("Default admin user created (admin / 123456)")
		}
	}
}
