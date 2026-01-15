package main

import (
	"fmt"
	"log/slog"
	"os"

	"InteractiveScraper/internal/normalizer"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	n := normalizer.NewNormalizer(normalizer.DefaultConfig(), logger)

	// Test 1: Basic normalization
	fmt.Println("=== Test 1: Basic Normalization ===")
	text1 := "The   EXPLOIT   is   HERE!"
	result1 := n.Normalize(text1)
	fmt.Printf("Original:   '%s'\n", text1)
	fmt.Printf("Normalized: '%s'\n", result1.Text)
	fmt.Printf("Word Count: %d\n\n", result1.WordCount)

	// Test 2: Unicode normalization
	fmt.Println("=== Test 2: Unicode Normalization ===")
	text2 := "café résumé naïve" // Accented characters
	result2 := n.Normalize(text2)
	fmt.Printf("Original:   '%s'\n", text2)
	fmt.Printf("Normalized: '%s'\n\n", result2.Text)

	// Test 3: Normalize for scoring
	fmt.Println("=== Test 3: Normalize for Scoring ===")
	text3 := "The exploit is NOT vulnerable!"
	scoring := n.NormalizeForScoring(text3)
	fmt.Printf("Original: '%s'\n", text3)
	fmt.Printf("Scoring:  '%s'\n\n", scoring)

	// Test 4: Normalize for display
	fmt.Println("=== Test 4: Normalize for Display ===")
	text4 := "Multiple    spaces   and\n\n\nnewlines"
	display := n.NormalizeForDisplay(text4)
	fmt.Printf("Original: '%s'\n", text4)
	fmt.Printf("Display:  '%s'\n\n", display)

	// Test 5: Keyword extraction
	fmt.Println("=== Test 5: Keyword Extraction ===")
	text5 := "This is a critical exploit affecting Windows systems"
	keywords := n.NormalizeKeywords(text5)
	fmt.Printf("Original: '%s'\n", text5)
	fmt.Printf("Keywords: %v\n", keywords)
}
