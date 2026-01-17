package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/lib/pq"
)

type PostgresStorage struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewPostgresStorage(dsn string, logger *slog.Logger) (*PostgresStorage, error) {

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("Veritabanı açılamadı: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("Veritabanına ping atılamadı: %w", err)
	}

	logger.Info("Veritabanı bağlantısı kuruldu",
		slog.Int("max_open_conns", 25),
		slog.Int("max_idle_conns", 5),
	)

	return &PostgresStorage{
		db:     db,
		logger: logger,
	}, nil
}

func (s *PostgresStorage) GetEnabledSources(ctx context.Context) ([]*Source, error) {
	rows, err := s.db.QueryContext(ctx, queryGetEnabledSources)
	if err != nil {
		return nil, fmt.Errorf("Sorgu Başarısız Oldu: %w", err)
	}
	defer rows.Close()

	var sources []*Source
	for rows.Next() {
		source := &Source{}
		err := rows.Scan(
			&source.ID,
			&source.Name,
			&source.URL,
			&source.Category,
			&source.Criticality,
			&source.Enabled,
			&source.ScrapeInterval,
			&source.LastScrapedAt,
			&source.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("tarama başarısız oldu: %w", err)
		}
		sources = append(sources, source)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("satır yinelemesi başarısız oldu: %w", err)
	}

	s.logger.Debug("etkin kaynaklar alındı",
		slog.Int("count", len(sources)),
	)

	return sources, nil
}

func (s *PostgresStorage) GetSourceByID(ctx context.Context, id int) (*Source, error) {
	source := &Source{}

	err := s.db.QueryRowContext(ctx, queryGetSourceByID, id).Scan(
		&source.ID,
		&source.Name,
		&source.URL,
		&source.Category,
		&source.Criticality,
		&source.Enabled,
		&source.ScrapeInterval,
		&source.LastScrapedAt,
		&source.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("kaynak bulunamadı: %d", id)
	}
	if err != nil {
		return nil, fmt.Errorf("Sorgu Başarısız Oldu: %w", err)
	}

	return source, nil
}

func (s *PostgresStorage) UpdateSourceLastScraped(ctx context.Context, sourceID int, scrapedAt time.Time) error {
	result, err := s.db.ExecContext(ctx, queryUpdateSourceLastScraped, scrapedAt, sourceID)
	if err != nil {
		return fmt.Errorf("Güncelleme Başarısız Oldu: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("kaynak bulunamadı: %d", sourceID)
	}

	s.logger.Debug("kaynak son tarama zamanı güncellendi",
		slog.Int("source_id", sourceID),
		slog.Time("scraped_at", scrapedAt),
	)

	return nil
}

func (s *PostgresStorage) SaveIntelligence(ctx context.Context, input *IntelligenceInput) (*IntelligenceData, error) {
	var intelligence *IntelligenceData

	err := WithTransaction(ctx, s.db, func(tx *sql.Tx) error {

		var id int
		var createdAt time.Time

		err := tx.QueryRowContext(
			ctx,
			queryInsertIntelligence,
			input.SourceID,
			input.Title,
			input.RawContent,
			input.SourceURL,
			input.CriticalityScore,
			input.PublishedAt,
		).Scan(&id, &createdAt)

		if err != nil {
			return fmt.Errorf("intelligence eklenemedi: %w", err)
		}

		intelligence = &IntelligenceData{
			ID:               id,
			SourceID:         input.SourceID,
			Title:            input.Title,
			RawContent:       input.RawContent,
			SourceURL:        input.SourceURL,
			CriticalityScore: input.CriticalityScore,
			PublishedAt:      input.PublishedAt,
			CreatedAt:        createdAt,
		}

		if input.Features != nil {
			_, err = tx.ExecContext(
				ctx,
				queryInsertFeatures,
				id,
				pq.Array(input.Features.BitcoinAddrs),
				pq.Array(input.Features.EthereumAddrs),
				pq.Array(input.Features.MoneroAddrs),
				pq.Array(input.Features.OnionURLs),
				pq.Array(input.Features.IPAddresses),
				pq.Array(input.Features.Emails),
				pq.Array(input.Features.CVEs),
				pq.Array(input.Features.Keywords),
			)

			if err != nil {
				return fmt.Errorf("Özellikler eklenemedi: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	s.logger.Info("intelligence verisi kaydedildi",
		slog.Int("id", intelligence.ID),
		slog.Int("source_id", intelligence.SourceID),
		slog.String("title", intelligence.Title),
		slog.Int("criticality_score", intelligence.CriticalityScore),
	)

	return intelligence, nil
}

func (s *PostgresStorage) GetIntelligenceByID(ctx context.Context, id int) (*IntelligenceData, error) {
	data := &IntelligenceData{}

	err := s.db.QueryRowContext(ctx, queryGetIntelligenceByID, id).Scan(
		&data.ID,
		&data.SourceID,
		&data.Title,
		&data.RawContent,
		&data.SourceURL,
		&data.CriticalityScore,
		&data.PublishedAt,
		&data.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("intelligence verisi bulunamadı: %d", id)
	}
	if err != nil {
		return nil, fmt.Errorf("Sorgu Başarısız Oldu: %w", err)
	}

	return data, nil
}

func (s *PostgresStorage) GetRecentIntelligence(ctx context.Context, limit int) ([]*IntelligenceData, error) {
	rows, err := s.db.QueryContext(ctx, queryGetRecentIntelligence, limit)
	if err != nil {
		return nil, fmt.Errorf("Sorgu Başarısız Oldu: %w", err)
	}
	defer rows.Close()

	var results []*IntelligenceData
	for rows.Next() {
		data := &IntelligenceData{}
		err := rows.Scan(
			&data.ID,
			&data.SourceID,
			&data.Title,
			&data.RawContent,
			&data.SourceURL,
			&data.CriticalityScore,
			&data.PublishedAt,
			&data.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("tarama başarısız oldu: %w", err)
		}
		results = append(results, data)
	}

	return results, nil
}

func (s *PostgresStorage) GetIntelligenceByCriticality(ctx context.Context, minScore int, limit int) ([]*IntelligenceData, error) {
	rows, err := s.db.QueryContext(ctx, queryGetIntelligenceByCriticality, minScore, limit)
	if err != nil {
		return nil, fmt.Errorf("Sorgu Başarısız Oldu: %w", err)
	}
	defer rows.Close()

	var results []*IntelligenceData
	for rows.Next() {
		data := &IntelligenceData{}
		err := rows.Scan(
			&data.ID,
			&data.SourceID,
			&data.Title,
			&data.RawContent,
			&data.SourceURL,
			&data.CriticalityScore,
			&data.PublishedAt,
			&data.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("tarama başarısız oldu: %w", err)
		}
		results = append(results, data)
	}

	return results, nil
}

func (s *PostgresStorage) SaveFeatures(ctx context.Context, intelligenceID int, features *FeatureInput) error {
	_, err := s.db.ExecContext(
		ctx,
		queryInsertFeatures,
		intelligenceID,
		pq.Array(features.BitcoinAddrs),
		pq.Array(features.EthereumAddrs),
		pq.Array(features.MoneroAddrs),
		pq.Array(features.OnionURLs),
		pq.Array(features.IPAddresses),
		pq.Array(features.Emails),
		pq.Array(features.CVEs),
		pq.Array(features.Keywords),
	)

	if err != nil {
		return fmt.Errorf("Özellikler kaydedilemedi: %w", err)
	}

	return nil
}

func (s *PostgresStorage) GetFeaturesByIntelligenceID(ctx context.Context, intelligenceID int) (*ExtractedFeatures, error) {
	features := &ExtractedFeatures{}

	err := s.db.QueryRowContext(ctx, queryGetFeaturesByIntelligenceID, intelligenceID).Scan(
		&features.ID,
		&features.IntelligenceID,
		pq.Array(&features.BitcoinAddrs),
		pq.Array(&features.EthereumAddrs),
		pq.Array(&features.MoneroAddrs),
		pq.Array(&features.OnionURLs),
		pq.Array(&features.IPAddresses),
		pq.Array(&features.Emails),
		pq.Array(&features.CVEs),
		pq.Array(&features.Keywords),
		&features.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("özellikler bulunamadı: %d", intelligenceID)
	}
	if err != nil {
		return nil, fmt.Errorf("Sorgu başarısız oldu: %w", err)
	}

	return features, nil
}

func (s *PostgresStorage) GetTotalCount(ctx context.Context) (int, error) {
	var count int
	err := s.db.QueryRowContext(ctx, queryGetTotalIntelligenceCount).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("Sorgu Başarısız Oldu: %w", err)
	}
	return count, nil
}

func (s *PostgresStorage) GetCriticalityDistribution(ctx context.Context) (map[string]int, error) {
	rows, err := s.db.QueryContext(ctx, queryGetCriticalityDistribution)
	if err != nil {
		return nil, fmt.Errorf("Sorgu Başarısız Oldu: %w", err)
	}
	defer rows.Close()

	distribution := make(map[string]int)
	for rows.Next() {
		var criticality string
		var count int
		if err := rows.Scan(&criticality, &count); err != nil {
			return nil, fmt.Errorf("tarama başarısız oldu: %w", err)
		}
		distribution[criticality] = count
	}

	return distribution, nil
}

func (s *PostgresStorage) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

func (s *PostgresStorage) Close() error {
	s.logger.Info("veritabanı bağlantısı kapatılıyor")
	return s.db.Close()
}
