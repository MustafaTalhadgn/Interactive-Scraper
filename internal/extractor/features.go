package extractor

type Features struct {
	OnionURLs   []string
	Domains     []string
	IPAddresses []string

	BitcoinAddrs  []string
	MoneroAddrs   []string
	EthereumAddrs []string

	Emails []string
	Phones []string

	CVEs []string

	Keywords      []string
	KeywordCounts map[string]int

	TotalMatches int
	UniqueTokens int
}

func NewFeatures() *Features {
	return &Features{
		OnionURLs:     make([]string, 0),
		Domains:       make([]string, 0),
		IPAddresses:   make([]string, 0),
		BitcoinAddrs:  make([]string, 0),
		MoneroAddrs:   make([]string, 0),
		EthereumAddrs: make([]string, 0),
		Emails:        make([]string, 0),
		Phones:        make([]string, 0),
		CVEs:          make([]string, 0),
		Keywords:      make([]string, 0),
		KeywordCounts: make(map[string]int),
	}
}

func (f *Features) HasCrypto() bool {
	return len(f.BitcoinAddrs) > 0 ||
		len(f.MoneroAddrs) > 0 ||
		len(f.EthereumAddrs) > 0
}

func (f *Features) HasVulnerabilityInfo() bool {
	return len(f.CVEs) > 0
}

func (f *Features) HasContactInfo() bool {
	return len(f.Emails) > 0 || len(f.Phones) > 0
}

func (f *Features) CountTotalMatches() int {
	return len(f.OnionURLs) +
		len(f.Domains) +
		len(f.IPAddresses) +
		len(f.BitcoinAddrs) +
		len(f.MoneroAddrs) +
		len(f.EthereumAddrs) +
		len(f.Emails) +
		len(f.Phones) +
		len(f.CVEs)
}
