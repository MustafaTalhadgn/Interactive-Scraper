package routes

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"InteractiveScraper/internal/api/handlers"
	"InteractiveScraper/internal/api/middleware"
	"InteractiveScraper/internal/storage"
)

func SetupRoutes(router *gin.Engine, storage storage.Storage, logger *slog.Logger) {

	router.Use(middleware.CORS())
	router.Use(middleware.Logger(logger))
	router.Use(gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	api := router.Group("/api")

	intelligenceHandler := handlers.NewIntelligenceHandler(storage, logger)
	statsHandler := handlers.NewStatsHandler(storage, logger)
	sourceHandler := handlers.NewSourceHandler(storage, logger)

	api.GET("/intelligence", intelligenceHandler.GetIntelligenceFeed)
	api.GET("/intelligence/:id", intelligenceHandler.GetIntelligenceDetail)

	api.GET("/stats/overview", statsHandler.GetOverview)
	api.GET("/stats/timeline", statsHandler.GetTimeline)

	sources := api.Group("/sources")
	{
		sources.GET("", sourceHandler.ListSources)
		sources.POST("", sourceHandler.CreateSource)
		sources.PATCH("/:id", sourceHandler.UpdateSource)
		sources.DELETE("/:id", sourceHandler.DeleteSource)
		sources.POST("/:id/scrape", sourceHandler.TriggerScrape)
	}
}
