package normalizer

import "strings"

type StopwordManager struct {
	stopwords map[string]bool
	enabled   bool
}

func NewStopwordManager(enabled bool) *StopwordManager {
	return &StopwordManager{
		stopwords: getDefaultStopwords(),
		enabled:   enabled,
	}
}

func (s *StopwordManager) RemoveStopwords(text string) string {
	if !s.enabled {
		return text
	}

	words := strings.Fields(text)
	var filtered []string

	for _, word := range words {
		if !s.stopwords[word] {
			filtered = append(filtered, word)
		}
	}

	return strings.Join(filtered, " ")
}

func (s *StopwordManager) IsStopword(word string) bool {
	return s.stopwords[word]
}

func (s *StopwordManager) AddStopwords(words []string) {
	for _, word := range words {
		s.stopwords[strings.ToLower(word)] = true
	}
}

func getDefaultStopwords() map[string]bool {
	words := []string{
		// Articles
		"a", "an", "the",

		// zamirler
		"i", "you", "he", "she", "it", "we", "they",
		"me", "him", "her", "us", "them",
		"my", "your", "his", "its", "our", "their",

		// edatlar
		"in", "on", "at", "to", "for", "of", "with",
		"from", "by", "about", "as", "into", "through",

		// Bağlaçlar
		"and", "or", "but", "if", "while", "because",

		// Fiiller
		"is", "are", "was", "were", "be", "been", "being",
		"have", "has", "had", "do", "does", "did",
		"will", "would", "could", "should", "may", "might",

		// Diğerleri
		"this", "that", "these", "those",
		"what", "which", "who", "where", "when", "why", "how",
		"not", "no", "yes",
	}

	stopwordMap := make(map[string]bool)
	for _, word := range words {
		stopwordMap[word] = true
	}

	return stopwordMap
}
