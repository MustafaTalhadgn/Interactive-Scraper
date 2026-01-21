package handlers

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"InteractiveScraper/internal/api/models"
	"InteractiveScraper/internal/storage"
)

type StatsHandler struct {
	storage storage.Storage
	logger  *slog.Logger
}

func NewStatsHandler(storage storage.Storage, logger *slog.Logger) *StatsHandler {
	return &StatsHandler{
		storage: storage,
		logger:  logger,
	}
}

func (h *StatsHandler) GetOverview(c *gin.Context) {
	ctx := c.Request.Context()

	total, err := h.storage.GetTotalCount(ctx)
	if err != nil {
		h.logger.Error("failed to get total count", slog.String("error", err.Error()))
		total = 0
	}

	criticalCount, _ := h.storage.GetIntelligenceCountByCriticality(ctx, 76)

	last24h, _ := h.storage.GetIntelligenceCountSince(ctx, time.Now().Add(-24*time.Hour))

	sources, _ := h.storage.GetEnabledSources(ctx)
	activeSources := len(sources)

	criticalityDist, _ := h.storage.GetCriticalityDistribution(ctx)

	categoryDist, _ := h.storage.GetCategoryDistribution(ctx)

	response := models.StatsOverview{
		TotalIntelligence: total,
		CriticalCount:     criticalCount,
		Last24Hours:       last24h,
		ActiveSources:     activeSources,
		CriticalityDist:   criticalityDist,
		CategoryDist:      categoryDist,
	}

	c.JSON(http.StatusOK, models.StandardResponse{
		Success: true,
		Data:    response,
	})
}

func (h *StatsHandler) GetTimeline(c *gin.Context) {
	ctx := c.Request.Context()

	days := 7
	if d := c.Query("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil && parsed > 0 && parsed <= 30 {
			days = parsed
		}
	}

	timeline, err := h.storage.GetTimelineData(ctx, days)
	if err != nil {
		h.logger.Error("	", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, models.StandardResponse{
			Success: false,
			Error: &models.ErrorInfo{
				Code:    "DATABASE_ERROR",
				Message: "Failed to fetch timeline data",
			},
		})
		return
	}

	// Convert storage.TimeSeriesData to models.TimeSeriesData
	modelTimeline := make([]models.TimeSeriesData, len(timeline))
	for i, t := range timeline {
		modelTimeline[i] = models.TimeSeriesData{
			Date:     t.Date,
			Critical: t.Critical,
			High:     t.High,
			Medium:   t.Medium,
			Low:      t.Low,
		}
	}

	response := models.TimelineResponse{
		Timeline: modelTimeline,
	}

	c.JSON(http.StatusOK, models.StandardResponse{
		Success: true,
		Data:    response,
	})
}
