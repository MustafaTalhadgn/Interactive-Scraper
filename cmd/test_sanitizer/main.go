package main

import (
	"fmt"
	"log/slog"
	"os"

	"InteractiveScraper/internal/sanitizer"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	s := sanitizer.NewSanitizer(sanitizer.DefaultConfig(), logger)

	// Test 1: XSS Attack
	fmt.Println("=== Test 1: XSS Attack ===")
	xss := `<div>Hello <script>alert('XSS')</script> World</div>`
	result := s.Sanitize(xss)
	fmt.Printf("PlainText: %s\n", result.PlainText)
	fmt.Printf("SafeHTML:  %s\n", result.SafeHTML)
	fmt.Printf("Dangerous: %v\n\n", result.WasDangerous)

	// Test 2: SQL Injection
	fmt.Println("=== Test 2: SQL Injection ===")
	sqli := `Username: admin'; DROP TABLE users; --`
	result = s.Sanitize(sqli)
	fmt.Printf("PlainText: %s\n", result.PlainText)
	fmt.Printf("Dangerous: %v\n\n", result.WasDangerous)

	// Test 3: Iframe Injection
	fmt.Println("=== Test 3: Iframe Injection ===")
	iframe := `<iframe src="http://evil.com"></iframe>`
	result = s.Sanitize(iframe)
	fmt.Printf("PlainText: %s\n", result.PlainText)
	fmt.Printf("SafeHTML:  %s\n\n", result.SafeHTML)

	// Test 4: Control Characters
	fmt.Println("=== Test 4: Control Characters ===")
	control := "Text\x00with\x01null\x02bytes"
	result = s.Sanitize(control)
	fmt.Printf("PlainText: %s\n\n", result.PlainText)

	// Test 5: Safe Content
	fmt.Println("=== Test 5: Safe Content ===")
	safe := `<p>This is <strong>safe</strong> content</p>`
	result = s.Sanitize(safe)
	fmt.Printf("PlainText: %s\n", result.PlainText)
	fmt.Printf("SafeHTML:  %s\n", result.SafeHTML)
	fmt.Printf("Dangerous: %v\n", result.WasDangerous)
}
