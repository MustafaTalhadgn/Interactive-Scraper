package storage

import (
	"time"
)

// Existing models...
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
	Summary          string     `db:"summary"`
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
	Summary          string
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

type IntelligenceFilters struct {
	Page        int
	Limit       int
	Criticality string
	SourceID    int
	Category    string
	Search      string
	DateFrom    *time.Time
	DateTo      *time.Time
}

type IntelligenceFeedResult struct {
	Items []IntelligenceFeedItem
	Total int
}

type IntelligenceFeedItem struct {
	ID               int
	Title            string
	SourceName       string
	SourceID         int
	Category         string
	CriticalityScore int
	CreatedAt        time.Time
}

type SourceCreateInput struct {
	Name           string
	URL            string
	Category       string
	Criticality    string
	ScrapeInterval string
}

type SourceUpdateInput struct {
	Name           string
	Criticality    string
	Enabled        *bool
	ScrapeInterval string
}

type TimeSeriesData struct {
	Date     string
	Critical int
	High     int
	Medium   int
	Low      int
}
