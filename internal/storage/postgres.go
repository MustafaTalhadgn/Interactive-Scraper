package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"InteractiveScraper/internal/api/models" // <--- BU EKLENDİ

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

// User Methods
func (s *PostgresStorage) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, password_hash, role, created_at FROM users WHERE username = $1`

	err := s.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID, &user.Username, &user.PasswordHash, &user.Role, &user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("kullanıcı bulunamadı")
	}
	return user, err
}

// CreateUser güncellendi: Artık struct alıyor
func (s *PostgresStorage) CreateUser(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (username, password_hash, role) VALUES ($1, $2, $3) RETURNING id, created_at`

	// ID ve CreatedAt veritabanından dönecek, struct'ı güncelleyelim
	err := s.db.QueryRowContext(ctx, query, user.Username, user.PasswordHash, user.Role).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return fmt.Errorf("kullanıcı oluşturulamadı: %w", err)
	}
	return nil
}

// Source Methods
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
			input.Summary,
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
			Summary:          input.Summary,
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
		&data.Summary,
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
			&data.Summary,
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
			&data.Summary,
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

func (s *PostgresStorage) GetIntelligenceFeed(ctx context.Context, filters IntelligenceFilters) (*IntelligenceFeedResult, error) {

	query := queryGetIntelligenceFeed
	countQuery := queryCountIntelligenceFeed
	args := []interface{}{}
	argCount := 1

	if filters.Criticality != "" {
		var minScore int
		switch filters.Criticality {
		case "critical":
			minScore = 76
		case "high":
			minScore = 51
		case "medium":
			minScore = 26
		case "low":
			minScore = 0
		}

		clause := fmt.Sprintf(" AND i.criticality_score >= $%d", argCount)
		query += clause
		countQuery += clause
		args = append(args, minScore)
		argCount++
	}

	if filters.SourceID > 0 {
		clause := fmt.Sprintf(" AND i.source_id = $%d", argCount)
		query += clause
		countQuery += clause
		args = append(args, filters.SourceID)
		argCount++
	}

	if filters.Category != "" {
		clause := fmt.Sprintf(" AND s.category = $%d", argCount)
		query += clause
		countQuery += clause
		args = append(args, filters.Category)
		argCount++
	}

	if filters.Search != "" {
		clause := fmt.Sprintf(" AND (i.title ILIKE $%d OR i.summary ILIKE $%d)", argCount, argCount)
		query += clause
		countQuery += clause
		searchPattern := "%" + filters.Search + "%"
		args = append(args, searchPattern)
		argCount++
	}

	if filters.DateFrom != nil {
		clause := fmt.Sprintf(" AND i.created_at >= $%d", argCount)
		query += clause
		countQuery += clause
		args = append(args, filters.DateFrom)
		argCount++
	}

	if filters.DateTo != nil {
		clause := fmt.Sprintf(" AND i.created_at <= $%d", argCount)
		query += clause
		countQuery += clause
		args = append(args, filters.DateTo)
		argCount++
	}

	var total int
	err := s.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("count query failed: %w", err)
	}

	query += " ORDER BY i.created_at DESC"
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCount, argCount+1)

	offset := (filters.Page - 1) * filters.Limit
	args = append(args, filters.Limit, offset)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var items []IntelligenceFeedItem
	for rows.Next() {
		var item IntelligenceFeedItem
		err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.SourceID,
			&item.SourceName,
			&item.Category,
			&item.CriticalityScore,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		items = append(items, item)
	}

	return &IntelligenceFeedResult{
		Items: items,
		Total: total,
	}, nil
}

func (s *PostgresStorage) GetAllSources(ctx context.Context) ([]*Source, error) {
	rows, err := s.db.QueryContext(ctx, queryGetAllSources)
	if err != nil {
		return nil, fmt.Errorf("sorgu başarısız oldu: %w", err)
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

	return sources, nil
}

func (s *PostgresStorage) CreateSource(ctx context.Context, input *SourceCreateInput) (*Source, error) {
	source := &Source{}

	err := s.db.QueryRowContext(
		ctx,
		queryCreateSource,
		input.Name,
		input.URL,
		input.Category,
		input.Criticality,
		input.ScrapeInterval,
	).Scan(
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
		return nil, fmt.Errorf("create failed: %w", err)
	}

	s.logger.Info("kaynak oluşturuldu",
		slog.Int("id", source.ID),
		slog.String("name", source.Name),
	)

	return source, nil
}

func (s *PostgresStorage) UpdateSource(ctx context.Context, id int, input *SourceUpdateInput) (*Source, error) {

	const query = `
        UPDATE sources 
        SET name = COALESCE(NULLIF($1, ''), name), 
            criticality = COALESCE(NULLIF($2, ''), criticality),
            enabled = COALESCE($3, enabled),
            scrape_interval = COALESCE(NULLIF($4, '')::INTERVAL, scrape_interval)
        WHERE id = $5
        RETURNING id, name, url, category, criticality, enabled, 
                  scrape_interval, last_scraped_at, created_at`

	source := &Source{}

	err := s.db.QueryRowContext(
		ctx,
		query,
		input.Name,
		input.Criticality,
		input.Enabled,
		input.ScrapeInterval,
		id,
	).Scan(
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
		// Hata detayını görelim
		return nil, fmt.Errorf("güncelleme başarısız oldu (SQL Error): %w", err)
	}

	s.logger.Info("kaynak güncellendi",
		slog.Int("id", source.ID),
	)

	return source, nil
}

func (s *PostgresStorage) DeleteSource(ctx context.Context, id int) error {
	result, err := s.db.ExecContext(ctx, queryDeleteSource, id)
	if err != nil {
		return fmt.Errorf("silme işlemi başarısız oldu: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("kaynak bulunamadı: %d", id)
	}

	s.logger.Info("kaynak silindi", slog.Int("id", id))

	return nil
}

func (s *PostgresStorage) GetIntelligenceCountByCriticality(ctx context.Context, minScore int) (int, error) {
	var count int
	err := s.db.QueryRowContext(ctx, queryGetIntelligenceCountByCriticality, minScore).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("query failed: %w", err)
	}
	return count, nil
}

func (s *PostgresStorage) GetIntelligenceCountSince(ctx context.Context, since time.Time) (int, error) {
	var count int
	err := s.db.QueryRowContext(ctx, queryGetIntelligenceCountSince, since).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("query failed: %w", err)
	}
	return count, nil
}

func (s *PostgresStorage) GetCategoryDistribution(ctx context.Context) (map[string]int, error) {
	rows, err := s.db.QueryContext(ctx, queryGetCategoryDistribution)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	distribution := make(map[string]int)
	for rows.Next() {
		var category string
		var count int
		if err := rows.Scan(&category, &count); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		distribution[category] = count
	}

	return distribution, nil
}

func (s *PostgresStorage) GetTimelineData(ctx context.Context, days int) ([]TimeSeriesData, error) {
	rows, err := s.db.QueryContext(ctx, queryGetTimelineData, days)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var timeline []TimeSeriesData
	for rows.Next() {
		var data TimeSeriesData
		err := rows.Scan(
			&data.Date,
			&data.Critical,
			&data.High,
			&data.Medium,
			&data.Low,
		)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		timeline = append(timeline, data)
	}

	return timeline, nil
}

func (s *PostgresStorage) TriggerManualScrape(ctx context.Context, sourceID int) error {

	s.logger.Info("manuel scraper tetiklendi", slog.Int("source_id", sourceID))
	return nil
}
