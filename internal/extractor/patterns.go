package extractor

import "regexp"

var (
	OnionURLPattern = regexp.MustCompile(`(?i)\b[a-z2-7]{16,56}\.onion\b`)

	BitcoinPattern = regexp.MustCompile(`\b(bc1|[13])[a-zA-HJ-NP-Z0-9]{25,89}\b`)

	MoneroPattern = regexp.MustCompile(`\b[48][a-zA-Z0-9]{94}\b`)

	EthereumPattern = regexp.MustCompile(`\b0x[a-fA-F0-9]{40}\b`)

	EmailPattern = regexp.MustCompile(`\b[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}\b`)

	CVEPattern = regexp.MustCompile(`\bCVE-\d{4}-\d{4,7}\b`)

	IPv4Pattern = regexp.MustCompile(`\b(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\b`)

	DomainPattern = regexp.MustCompile(`\b(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]\b`)

	PhonePattern = regexp.MustCompile(`\+?\d{1,4}[-.\s]?\(?\d{1,4}\)?[-.\s]?\d{1,4}[-.\s]?\d{1,9}`)
)

type PatternType string

const (
	PatternOnion    PatternType = "onion"
	PatternBitcoin  PatternType = "bitcoin"
	PatternMonero   PatternType = "monero"
	PatternEthereum PatternType = "ethereum"
	PatternEmail    PatternType = "email"
	PatternCVE      PatternType = "cve"
	PatternIPv4     PatternType = "ipv4"
	PatternDomain   PatternType = "domain"
	PatternPhone    PatternType = "phone"
)

type Pattern struct {
	Type    PatternType
	Regex   *regexp.Regexp
	Enabled bool
}

func GetAllPatterns() []*Pattern {
	return []*Pattern{
		{Type: PatternOnion, Regex: OnionURLPattern, Enabled: true},
		{Type: PatternBitcoin, Regex: BitcoinPattern, Enabled: true},
		{Type: PatternMonero, Regex: MoneroPattern, Enabled: true},
		{Type: PatternEthereum, Regex: EthereumPattern, Enabled: true},
		{Type: PatternEmail, Regex: EmailPattern, Enabled: true},
		{Type: PatternCVE, Regex: CVEPattern, Enabled: true},
		{Type: PatternIPv4, Regex: IPv4Pattern, Enabled: true},
		{Type: PatternDomain, Regex: DomainPattern, Enabled: false}, // Disabled by default (too noisy)
		{Type: PatternPhone, Regex: PhonePattern, Enabled: false},   // Disabled by default (optional)
	}
}
