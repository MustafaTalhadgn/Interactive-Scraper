CREATE TABLE IF NOT EXISTS sources (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    url TEXT NOT NULL UNIQUE,
    category VARCHAR(50), 
    criticality VARCHAR(20) 
        CHECK (criticality IN ('low', 'medium', 'high', 'critical')),
    enabled BOOLEAN DEFAULT true,
    scrape_interval INTERVAL DEFAULT '1 hour',
    last_scraped_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS intelligence_data (
    id SERIAL PRIMARY KEY,
    source_id INTEGER REFERENCES sources(id) ON DELETE CASCADE,
    title VARCHAR(500),
    summary TEXT NOT NULL, 
    source_url TEXT NOT NULL,
    criticality_score INT CHECK (criticality_score BETWEEN 0 AND 100),
    published_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE TABLE extracted_features (
    id SERIAL PRIMARY KEY,
    intelligence_id INTEGER REFERENCES intelligence_data(id) ON DELETE CASCADE,
    bitcoin_addrs TEXT[], 
    ethereum_addrs TEXT[],
    monero_addrs TEXT[],
    onion_urls TEXT[],
    ip_addresses TEXT[],
    emails TEXT[],
    cves TEXT[],
    keywords TEXT[],
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) DEFAULT 'admin',
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_sources_enabled ON sources(enabled);
CREATE INDEX idx_intel_criticality ON intelligence_data(criticality_score);
CREATE INDEX idx_intel_created_at ON intelligence_data(created_at);

INSERT INTO sources (name, url, category, criticality, scrape_interval) VALUES
('Abyss-Data', 'http://3ev4metjirohtdpshsqlkrqcmxq6zu3d7obrdhglpy5jpbr7whmlfgqd.onion/', 'ransomware', 'high', '1 hour'),
('Beast Leaks', 'http://beast6azu4f7fxjakiayhnssybibsgjnmy77a6duufqw5afjzfjhzuqd.onion/', 'ransomware', 'medium', '1 hour'),
('Dread Forum', 'https://dreadytofatroptsdj6io7l3xptbet6onoyno2yv7jicoxknyazubrad.onion/', 'forum', 'low', '1 hour'),
('Benzona', 'http://benzona6x5ggng3hx52h4mak5sgx5vukrdlrrd3of54g2uppqog2joyd.onion/', 'ransomware', 'medium', '1 hour'),
('Blackout', 'http://black3gnkizshuynieigw6ejgpblb53mpasftzd6pydqpmq2vn2xf6yd.onion/', 'ransomware', 'medium', '1 hour'),
('Brain Cipher', 'http://vkvsgl7lhipjirmz6j5ubp3w3bwvxgcdbpi3fsbqngfynetqtw4w5hyd.onion/', 'ransomware', 'medium', '1 hour'),
('CHAOS', 'http://hptqq2o2qjva7lcaaq67w36jihzivkaitkexorauw7b2yul2z6zozpqd.onion/', 'ransomware', 'medium', '1 hour'),
('CiphBit', 'http://ciphbitqyg26jor7eeo6xieyq7reouctefrompp6ogvhqjba7uo4xdid.onion/', 'ransomware', 'medium', '1 hour'),
('ContFR', 'http://zprxx7sfc26rufggreanowmme5qqouqegr2efnko6erycquwvpq5egid.onion/', 'ransomware', 'medium', '1 hour');



