package scorer

type Weights struct {
	ThreatKeywordBase int
	ThreatKeywordMax  int

	BitcoinAddress  int
	EthereumAddress int
	MoneroAddress   int
	MaxCryptoScore  int

	CVEIdentifier int
	MaxCVEScore   int

	EmailAddress int
	PhoneNumber  int

	OnionURL  int
	IPAddress int

	KeywordDensityMultiplier float64
	LongContentBonus         int

	EmptyContentPenalty int
}

func DefaultWeights() *Weights {
	return &Weights{

		ThreatKeywordBase: 5,
		ThreatKeywordMax:  40,

		BitcoinAddress:  8,
		EthereumAddress: 6,
		MoneroAddress:   10,
		MaxCryptoScore:  25,

		CVEIdentifier: 15,
		MaxCVEScore:   30,

		EmailAddress: 3,
		PhoneNumber:  2,

		OnionURL:  5,
		IPAddress: 3,

		KeywordDensityMultiplier: 1.5,
		LongContentBonus:         5,

		EmptyContentPenalty: -10,
	}
}

type CustomWeights struct {
	*Weights
	CustomKeywords map[string]int
}

func NewCustomWeights() *CustomWeights {
	return &CustomWeights{
		Weights:        DefaultWeights(),
		CustomKeywords: make(map[string]int),
	}
}

func (w *CustomWeights) AddCustomKeyword(keyword string, weight int) {
	w.CustomKeywords[keyword] = weight
}
