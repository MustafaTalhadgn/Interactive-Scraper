package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"InteractiveScraper/internal/config"
	"InteractiveScraper/internal/storage"
)

func main() {
	// Setup logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// Load config
	cfg := config.Load()
	dsn := cfg.GetDSN()

	// Create storage
	store, err := storage.NewPostgresStorage(dsn, logger)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	ctx := context.Background()

	// Test 1: Get enabled sources
	fmt.Println("=== Test 1: Get Enabled Sources ===")
	sources, err := store.GetEnabledSources(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d enabled sources\n", len(sources))
	for _, source := range sources {
		fmt.Printf("  - %s (%s)\n", source.Name, source.URL)
	}

	// Test 2: Save intelligence data
	fmt.Println("\n=== Test 2: Save Intelligence Data ===")
	now := time.Now()
	input := &storage.IntelligenceInput{
		SourceID:         sources[0].ID,
		Title:            "Test Exploit Found",
		RawContent:       "This is a test exploit targeting Windows...",
		SourceURL:        sources[0].URL + "/post/123",
		CriticalityScore: 85,
		PublishedAt:      &now,
		Features: &storage.FeatureInput{
			BitcoinAddrs: []string{"1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"},
			CVEs:         []string{"CVE-2024-1234"},
			Keywords:     []string{"exploit", "windows", "0day"},
		},
	}

	intelligence, err := store.SaveIntelligence(ctx, input)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Saved intelligence: ID=%d, Score=%d\n",
		intelligence.ID, intelligence.CriticalityScore)

	// Test 3: Update source last_scraped_at
	fmt.Println("\n=== Test 3: Update Last Scraped ===")
	err = store.UpdateSourceLastScraped(ctx, sources[0].ID, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Updated source last_scraped_at ✅")

	// Test 4: Get statistics
	fmt.Println("\n=== Test 4: Statistics ===")
	total, _ := store.GetTotalCount(ctx)
	fmt.Printf("Total intelligence records: %d\n", total)

	distribution, _ := store.GetCriticalityDistribution(ctx)
	fmt.Println("Criticality distribution:")
	for level, count := range distribution {
		fmt.Printf("  %s: %d\n", level, count)
	}

	fmt.Println("\n✅ All tests passed!")
}
