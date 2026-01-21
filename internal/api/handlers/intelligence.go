package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"InteractiveScraper/internal/api/models"
	"InteractiveScraper/internal/storage"
)

type IntelligenceHandler struct {
	storage storage.Storage
	logger  *slog.Logger
}

func NewIntelligenceHandler(storage storage.Storage, logger *slog.Logger) *IntelligenceHandler {
	return &IntelligenceHandler{
		storage: storage,
		logger:  logger,
	}
}

func (h *IntelligenceHandler) GetIntelligenceFeed(c *gin.Context) {
	ctx := c.Request.Context()

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	criticality := c.Query("criticality")
	sourceID, _ := strconv.Atoi(c.Query("source_id"))
	category := c.Query("category")
	search := c.Query("search")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	filters := storage.IntelligenceFilters{
		Page:        page,
		Limit:       limit,
		Criticality: criticality,
		SourceID:    sourceID,
		Category:    category,
		Search:      search,
	}

	result, err := h.storage.GetIntelligenceFeed(ctx, filters)
	if err != nil {
		h.logger.Error("failed to fetch intelligence feed", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, models.StandardResponse{
			Success: false,
			Error: &models.ErrorInfo{
				Code:    "DATABASE_ERROR",
				Message: "Failed to fetch intelligence data",
			},
		})
		return
	}

	items := make([]models.IntelligenceFeedItem, len(result.Items))
	for i, item := range result.Items {
		items[i] = models.IntelligenceFeedItem{
			ID:               item.ID,
			Title:            item.Title,
			SourceName:       item.SourceName,
			CriticalityScore: item.CriticalityScore,
			Criticality:      getCriticalityLevel(item.CriticalityScore),
			Category:         item.Category,
			CreatedAt:        item.CreatedAt,
		}
	}

	response := models.IntelligenceFeedResponse{
		Intelligence: items,
		Pagination: models.PaginationInfo{
			Page:       page,
			Limit:      limit,
			Total:      result.Total,
			TotalPages: (result.Total + limit - 1) / limit,
		},
	}

	c.JSON(http.StatusOK, models.StandardResponse{
		Success: true,
		Data:    response,
	})
}

func (h *IntelligenceHandler) GetIntelligenceDetail(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.StandardResponse{
			Success: false,
			Error: &models.ErrorInfo{
				Code:    "INVALID_ID",
				Message: "Invalid intelligence ID",
			},
		})
		return
	}

	intelligence, err := h.storage.GetIntelligenceByID(ctx, id)
	if err != nil {
		h.logger.Error("failed to fetch intelligence detail",
			slog.Int("id", id),
			slog.String("error", err.Error()),
		)
		c.JSON(http.StatusNotFound, models.StandardResponse{
			Success: false,
			Error: &models.ErrorInfo{
				Code:    "NOT_FOUND",
				Message: "Intelligence not found",
			},
		})
		return
	}

	features, err := h.storage.GetFeaturesByIntelligenceID(ctx, id)
	var extractedFeatures *models.ExtractedFeatures
	if err == nil && features != nil {
		extractedFeatures = &models.ExtractedFeatures{
			BitcoinAddrs:  features.BitcoinAddrs,
			EthereumAddrs: features.EthereumAddrs,
			MoneroAddrs:   features.MoneroAddrs,
			OnionURLs:     features.OnionURLs,
			IPAddresses:   features.IPAddresses,
			Emails:        features.Emails,
			CVEs:          features.CVEs,
			Keywords:      features.Keywords,
		}
	}

	source, _ := h.storage.GetSourceByID(ctx, intelligence.SourceID)
	sourceName := "Unknown"
	category := ""
	if source != nil {
		sourceName = source.Name
		category = source.Category
	}

	response := models.IntelligenceDetailResponse{
		ID:                intelligence.ID,
		Title:             intelligence.Title,
		Summary:           intelligence.Summary,
		SourceName:        sourceName,
		SourceURL:         intelligence.SourceURL,
		CriticalityScore:  intelligence.CriticalityScore,
		Criticality:       getCriticalityLevel(intelligence.CriticalityScore),
		Category:          category,
		PublishedAt:       intelligence.PublishedAt,
		CreatedAt:         intelligence.CreatedAt,
		ExtractedFeatures: extractedFeatures,
	}

	c.JSON(http.StatusOK, models.StandardResponse{
		Success: true,
		Data:    response,
	})
}

func getCriticalityLevel(score int) string {
	switch {
	case score >= 76:
		return "critical"
	case score >= 51:
		return "high"
	case score >= 26:
		return "medium"
	default:
		return "low"
	}
}
