package scorer

import (
	"strings"

	"InteractiveScraper/internal/extractor"
)

type ScoringRule interface {
	Apply(features *extractor.Features, content string) int
}

type KeywordRule struct {
	weights *Weights
}

func (r *KeywordRule) Apply(features *extractor.Features, content string) int {
	score := 0

	for keyword, count := range features.KeywordCounts {

		if isThreatKeyword(keyword) {
			weight := getThreatKeywordWeight(keyword)
			score += weight * count
		}
	}

	if score > r.weights.ThreatKeywordMax {
		score = r.weights.ThreatKeywordMax
	}

	return score
}

type CryptoRule struct {
	weights *Weights
}

func (r *CryptoRule) Apply(features *extractor.Features, content string) int {
	score := 0

	score += len(features.BitcoinAddrs) * r.weights.BitcoinAddress

	score += len(features.EthereumAddrs) * r.weights.EthereumAddress

	score += len(features.MoneroAddrs) * r.weights.MoneroAddress

	if score > r.weights.MaxCryptoScore {
		score = r.weights.MaxCryptoScore
	}

	return score
}

type CVERule struct {
	weights *Weights
}

func (r *CVERule) Apply(features *extractor.Features, content string) int {
	score := len(features.CVEs) * r.weights.CVEIdentifier

	if score > r.weights.MaxCVEScore {
		score = r.weights.MaxCVEScore
	}

	return score
}

type ContactRule struct {
	weights *Weights
}

func (r *ContactRule) Apply(features *extractor.Features, content string) int {
	score := 0

	score += len(features.Emails) * r.weights.EmailAddress
	score += len(features.Phones) * r.weights.PhoneNumber

	return score
}

type NetworkRule struct {
	weights *Weights
}

func (r *NetworkRule) Apply(features *extractor.Features, content string) int {
	score := 0

	score += len(features.OnionURLs) * r.weights.OnionURL
	score += len(features.IPAddresses) * r.weights.IPAddress

	return score
}

type ContentRule struct {
	weights *Weights
}

func (r *ContentRule) Apply(features *extractor.Features, content string) int {
	score := 0

	wordCount := len(strings.Fields(content))
	if wordCount > 0 {
		density := float64(features.UniqueTokens) / float64(wordCount)
		if density > 0.3 {
			score += int(float64(r.weights.ThreatKeywordBase) * r.weights.KeywordDensityMultiplier)
		}
	}

	if wordCount > 500 {
		score += r.weights.LongContentBonus
	}

	hasImportantData := len(features.CVEs) > 0 ||
		len(features.BitcoinAddrs) > 0 ||
		len(features.OnionURLs) > 0 ||
		len(features.KeywordCounts) > 0

	if wordCount < 50 && !hasImportantData {
		score += r.weights.EmptyContentPenalty
	}
	return score
}

func isThreatKeyword(keyword string) bool {
	threatKeywords := map[string]bool{
		"exploit": true, "0day": true, "zeroday": true, "ransomware": true,
		"malware": true, "trojan": true, "backdoor": true, "rootkit": true,
		"leak": true, "breach": true, "dump": true, "phishing": true,
		"ddos": true, "vulnerability": true, "bypass": true, "injection": true,
		"rce": true, "lfi": true, "xss": true, "sqli": true,
	}
	return threatKeywords[strings.ToLower(keyword)]
}

func getThreatKeywordWeight(keyword string) int {
	weights := map[string]int{
		"0day": 25, "zeroday": 25, "ransomware": 25,
		"exploit": 20, "vulnerability": 20, "backdoor": 20, "rootkit": 20,
		"leak": 20, "breach": 20,
		"malware": 18, "trojan": 18, "dump": 18,
		"phishing": 15, "ddos": 15, "bypass": 15, "injection": 15,
		"rce": 20, "lfi": 15, "xss": 15, "sqli": 15,
	}

	if weight, exists := weights[strings.ToLower(keyword)]; exists {
		return weight
	}
	return 5
}
