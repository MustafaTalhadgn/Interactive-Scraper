package routes

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"InteractiveScraper/internal/api/handlers"
	"InteractiveScraper/internal/api/middleware"
	"InteractiveScraper/internal/config"
	"InteractiveScraper/internal/storage"
)

func SetupRoutes(router *gin.Engine, storage storage.Storage, logger *slog.Logger, cfg *config.Config) {

	router.Use(middleware.CORS())
	router.Use(middleware.Logger(logger))
	router.Use(gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	api := router.Group("/api")

	authHandler := handlers.NewAuthHandler(storage, logger, cfg.JWTSecret)
	intelligenceHandler := handlers.NewIntelligenceHandler(storage, logger)
	statsHandler := handlers.NewStatsHandler(storage, logger)
	sourceHandler := handlers.NewSourceHandler(storage, logger)

	api.POST("/auth/login", authHandler.Login)
	api.POST("/auth/register", authHandler.Register)
	api.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "healthy"}) })

	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{

		protected.GET("/intelligence", intelligenceHandler.GetIntelligenceFeed)
		protected.GET("/intelligence/:id", intelligenceHandler.GetIntelligenceDetail)

		protected.GET("/stats/overview", statsHandler.GetOverview)
		protected.GET("/stats/timeline", statsHandler.GetTimeline)

		sources := protected.Group("/sources")
		{
			sources.GET("", sourceHandler.ListSources)
			sources.POST("", sourceHandler.CreateSource)
			sources.PATCH("/:id", sourceHandler.UpdateSource)
			sources.DELETE("/:id", sourceHandler.DeleteSource)
			sources.POST("/:id/scrape", sourceHandler.TriggerScrape)
		}
	}

	go authHandler.SeedAdmin()
}
