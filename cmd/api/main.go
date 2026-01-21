package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"

	"InteractiveScraper/internal/api/routes"
	"InteractiveScraper/internal/config"
	"InteractiveScraper/internal/storage"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	slog.SetDefault(logger)

	logger.Info("CTI API server başlatılıyor")

	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatal("config doğrulaması başarısız oldu:", err)
	}

	store, err := storage.NewPostgresStorage(cfg.GetDSN(), logger)
	if err != nil {
		log.Fatal("depolama başlatılamadı:", err)
	}
	defer store.Close()

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	routes.SetupRoutes(router, store, logger)

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("API sunucusu dinleniyor", slog.String("port", port))

	if err := router.Run(":" + port); err != nil {
		log.Fatal("sunucu başarısız oldu:", err)
	}
}
