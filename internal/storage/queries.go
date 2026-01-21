package storage

const (
	// Existing queries...
	queryGetEnabledSources = `
		SELECT id, name, url, category, criticality, enabled, 
		       scrape_interval, last_scraped_at, created_at
		FROM sources
		WHERE enabled = true
		ORDER BY last_scraped_at ASC NULLS FIRST
	`

	queryGetSourceByID = `
		SELECT id, name, url, category, criticality, enabled,
		       scrape_interval, last_scraped_at, created_at
		FROM sources
		WHERE id = $1
	`

	queryUpdateSourceLastScraped = `
		UPDATE sources
		SET last_scraped_at = $1
		WHERE id = $2
	`

	queryInsertIntelligence = `
		INSERT INTO intelligence_data 
		(source_id, title, raw_content, source_url, criticality_score, published_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`

	queryGetIntelligenceByID = `
		SELECT id, source_id, title, raw_content, source_url,
		       criticality_score, published_at, created_at
		FROM intelligence_data
		WHERE id = $1
	`

	queryGetRecentIntelligence = `
		SELECT id, source_id, title, raw_content, source_url,
		       criticality_score, published_at, created_at
		FROM intelligence_data
		ORDER BY created_at DESC
		LIMIT $1
	`

	queryGetIntelligenceByCriticality = `
		SELECT id, source_id, title, raw_content, source_url,
		       criticality_score, published_at, created_at
		FROM intelligence_data
		WHERE criticality_score >= $1
		ORDER BY criticality_score DESC, created_at DESC
		LIMIT $2
	`

	queryInsertFeatures = `
		INSERT INTO extracted_features
		(intelligence_id, bitcoin_addrs, ethereum_addrs, monero_addrs,
		 onion_urls, ip_addresses, emails, cves, keywords)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at
	`

	queryGetFeaturesByIntelligenceID = `
		SELECT id, intelligence_id, bitcoin_addrs, ethereum_addrs, monero_addrs,
		       onion_urls, ip_addresses, emails, cves, keywords, created_at
		FROM extracted_features
		WHERE intelligence_id = $1
	`

	queryGetTotalIntelligenceCount = `
		SELECT COUNT(*) FROM intelligence_data
	`
	queryGetCriticalityDistribution = `
		SELECT 
			CASE 
				WHEN criticality_score >= 76 THEN 'critical'
				WHEN criticality_score >= 51 THEN 'high'
				WHEN criticality_score >= 26 THEN 'medium'
				ELSE 'low'
			END as criticality,
			COUNT(*) as count
		FROM intelligence_data
		GROUP BY criticality
		ORDER BY criticality DESC
	`

	queryGetIntelligenceFeed = `
		SELECT 
			i.id,
			i.title,
			i.source_id,
			s.name as source_name,
			s.category,
			i.criticality_score,
			i.created_at
		FROM intelligence_data i
		JOIN sources s ON i.source_id = s.id
		WHERE 1=1
	`

	queryCountIntelligenceFeed = `
		SELECT COUNT(*)
		FROM intelligence_data i
		JOIN sources s ON i.source_id = s.id
		WHERE 1=1
	`

	queryGetAllSources = `
		SELECT id, name, url, category, criticality, enabled,
		       scrape_interval, last_scraped_at, created_at
		FROM sources
		ORDER BY name ASC
	`

	queryCreateSource = `
		INSERT INTO sources 
		(name, url, category, criticality, scrape_interval, enabled)
		VALUES ($1, $2, $3, $4, $5, true)
		RETURNING id, name, url, category, criticality, enabled, 
		          scrape_interval, last_scraped_at, created_at
	`

	queryUpdateSource = `
		UPDATE sources
		SET name = COALESCE(NULLIF($1, ''), name),
		    criticality = COALESCE(NULLIF($2, ''), criticality),
		    enabled = COALESCE($3, enabled),
		    scrape_interval = COALESCE(NULLIF($4, ''), scrape_interval)
		WHERE id = $5
		RETURNING id, name, url, category, criticality, enabled,
		          scrape_interval, last_scraped_at, created_at
	`

	queryDeleteSource = `
		DELETE FROM sources WHERE id = $1
	`
	queryTriggerScrape = `UPDATE sources SET last_scraped_at = NULL WHERE id = $1`

	queryGetIntelligenceCountByCriticality = `
		SELECT COUNT(*)
		FROM intelligence_data
		WHERE criticality_score >= $1
	`

	queryGetIntelligenceCountSince = `
		SELECT COUNT(*)
		FROM intelligence_data
		WHERE created_at >= $1
	`

	queryGetCategoryDistribution = `
		SELECT s.category, COUNT(i.id) as count
		FROM intelligence_data i
		JOIN sources s ON i.source_id = s.id
		GROUP BY s.category
		ORDER BY count DESC
	`

	queryGetTimelineData = `
		SELECT 
			DATE(created_at) as date,
			COUNT(CASE WHEN criticality_score >= 76 THEN 1 END) as critical,
			COUNT(CASE WHEN criticality_score >= 51 AND criticality_score < 76 THEN 1 END) as high,
			COUNT(CASE WHEN criticality_score >= 26 AND criticality_score < 51 THEN 1 END) as medium,
			COUNT(CASE WHEN criticality_score < 26 THEN 1 END) as low
		FROM intelligence_data
		WHERE created_at >= CURRENT_DATE - $1::integer * INTERVAL '1 day'
		GROUP BY DATE(created_at)
		ORDER BY date ASC
	`
)
