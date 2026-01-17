package scraper

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"InteractiveScraper/internal/config"
	"InteractiveScraper/internal/extractor"
	"InteractiveScraper/internal/normalizer"
	"InteractiveScraper/internal/parser"
	"InteractiveScraper/internal/sanitizer"
	"InteractiveScraper/internal/scorer"
	"InteractiveScraper/internal/storage"
	"InteractiveScraper/internal/transport"
	"InteractiveScraper/internal/validation"
)

// Scraper orchestrates the complete scraping process
type Scraper struct {
	config    *config.Config
	storage   storage.Storage
	pipeline  *Pipeline
	scheduler *Scheduler
	logger    *slog.Logger

	stats      *Stats
	statsMutex sync.RWMutex
}

type Stats struct {
	TotalScraped int
	TotalErrors  int
	LastScrapeAt *time.Time
	LastError    error
}

func NewScraper(
	cfg *config.Config,
	storage storage.Storage,
	logger *slog.Logger,
) (*Scraper, error) {

	torClient, err := transport.NewTorClient(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create tor client: %w", err)
	}

	validator := validation.NewValidator()
	parserInstance := parser.NewParser(parser.DefaultParserConfig())
	sanitizerInstance := sanitizer.NewSanitizer(sanitizer.DefaultConfig(), logger)
	normalizerInstance := normalizer.NewNormalizer(normalizer.DefaultConfig(), logger)
	extractorInstance := extractor.NewExtractor(extractor.DefaultConfig(), logger)
	scorerInstance := scorer.NewScorer(scorer.DefaultConfig(), logger)

	pipeline := NewPipeline(
		torClient,
		validator,
		parserInstance,
		sanitizerInstance,
		normalizerInstance,
		extractorInstance,
		scorerInstance,
		logger,
	)

	scheduler := NewScheduler(storage, logger)

	return &Scraper{
		config:    cfg,
		storage:   storage,
		pipeline:  pipeline,
		scheduler: scheduler,
		logger:    logger,
		stats:     &Stats{},
	}, nil
}

func (s *Scraper) Run(ctx context.Context) error {
	s.logger.Info("scraper started",
		slog.Duration("interval", s.config.ScraperInterval),
	)

	ticker := time.NewTicker(s.config.ScraperInterval)
	defer ticker.Stop()

	if err := s.scrapeOnce(ctx); err != nil {
		s.logger.Error("initial scrape failed", slog.String("error", err.Error()))
	}

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("scraper shutting down")
			return ctx.Err()

		case <-ticker.C:
			if err := s.scrapeOnce(ctx); err != nil {
				s.logger.Error("scrape cycle failed", slog.String("error", err.Error()))
			}
		}
	}
}

func (s *Scraper) scrapeOnce(ctx context.Context) error {
	start := time.Now()

	sources, err := s.scheduler.GetDueSources(ctx)
	if err != nil {
		return fmt.Errorf("failed to get due sources: %w", err)
	}

	if len(sources) == 0 {
		s.logger.Debug("no sources due for scraping")
		return nil
	}

	s.logger.Info("starting scrape cycle",
		slog.Int("sources_count", len(sources)),
	)

	successCount := 0
	errorCount := 0

	for _, source := range sources {

		if ctx.Err() != nil {
			s.logger.Warn("scrape cycle cancelled")
			break
		}

		result := s.pipeline.Process(ctx, source)

		if result.Success {

			if err := s.saveIntelligence(ctx, result); err != nil {
				s.logger.Error("failed to save intelligence",
					slog.Int("source_id", source.ID),
					slog.String("error", err.Error()),
				)
				errorCount++
				s.updateStats(false, err)
				continue
			}

			if err := s.storage.UpdateSourceLastScraped(ctx, source.ID, time.Now()); err != nil {
				s.logger.Warn("failed to update source timestamp",
					slog.Int("source_id", source.ID),
					slog.String("error", err.Error()),
				)
			}

			successCount++
			s.updateStats(true, nil)

		} else {
			s.logger.Error("pipeline processing failed",
				slog.Int("source_id", source.ID),
				slog.String("stage", result.Stage),
				slog.String("error", result.Error.Error()),
			)
			errorCount++
			s.updateStats(false, result.Error)
		}

		time.Sleep(2 * time.Second)
	}

	duration := time.Since(start)

	s.logger.Info("scrape cycle completed",
		slog.Int("total_sources", len(sources)),
		slog.Int("success", successCount),
		slog.Int("errors", errorCount),
		slog.Duration("duration", duration),
	)

	return nil
}

func (s *Scraper) saveIntelligence(ctx context.Context, result *ProcessingResult) error {
	_, err := s.storage.SaveIntelligence(ctx, result.IntelligenceData)
	if err != nil {
		return fmt.Errorf("storage failed: %w", err)
	}

	s.logger.Info("intelligence saved",
		slog.Int("source_id", result.SourceID),
		slog.String("title", result.IntelligenceData.Title),
		slog.Int("score", result.IntelligenceData.CriticalityScore),
	)

	return nil
}

func (s *Scraper) updateStats(success bool, err error) {
	s.statsMutex.Lock()
	defer s.statsMutex.Unlock()

	now := time.Now()
	s.stats.LastScrapeAt = &now

	if success {
		s.stats.TotalScraped++
	} else {
		s.stats.TotalErrors++
		s.stats.LastError = err
	}
}

func (s *Scraper) GetStats() Stats {
	s.statsMutex.RLock()
	defer s.statsMutex.RUnlock()

	return Stats{
		TotalScraped: s.stats.TotalScraped,
		TotalErrors:  s.stats.TotalErrors,
		LastScrapeAt: s.stats.LastScrapeAt,
		LastError:    s.stats.LastError,
	}
}

func (s *Scraper) ScrapeSource(ctx context.Context, sourceID int) error {
	source, err := s.storage.GetSourceByID(ctx, sourceID)
	if err != nil {
		return fmt.Errorf("failed to get source: %w", err)
	}

	result := s.pipeline.Process(ctx, source)

	if !result.Success {
		return result.Error
	}

	if err := s.saveIntelligence(ctx, result); err != nil {
		return err
	}

	return s.storage.UpdateSourceLastScraped(ctx, sourceID, time.Now())
}
