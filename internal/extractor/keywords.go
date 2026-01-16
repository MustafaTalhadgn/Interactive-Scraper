package extractor

import (
	"sort"
	"strings"
)

type KeywordExtractor struct {
	minWordLength  int
	maxKeywords    int
	threatKeywords map[string]int
}

func NewKeywordExtractor() *KeywordExtractor {
	return &KeywordExtractor{
		minWordLength:  3,
		maxKeywords:    50,
		threatKeywords: getThreatKeywords(),
	}
}

func (k *KeywordExtractor) ExtractKeywords(text string) ([]string, map[string]int) {
	words := strings.Fields(text)

	wordCounts := make(map[string]int)
	for _, word := range words {
		if len(word) >= k.minWordLength {
			wordCounts[word]++
		}
	}

	type wordFreq struct {
		word  string
		count int
		score int
	}

	var frequencies []wordFreq
	for word, count := range wordCounts {
		score := count

		if weight, exists := k.threatKeywords[word]; exists {
			score += weight * 10
		}

		frequencies = append(frequencies, wordFreq{
			word:  word,
			count: count,
			score: score,
		})
	}

	sort.Slice(frequencies, func(i, j int) bool {
		return frequencies[i].score > frequencies[j].score
	})

	maxKeywords := k.maxKeywords
	if len(frequencies) < maxKeywords {
		maxKeywords = len(frequencies)
	}

	keywords := make([]string, maxKeywords)
	for i := 0; i < maxKeywords; i++ {
		keywords[i] = frequencies[i].word
	}

	return keywords, wordCounts
}

func getThreatKeywords() map[string]int {
	return map[string]int{

		"exploit":       20,
		"vulnerability": 20,
		"0day":          25,
		"zero-day":      25,
		"zeroday":       25,
		"rce":           20,
		"lfi":           15,
		"xss":           15,
		"sqli":          15,
		"injection":     15,
		"bypass":        15,

		"ransomware": 25,
		"malware":    18,
		"trojan":     18,
		"backdoor":   20,
		"rootkit":    20,
		"botnet":     18,
		"rat":        15,

		"leak":        20,
		"breach":      20,
		"dump":        18,
		"database":    15,
		"credentials": 18,
		"passwords":   18,

		"phishing":   15,
		"ddos":       15,
		"bruteforce": 15,
		"cracking":   15,
		"scanning":   12,

		"bitcoin":  12,
		"btc":      12,
		"crypto":   10,
		"wallet":   10,
		"monero":   12,
		"ethereum": 10,

		"windows":   8,
		"linux":     8,
		"apache":    8,
		"nginx":     8,
		"mysql":     8,
		"wordpress": 8,

		"download":   10,
		"upload":     10,
		"execute":    12,
		"shell":      15,
		"admin":      10,
		"root":       12,
		"privilege":  12,
		"escalation": 15,
	}
}

// IsThreatKeyword checks if a word is a high-value threat keyword
func (k *KeywordExtractor) IsThreatKeyword(word string) bool {
	_, exists := k.threatKeywords[strings.ToLower(word)]
	return exists
}

// GetThreatKeywordWeight returns the weight of a threat keyword
func (k *KeywordExtractor) GetThreatKeywordWeight(word string) int {
	weight, _ := k.threatKeywords[strings.ToLower(word)]
	return weight
}
