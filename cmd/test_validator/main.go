package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"InteractiveScraper/internal/validation"
)

func main() {
	validator := validation.NewValidator()

	// Test 1: Valid HTML response
	fmt.Println("=== Test 1: Valid HTML ===")
	testValidHTML(validator)

	// Test 2: Empty body
	fmt.Println("\n=== Test 2: Empty Body ===")
	testEmptyBody(validator)

	// Test 3: Body too large
	fmt.Println("\n=== Test 3: Body Too Large ===")
	testBodyTooLarge(validator)

	// Test 4: Invalid Content-Type
	fmt.Println("\n=== Test 4: Invalid Content-Type ===")
	testInvalidContentType(validator)
}

func testValidHTML(v *validation.Validator) {
	html := "<html><body>Hello Dark Web</body></html>"
	resp := createMockResponse(200, "text/html", html)

	result, err := v.Validate(resp)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Printf("✅ Valid! Body size: %d bytes\n", len(result.Body))
	}
}

func testEmptyBody(v *validation.Validator) {
	resp := createMockResponse(200, "text/html", "")

	_, err := v.Validate(resp)
	if err != nil {
		fmt.Printf("✅ Expected error: %v\n", err)
	} else {
		fmt.Println("❌ Should have failed!")
	}
}

func testBodyTooLarge(v *validation.Validator) {
	// Create 6MB body
	largeBody := bytes.Repeat([]byte("x"), 6*1024*1024)
	resp := createMockResponse(200, "text/html", string(largeBody))

	_, err := v.Validate(resp)
	if err != nil {
		fmt.Printf("✅ Expected error: %v\n", err)
	} else {
		fmt.Println("❌ Should have failed!")
	}
}

func testInvalidContentType(v *validation.Validator) {
	resp := createMockResponse(200, "application/json", `{"error": "not html"}`)

	_, err := v.Validate(resp)
	if err != nil {
		fmt.Printf("✅ Expected error: %v\n", err)
	} else {
		fmt.Println("❌ Should have failed!")
	}
}

func createMockResponse(status int, contentType, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header: http.Header{
			"Content-Type": []string{contentType},
		},
		Body:          io.NopCloser(bytes.NewBufferString(body)),
		ContentLength: int64(len(body)),
	}
}
