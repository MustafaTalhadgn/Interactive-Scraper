package extractor

import (
	"log/slog"
	"strings"
)

type Extractor struct {
	patterns         []*Pattern
	cryptoValidator  *CryptoValidator
	keywordExtractor *KeywordExtractor
	logger           *slog.Logger
	config           *Config
}

type Config struct {
	ExtractOnion    bool
	ExtractCrypto   bool
	ExtractEmail    bool
	ExtractCVE      bool
	ExtractIP       bool
	ExtractDomain   bool
	ExtractPhone    bool
	ExtractKeywords bool
	MaxKeywords     int
}

func DefaultConfig() *Config {
	return &Config{
		ExtractOnion:    true,
		ExtractCrypto:   true,
		ExtractEmail:    true,
		ExtractCVE:      true,
		ExtractIP:       true,
		ExtractDomain:   false,
		ExtractPhone:    false,
		ExtractKeywords: true,
		MaxKeywords:     50,
	}
}

func NewExtractor(config *Config, logger *slog.Logger) *Extractor {
	if config == nil {
		config = DefaultConfig()
	}

	return &Extractor{
		patterns:         GetAllPatterns(),
		cryptoValidator:  NewCryptoValidator(),
		keywordExtractor: NewKeywordExtractor(),
		logger:           logger,
		config:           config,
	}
}

func (e *Extractor) Extract(text string) *Features {
	features := NewFeatures()

	if e.config.ExtractOnion {
		features.OnionURLs = e.extractPattern(text, PatternOnion)
	}

	if e.config.ExtractCrypto {
		features.BitcoinAddrs = e.extractAndValidateBitcoin(text)
		features.MoneroAddrs = e.extractAndValidateMonero(text)
		features.EthereumAddrs = e.extractAndValidateEthereum(text)
	}

	if e.config.ExtractEmail {
		features.Emails = e.extractPattern(text, PatternEmail)
	}

	if e.config.ExtractCVE {
		features.CVEs = e.extractPattern(text, PatternCVE)
	}

	if e.config.ExtractIP {
		features.IPAddresses = e.extractPattern(text, PatternIPv4)
	}

	if e.config.ExtractDomain {
		features.Domains = e.extractPattern(text, PatternDomain)

		features.Domains = e.filterOnionDomains(features.Domains)
	}

	if e.config.ExtractPhone {
		features.Phones = e.extractPattern(text, PatternPhone)
	}

	if e.config.ExtractKeywords {
		keywords, counts := e.keywordExtractor.ExtractKeywords(text)
		features.Keywords = keywords
		features.KeywordCounts = counts
	}

	features.TotalMatches = features.CountTotalMatches()
	features.UniqueTokens = len(features.KeywordCounts)

	e.logger.Debug("feature extraction completed",
		slog.Int("onion_urls", len(features.OnionURLs)),
		slog.Int("bitcoin", len(features.BitcoinAddrs)),
		slog.Int("emails", len(features.Emails)),
		slog.Int("cves", len(features.CVEs)),
		slog.Int("keywords", len(features.Keywords)),
		slog.Int("total_matches", features.TotalMatches),
	)

	return features
}

func (e *Extractor) extractPattern(text string, patternType PatternType) []string {
	for _, pattern := range e.patterns {
		if pattern.Type == patternType && pattern.Enabled {
			matches := pattern.Regex.FindAllString(text, -1)
			return e.deduplicate(matches)
		}
	}
	return []string{}
}

func (e *Extractor) extractAndValidateBitcoin(text string) []string {
	matches := BitcoinPattern.FindAllString(text, -1)

	var validated []string
	for _, match := range matches {
		if e.cryptoValidator.ValidateBitcoin(match) {
			validated = append(validated, match)
		}
	}

	return e.deduplicate(validated)
}

func (e *Extractor) extractAndValidateEthereum(text string) []string {
	matches := EthereumPattern.FindAllString(text, -1)

	var validated []string
	for _, match := range matches {
		if e.cryptoValidator.ValidateEthereum(match) {
			validated = append(validated, match)
		}
	}

	return e.deduplicate(validated)
}

func (e *Extractor) extractAndValidateMonero(text string) []string {
	matches := MoneroPattern.FindAllString(text, -1)

	var validated []string
	for _, match := range matches {
		if e.cryptoValidator.ValidateMonero(match) {
			validated = append(validated, match)
		}
	}

	return e.deduplicate(validated)
}

func (e *Extractor) deduplicate(items []string) []string {
	seen := make(map[string]bool)
	var unique []string

	for _, item := range items {
		if !seen[item] {
			unique = append(unique, item)
			seen[item] = true
		}
	}

	return unique
}

func (e *Extractor) filterOnionDomains(domains []string) []string {
	var filtered []string
	for _, domain := range domains {
		if !strings.HasSuffix(domain, ".onion") {
			filtered = append(filtered, domain)
		}
	}
	return filtered
}
