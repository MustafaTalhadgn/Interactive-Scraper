package scraper

import (
	"context"
	"log/slog"
	"time"

	"InteractiveScraper/internal/storage"
)

type Scheduler struct {
	storage storage.Storage
	logger  *slog.Logger
}

func NewScheduler(storage storage.Storage, logger *slog.Logger) *Scheduler {
	return &Scheduler{
		storage: storage,
		logger:  logger,
	}
}

func (s *Scheduler) GetDueSources(ctx context.Context) ([]*storage.Source, error) {
	sources, err := s.storage.GetEnabledSources(ctx)
	if err != nil {
		return nil, err
	}

	var dueSources []*storage.Source
	now := time.Now()

	for _, source := range sources {
		if s.isDue(source, now) {
			dueSources = append(dueSources, source)
		}
	}

	s.logger.Debug("checked source schedule",
		slog.Int("total_sources", len(sources)),
		slog.Int("due_sources", len(dueSources)),
	)

	return dueSources, nil
}

func (s *Scheduler) isDue(source *storage.Source, now time.Time) bool {

	if source.LastScrapedAt == nil {
		return true
	}

	interval, err := parseInterval(source.ScrapeInterval)
	if err != nil {
		s.logger.Warn("Interval çözümlenemedi, 1 saat varsayılan olarak ayarlandı",
			slog.String("source", source.Name),
			slog.String("interval", source.ScrapeInterval),
		)
		interval = 1 * time.Hour
	}

	nextScrapeTime := source.LastScrapedAt.Add(interval)
	return now.After(nextScrapeTime)
}

func parseInterval(interval string) (time.Duration, error) {
	switch interval {
	case "15 minutes":
		return 15 * time.Minute, nil
	case "30 minutes":
		return 30 * time.Minute, nil
	case "1 hour":
		return 1 * time.Hour, nil
	case "2 hours":
		return 2 * time.Hour, nil
	case "4 hours":
		return 4 * time.Hour, nil
	case "6 hours":
		return 6 * time.Hour, nil
	case "12 hours":
		return 12 * time.Hour, nil
	case "1 day":
		return 24 * time.Hour, nil
	default:

		return 1 * time.Hour, nil
	}
}
