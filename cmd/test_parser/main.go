package main

import (
	"fmt"
	"log"

	"InteractiveScraper/internal/parser"
)

func main() {
	// Sample Dark Web HTML (bozuk olabilir)
	html := `
	<!DOCTYPE html>
	<html>
	<head><title>Dark Web Post</title></head>
	<body>
		<h1 class="post-title">New Exploit Released</h1>
		<span class="author">Anonymous123</span>
		<time datetime="2024-01-15T10:30:00Z">2024-01-15</time>
		
		<div class="post-content">
			This is a new zero-day exploit affecting Windows systems.
			Contact: hacker@example.onion
			Bitcoin: 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa
			
			<script>alert('evil')</script>
		</div>
		
		<a href="http://another.onion">Related Post</a>
		<img src="/image.jpg">
	</body>
	</html>
	`

	// Create parser
	config := parser.DefaultParserConfig()
	config.ExtractLinks = true
	config.ExtractImages = true

	p := parser.NewParser(config)

	// Parse
	data, err := p.Parse([]byte(html), "http://test.onion/post/123")
	if err != nil {
		log.Fatal(err)
	}

	// Print results
	fmt.Println("=== Parsed Data ===")
	fmt.Printf("Title:      %s\n", data.Title)
	fmt.Printf("Author:     %s\n", data.Author)
	fmt.Printf("Date:       %s\n", data.Date)
	fmt.Printf("Word Count: %d\n", data.WordCount)
	fmt.Printf("Links:      %v\n", data.Links)
	fmt.Printf("Images:     %v\n", data.ImageURLs)
	fmt.Printf("\nContent:\n%s\n", data.Content)
}
