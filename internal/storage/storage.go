package storage

import (
	"context"
	"time"
)

type Storage interface {
	GetEnabledSources(ctx context.Context) ([]*Source, error)
	GetSourceByID(ctx context.Context, id int) (*Source, error)
	UpdateSourceLastScraped(ctx context.Context, sourceID int, scrapedAt time.Time) error

	SaveIntelligence(ctx context.Context, input *IntelligenceInput) (*IntelligenceData, error)
	GetIntelligenceByID(ctx context.Context, id int) (*IntelligenceData, error)
	GetRecentIntelligence(ctx context.Context, limit int) ([]*IntelligenceData, error)
	GetIntelligenceByCriticality(ctx context.Context, minScore int, limit int) ([]*IntelligenceData, error)

	SaveFeatures(ctx context.Context, intelligenceID int, features *FeatureInput) error
	GetFeaturesByIntelligenceID(ctx context.Context, intelligenceID int) (*ExtractedFeatures, error)

	GetTotalCount(ctx context.Context) (int, error)
	GetCriticalityDistribution(ctx context.Context) (map[string]int, error)

	Ping(ctx context.Context) error
	Close() error
}
