package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"InteractiveScraper/internal/config"
	"InteractiveScraper/internal/transport"
)

func main() {
	// Setup logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// Load config
	cfg := config.Load()

	// Create Tor client
	client, err := transport.NewTorClient(cfg, logger)
	if err != nil {
		logger.Error("failed to create tor client", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Test fetch
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	testURL := "http://2gzyxa5ihm7nsggfxnu52rck2vv4rvmdlkiu3zzui5du4xyclen53wid.onion/" // Tor Project onion
	logger.Info("testing fetch", slog.String("url", testURL))

	resp, err := client.Fetch(ctx, testURL)
	if err != nil {
		logger.Error("fetch failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Read body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("failed to read body", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.Info("fetch successful",
		slog.Int("status", resp.StatusCode),
		slog.Int("body_size", len(body)),
	)

	fmt.Println("First 200 chars of body:")
	if len(body) > 200 {
		fmt.Println(string(body[:200]))
	} else {
		fmt.Println(string(body))
	}
}
