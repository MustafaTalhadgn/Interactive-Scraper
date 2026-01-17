package storage

import (
	"time"
)

type Source struct {
	ID             int        `db:"id"`
	Name           string     `db:"name"`
	URL            string     `db:"url"`
	Category       string     `db:"category"`
	Criticality    string     `db:"criticality"`
	Enabled        bool       `db:"enabled"`
	ScrapeInterval string     `db:"scrape_interval"`
	LastScrapedAt  *time.Time `db:"last_scraped_at"`
	CreatedAt      time.Time  `db:"created_at"`
}

type IntelligenceData struct {
	ID               int        `db:"id"`
	SourceID         int        `db:"source_id"`
	Title            string     `db:"title"`
	RawContent       string     `db:"raw_content"`
	SourceURL        string     `db:"source_url"`
	CriticalityScore int        `db:"criticality_score"`
	PublishedAt      *time.Time `db:"published_at"`
	CreatedAt        time.Time  `db:"created_at"`
}

type ExtractedFeatures struct {
	ID             int       `db:"id"`
	IntelligenceID int       `db:"intelligence_id"`
	BitcoinAddrs   []string  `db:"bitcoin_addrs"`
	EthereumAddrs  []string  `db:"ethereum_addrs"`
	MoneroAddrs    []string  `db:"monero_addrs"`
	OnionURLs      []string  `db:"onion_urls"`
	IPAddresses    []string  `db:"ip_addresses"`
	Emails         []string  `db:"emails"`
	CVEs           []string  `db:"cves"`
	Keywords       []string  `db:"keywords"`
	CreatedAt      time.Time `db:"created_at"`
}

type IntelligenceInput struct {
	SourceID         int
	Title            string
	RawContent       string
	SourceURL        string
	CriticalityScore int
	PublishedAt      *time.Time
	Features         *FeatureInput
}

type FeatureInput struct {
	BitcoinAddrs  []string
	EthereumAddrs []string
	MoneroAddrs   []string
	OnionURLs     []string
	IPAddresses   []string
	Emails        []string
	CVEs          []string
	Keywords      []string
}
