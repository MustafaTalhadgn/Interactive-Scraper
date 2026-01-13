CREATE TABLE IF NOT EXISTS sources (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    url TEXT NOT NULL UNIQUE,
    category VARCHAR(50), -- forum, market, paste
    criticality VARCHAR(20) 
        CHECK (criticality IN ('low', 'medium', 'high', 'critical')),
    enabled BOOLEAN DEFAULT true,
    scrape_interval INTERVAL DEFAULT '1 hour',
    last_scraped_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS intelligence_data (
    id SERIAL PRIMARY KEY,
    source_id INTEGER 
        REFERENCES sources(id) 
        ON DELETE CASCADE,
    title VARCHAR(500),
    analyzed_content  TEXT NOT NULL,
    source_url TEXT NOT NULL,
    criticality_score INT 
        CHECK (criticality_score BETWEEN 0 AND 100),
    published_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);


CREATE INDEX idx_sources_enabled ON sources(enabled);
CREATE INDEX idx_intel_criticality ON intelligence_data(criticality_score);
CREATE INDEX idx_intel_created_at ON intelligence_data(created_at);

INSERT INTO sources (name, url, category, criticality, scrape_interval) VALUES
('Tor News Forum', 'http://zsxjtsgzborzdllyp64c6pwnjz5eic76bsksbxzqefzogwcydnkjy3yd.onion/', 'forum', 'high', '1 hour'),
('Black Market HQ', 'https://dreadytofatroptsdj6io7l3xptbet6onoyno2yv7jicoxknyazubrad.onion/', 'marketplace', 'critical', '4 hours');
