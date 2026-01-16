package main

import (
	"fmt"
	"log/slog"
	"os"

	"InteractiveScraper/internal/extractor"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	e := extractor.NewExtractor(extractor.DefaultConfig(), logger)

	// Sample Dark Web content
	text := `
	New 0day exploit released for Windows! CVE-2024-1234 affects all versions.
	
	Contact: hacker@darkmail.com
	Onion mirror: abc123xyz456def.onion
	
	Payment: 
	BTC: 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa
	ETH: 0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb
	
	This ransomware is undetectable and bypasses Windows Defender.
	Download link: http://another789site.onion/malware.zip
	
	Server: 192.168.1.100
	`

	// Extract features
	features := e.Extract(text)

	// Print results
	fmt.Println("=== Extracted Features ===\n")

	if len(features.OnionURLs) > 0 {
		fmt.Printf("Onion URLs: %v\n", features.OnionURLs)
	}

	if len(features.BitcoinAddrs) > 0 {
		fmt.Printf("Bitcoin:    %v\n", features.BitcoinAddrs)
	}

	if len(features.EthereumAddrs) > 0 {
		fmt.Printf("Ethereum:   %v\n", features.EthereumAddrs)
	}

	if len(features.Emails) > 0 {
		fmt.Printf("Emails:     %v\n", features.Emails)
	}

	if len(features.CVEs) > 0 {
		fmt.Printf("CVEs:       %v\n", features.CVEs)
	}

	if len(features.IPAddresses) > 0 {
		fmt.Printf("IPs:        %v\n", features.IPAddresses)
	}

	fmt.Printf("\nTop 10 Keywords:\n")
	maxKeywords := 10
	if len(features.Keywords) < maxKeywords {
		maxKeywords = len(features.Keywords)
	}
	for i := 0; i < maxKeywords; i++ {
		keyword := features.Keywords[i]
		count := features.KeywordCounts[keyword]
		fmt.Printf("  %s (count: %d)\n", keyword, count)
	}

	fmt.Printf("\nMetadata:\n")
	fmt.Printf("  Total Matches: %d\n", features.TotalMatches)
	fmt.Printf("  Unique Tokens: %d\n", features.UniqueTokens)
	fmt.Printf("  Has Crypto:    %v\n", features.HasCrypto())
	fmt.Printf("  Has CVE:       %v\n", features.HasVulnerabilityInfo())
	fmt.Printf("  Has Contact:   %v\n", features.HasContactInfo())
}
