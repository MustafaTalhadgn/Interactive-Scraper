package main

import (
	"InteractiveScraper/internal/extractor"
	"InteractiveScraper/internal/scorer"
	"fmt"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	s := scorer.NewScorer(scorer.DefaultConfig(), logger)

	// Test scenarios
	testCases := []struct {
		name     string
		features *extractor.Features
		content  string
	}{
		{
			name: "Critical Threat (Ransomware + CVE + Crypto)",
			features: &extractor.Features{
				CVEs:         []string{"CVE-2024-1234", "CVE-2024-5678"},
				BitcoinAddrs: []string{"1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"},
				OnionURLs:    []string{"abc123.onion"},
				Keywords:     []string{"ransomware", "0day", "exploit", "windows"},
				KeywordCounts: map[string]int{
					"ransomware": 3,
					"0day":       2,
					"exploit":    2,
					"windows":    1,
				},
			},
			content: "New ransomware 0day exploit targeting Windows systems...",
		},
		{
			name: "Medium Threat (Single CVE)",
			features: &extractor.Features{
				CVEs:     []string{"CVE-2024-1234"},
				Emails:   []string{"contact@example.com"},
				Keywords: []string{"vulnerability", "patch", "update"},
				KeywordCounts: map[string]int{
					"vulnerability": 1,
					"patch":         1,
				},
			},
			content: "A vulnerability was discovered in the software...",
		},
		{
			name: "Low Threat (General discussion)",
			features: &extractor.Features{
				Keywords: []string{"discussion", "forum", "post"},
				KeywordCounts: map[string]int{
					"discussion": 1,
				},
			},
			content: "General discussion about security topics...",
		},
	}

	for i, tc := range testCases {
		fmt.Printf("\n=== Test %d: %s ===\n", i+1, tc.name)

		result, explanation := s.ScoreWithExplanation(tc.features, tc.content)

		fmt.Printf("Score:       %d/100\n", result.Score)
		fmt.Printf("Criticality: %s\n", result.Criticality)
		fmt.Printf("Color:       %s\n", result.CriticalityColor)
		fmt.Printf("\nBreakdown:\n")
		fmt.Printf("  Keywords:  +%d\n", result.Breakdown.KeywordScore)
		fmt.Printf("  Crypto:    +%d\n", result.Breakdown.CryptoScore)
		fmt.Printf("  CVE:       +%d\n", result.Breakdown.CVEScore)
		fmt.Printf("  Network:   +%d\n", result.Breakdown.NetworkScore)
		fmt.Printf("  Content:   +%d\n", result.Breakdown.ContentScore)
		fmt.Printf("  Raw Total: %d\n", result.Breakdown.TotalRawScore)
		fmt.Printf("\nExplanation: %s\n", explanation)
	}
}
