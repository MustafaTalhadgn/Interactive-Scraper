package models

import "time"

type StandardResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type IntelligenceFeedItem struct {
	ID               int       `json:"id"`
	Title            string    `json:"title"`
	SourceName       string    `json:"source_name"`
	CriticalityScore int       `json:"criticality_score"`
	Criticality      string    `json:"criticality"`
	Category         string    `json:"category"`
	CreatedAt        time.Time `json:"created_at"`
}

type IntelligenceFeedResponse struct {
	Intelligence []IntelligenceFeedItem `json:"intelligence"`
	Pagination   PaginationInfo         `json:"pagination"`
}

type PaginationInfo struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type IntelligenceDetailResponse struct {
	ID                int                `json:"id"`
	Title             string             `json:"title"`
	Summary           string             `json:"summary"`
	SourceName        string             `json:"source_name"`
	SourceURL         string             `json:"source_url"`
	CriticalityScore  int                `json:"criticality_score"`
	Criticality       string             `json:"criticality"`
	Category          string             `json:"category"`
	PublishedAt       *time.Time         `json:"published_at"`
	CreatedAt         time.Time          `json:"created_at"`
	ExtractedFeatures *ExtractedFeatures `json:"extracted_features,omitempty"`
}

type ExtractedFeatures struct {
	BitcoinAddrs  []string `json:"bitcoin_addresses"`
	EthereumAddrs []string `json:"ethereum_addresses"`
	MoneroAddrs   []string `json:"monero_addresses"`
	OnionURLs     []string `json:"onion_urls"`
	IPAddresses   []string `json:"ip_addresses"`
	Emails        []string `json:"emails"`
	CVEs          []string `json:"cves"`
	Keywords      []string `json:"keywords"`
}

type StatsOverview struct {
	TotalIntelligence int            `json:"total_intelligence"`
	CriticalCount     int            `json:"critical_count"`
	Last24Hours       int            `json:"last_24_hours"`
	ActiveSources     int            `json:"active_sources"`
	CriticalityDist   map[string]int `json:"criticality_distribution"`
	CategoryDist      map[string]int `json:"category_distribution"`
}

type TimeSeriesData struct {
	Date     string `json:"date"`
	Critical int    `json:"critical"`
	High     int    `json:"high"`
	Medium   int    `json:"medium"`
	Low      int    `json:"low"`
}

type TimelineResponse struct {
	Timeline []TimeSeriesData `json:"timeline"`
}

type SourceItem struct {
	ID             int        `json:"id"`
	Name           string     `json:"name"`
	URL            string     `json:"url"`
	Category       string     `json:"category"`
	Criticality    string     `json:"criticality"`
	Enabled        bool       `json:"enabled"`
	ScrapeInterval string     `json:"scrape_interval"`
	LastScrapedAt  *time.Time `json:"last_scraped_at"`
	CreatedAt      time.Time  `json:"created_at"`
}

type SourcesResponse struct {
	Sources []SourceItem `json:"sources"`
}

type SourceCreateRequest struct {
	Name           string `json:"name" binding:"required"`
	URL            string `json:"url" binding:"required,url"`
	Category       string `json:"category" binding:"required"`
	Criticality    string `json:"criticality" binding:"required,oneof=low medium high critical"`
	ScrapeInterval string `json:"scrape_interval"`
}

type SourceUpdateRequest struct {
	Name           string `json:"name"`
	Criticality    string `json:"criticality,omitempty" binding:"oneof=low medium high critical"`
	Enabled        *bool  `json:"enabled"`
	ScrapeInterval string `json:"scrape_interval"`
}
