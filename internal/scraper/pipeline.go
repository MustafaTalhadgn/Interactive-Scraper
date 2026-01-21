package scraper

import (
	"InteractiveScraper/internal/extractor"
	"InteractiveScraper/internal/normalizer"
	"InteractiveScraper/internal/parser"
	"InteractiveScraper/internal/sanitizer"
	"InteractiveScraper/internal/scorer"
	"InteractiveScraper/internal/storage"
	"InteractiveScraper/internal/transport"
	"InteractiveScraper/internal/validation"
	"context"
	"fmt"
	"log/slog"
	"time"
)

type Pipeline struct {
	transport  transport.HTTPClient
	validator  *validation.Validator
	parser     *parser.Parser
	sanitizer  *sanitizer.Sanitizer
	normalizer *normalizer.Normalizer
	extractor  *extractor.Extractor
	scorer     *scorer.Scorer
	logger     *slog.Logger
}

func NewPipeline(
	transport transport.HTTPClient,
	validator *validation.Validator,
	parser *parser.Parser,
	sanitizer *sanitizer.Sanitizer,
	normalizer *normalizer.Normalizer,
	extractor *extractor.Extractor,
	scorer *scorer.Scorer,
	logger *slog.Logger,
) *Pipeline {
	return &Pipeline{
		transport:  transport,
		validator:  validator,
		parser:     parser,
		sanitizer:  sanitizer,
		normalizer: normalizer,
		extractor:  extractor,
		scorer:     scorer,
		logger:     logger,
	}
}

type ProcessingResult struct {
	SourceID         int
	SourceURL        string
	Success          bool
	IntelligenceData *storage.IntelligenceInput
	Error            error
	Duration         time.Duration
	Stage            string
}

func (p *Pipeline) Process(ctx context.Context, source *storage.Source) *ProcessingResult {
	start := time.Now()

	result := &ProcessingResult{
		SourceID:  source.ID,
		SourceURL: source.URL,
		Success:   false,
	}

	p.logger.Info("starting pipeline",
		slog.Int("source_id", source.ID),
		slog.String("source_name", source.Name),
		slog.String("url", source.URL),
	)

	p.logger.Debug("stage: transport")
	resp, err := p.transport.Fetch(ctx, source.URL)
	if err != nil {
		result.Error = fmt.Errorf("transport failed: %w", err)
		result.Stage = "transport"
		result.Duration = time.Since(start)
		return result
	}
	defer resp.Body.Close()

	p.logger.Debug("stage: validation")
	validated, err := p.validator.Validate(resp)
	if err != nil {
		result.Error = fmt.Errorf("validation failed: %w", err)
		result.Stage = "validation"
		result.Duration = time.Since(start)
		return result
	}

	p.logger.Debug("stage: parser")
	parsed, err := p.parser.Parse(validated.Body, source.URL)
	if err != nil {
		result.Error = fmt.Errorf("parsing failed: %w", err)
		result.Stage = "parser"
		result.Duration = time.Since(start)
		return result
	}

	p.logger.Debug("stage: sanitizer")
	sanitizedTitle := p.sanitizer.SanitizeTitle(parsed.Title)
	sanitizedContent := p.sanitizer.SanitizeContent(parsed.Content)

	p.logger.Debug("stage: normalizer")
	normalizedForScoring := p.normalizer.NormalizeForScoring(sanitizedContent)
	normalizedForDisplay := p.normalizer.NormalizeForDisplay(sanitizedContent)

	p.logger.Debug("stage: extractor")
	features := p.extractor.Extract(normalizedForScoring)

	p.logger.Debug("stage: scorer")
	scoreResult := p.scorer.Score(features, normalizedForScoring)

	intelligenceInput := &storage.IntelligenceInput{
		SourceID:         source.ID,
		Title:            sanitizedTitle,
		Summary:          normalizedForDisplay,
		SourceURL:        source.URL,
		CriticalityScore: scoreResult.Score,
		PublishedAt:      &parsed.Date,
		Features: &storage.FeatureInput{
			BitcoinAddrs:  features.BitcoinAddrs,
			EthereumAddrs: features.EthereumAddrs,
			MoneroAddrs:   features.MoneroAddrs,
			OnionURLs:     features.OnionURLs,
			IPAddresses:   features.IPAddresses,
			Emails:        features.Emails,
			CVEs:          features.CVEs,
			Keywords:      features.Keywords[:min(20, len(features.Keywords))],
		},
	}

	result.Success = true
	result.IntelligenceData = intelligenceInput
	result.Duration = time.Since(start)

	p.logger.Info("pipeline başarıyla tamamlandı",
		slog.Int("source_id", source.ID),
		slog.String("title", sanitizedTitle),
		slog.Int("score", scoreResult.Score),
		slog.String("criticality", string(scoreResult.Criticality)),
		slog.Duration("duration", result.Duration),
	)

	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
