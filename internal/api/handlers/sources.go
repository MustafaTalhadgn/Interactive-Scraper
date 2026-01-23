package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"InteractiveScraper/internal/api/models"
	"InteractiveScraper/internal/storage"
)

type SourceHandler struct {
	storage storage.Storage
	logger  *slog.Logger
}

func NewSourceHandler(storage storage.Storage, logger *slog.Logger) *SourceHandler {
	return &SourceHandler{
		storage: storage,
		logger:  logger,
	}
}

func (h *SourceHandler) ListSources(c *gin.Context) {
	ctx := c.Request.Context()

	sources, err := h.storage.GetAllSources(ctx)
	if err != nil {
		h.logger.Error("failed to list sources", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, models.StandardResponse{
			Success: false,
			Error: &models.ErrorInfo{
				Code:    "DATABASE_ERROR",
				Message: "Failed to fetch sources",
			},
		})
		return
	}

	items := make([]models.SourceItem, len(sources))
	for i, src := range sources {
		items[i] = models.SourceItem{
			ID:             src.ID,
			Name:           src.Name,
			URL:            src.URL,
			Category:       src.Category,
			Criticality:    src.Criticality,
			Enabled:        src.Enabled,
			ScrapeInterval: src.ScrapeInterval,
			LastScrapedAt:  src.LastScrapedAt,
			CreatedAt:      src.CreatedAt,
		}
	}

	response := models.SourcesResponse{
		Sources: items,
	}

	c.JSON(http.StatusOK, models.StandardResponse{
		Success: true,
		Data:    response,
	})
}

func (h *SourceHandler) CreateSource(c *gin.Context) {
	ctx := c.Request.Context()

	var req models.SourceCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.StandardResponse{
			Success: false,
			Error: &models.ErrorInfo{
				Code:    "VALIDATION_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	if req.ScrapeInterval == "" {
		req.ScrapeInterval = "1 hour"
	}

	source, err := h.storage.CreateSource(ctx, &storage.SourceCreateInput{
		Name:           req.Name,
		URL:            req.URL,
		Category:       req.Category,
		Criticality:    req.Criticality,
		ScrapeInterval: req.ScrapeInterval,
	})

	if err != nil {
		h.logger.Error("failed to create source", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, models.StandardResponse{
			Success: false,
			Error: &models.ErrorInfo{
				Code:    "CREATE_FAILED",
				Message: "Failed to create source",
			},
		})
		return
	}

	response := models.SourceItem{
		ID:             source.ID,
		Name:           source.Name,
		URL:            source.URL,
		Category:       source.Category,
		Criticality:    source.Criticality,
		Enabled:        source.Enabled,
		ScrapeInterval: source.ScrapeInterval,
		CreatedAt:      source.CreatedAt,
	}

	c.JSON(http.StatusCreated, models.StandardResponse{
		Success: true,
		Data:    response,
	})
}

func (h *SourceHandler) UpdateSource(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.StandardResponse{
			Success: false,
			Error: &models.ErrorInfo{
				Code:    "INVALID_ID",
				Message: "Invalid source ID",
			},
		})
		return
	}

	var req models.SourceUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.StandardResponse{
			Success: false,
			Error: &models.ErrorInfo{
				Code:    "VALIDATION_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	source, err := h.storage.UpdateSource(ctx, id, &storage.SourceUpdateInput{
		Name:           req.Name,
		Criticality:    req.Criticality,
		Enabled:        req.Enabled,
		ScrapeInterval: req.ScrapeInterval,
	})

	if err != nil {
		h.logger.Error("failed to update source",
			slog.Int("id", id),
			slog.String("error", err.Error()),
		)
		c.JSON(http.StatusInternalServerError, models.StandardResponse{
			Success: false,
			Error: &models.ErrorInfo{
				Code:    "UPDATE_FAILED",
				Message: "Failed to update source",
			},
		})
		return
	}

	response := models.SourceItem{
		ID:             source.ID,
		Name:           source.Name,
		URL:            source.URL,
		Category:       source.Category,
		Criticality:    source.Criticality,
		Enabled:        source.Enabled,
		ScrapeInterval: source.ScrapeInterval,
		LastScrapedAt:  source.LastScrapedAt,
		CreatedAt:      source.CreatedAt,
	}

	c.JSON(http.StatusOK, models.StandardResponse{
		Success: true,
		Data:    response,
	})
}

func (h *SourceHandler) DeleteSource(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.StandardResponse{
			Success: false,
			Error: &models.ErrorInfo{
				Code:    "INVALID_ID",
				Message: "Invalid source ID",
			},
		})
		return
	}

	if err := h.storage.DeleteSource(ctx, id); err != nil {
		h.logger.Error("failed to delete source",
			slog.Int("id", id),
			slog.String("error", err.Error()),
		)
		c.JSON(http.StatusInternalServerError, models.StandardResponse{
			Success: false,
			Error: &models.ErrorInfo{
				Code:    "DELETE_FAILED",
				Message: "Failed to delete source",
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.StandardResponse{
		Success: true,
		Data:    map[string]interface{}{"deleted": true},
	})
}

func (h *SourceHandler) TriggerScrape(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.StandardResponse{
			Success: false,
			Error: &models.ErrorInfo{
				Code:    "INVALID_ID",
				Message: "Invalid source ID",
			},
		})
		return

	}

	go func() {
		if err := h.storage.TriggerManualScrape(context.Background(), id); err != nil {
			h.logger.Error("manual scrape failed",
				slog.Int("source_id", id),
				slog.String("error", err.Error()),
			)
		}
	}()

	c.JSON(http.StatusAccepted, models.StandardResponse{
		Success: true,
		Data:    map[string]interface{}{"message": "Scrape triggered"},
	})
}
