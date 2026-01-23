package storage

import (
	"context"
	"time"

	"InteractiveScraper/internal/api/models" // <--- BU EKLENDİ
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

	GetIntelligenceFeed(ctx context.Context, filters IntelligenceFilters) (*IntelligenceFeedResult, error)

	GetAllSources(ctx context.Context) ([]*Source, error)
	CreateSource(ctx context.Context, input *SourceCreateInput) (*Source, error)
	UpdateSource(ctx context.Context, id int, input *SourceUpdateInput) (*Source, error)
	DeleteSource(ctx context.Context, id int) error

	GetIntelligenceCountByCriticality(ctx context.Context, minScore int) (int, error)
	GetIntelligenceCountSince(ctx context.Context, since time.Time) (int, error)
	GetCategoryDistribution(ctx context.Context) (map[string]int, error)
	GetTimelineData(ctx context.Context, days int) ([]TimeSeriesData, error)
	TriggerManualScrape(ctx context.Context, sourceID int) error
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error

	Ping(ctx context.Context) error
	Close() error
}
