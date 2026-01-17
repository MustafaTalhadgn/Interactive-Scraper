package scraper

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

type HealthChecker struct {
	scraper *Scraper
	logger  *slog.Logger
}

func NewHealthChecker(scraper *Scraper, logger *slog.Logger) *HealthChecker {
	return &HealthChecker{
		scraper: scraper,
		logger:  logger,
	}
}

type HealthStatus struct {
	Status       string     `json:"status"`
	Timestamp    time.Time  `json:"timestamp"`
	DatabaseOK   bool       `json:"database_ok"`
	LastScrapeAt *time.Time `json:"last_scrape_at,omitempty"`
	TotalScraped int        `json:"total_scraped"`
	Error        string     `json:"error,omitempty"`
}

func (h *HealthChecker) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		status := h.Check(ctx)

		w.Header().Set("Content-Type", "application/json")

		if status.Status == "healthy" {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		json.NewEncoder(w).Encode(status)
	}
}

func (h *HealthChecker) Check(ctx context.Context) *HealthStatus {
	status := &HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
	}

	if err := h.scraper.storage.Ping(ctx); err != nil {
		status.Status = "unhealthy"
		status.DatabaseOK = false
		status.Error = "database ping failed: " + err.Error()
		return status
	}
	status.DatabaseOK = true

	stats := h.scraper.GetStats()
	status.LastScrapeAt = stats.LastScrapeAt
	status.TotalScraped = stats.TotalScraped

	return status
}
