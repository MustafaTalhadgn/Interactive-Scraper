package parser

import (
	"regexp"
	"strings"
	"unicode"
)

var (
	multipleSpaces   = regexp.MustCompile(`\s+`)
	multipleNewlines = regexp.MustCompile(`\n{3,}`)
)

func CleanText(text string) string {

	text = strings.TrimSpace(text)

	text = multipleSpaces.ReplaceAllString(text, "\n\n")
	text = multipleNewlines.ReplaceAllString(text, "\n\n")

	text = removeControlChars(text)

	return text
}

func removeControlChars(s string) string {
	var builder strings.Builder
	for _, r := range s {
		if r == '\n' || r == '\t' || !unicode.IsControl(r) {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

func TruncateText(text string, maxLen int) string {

	if len(text) <= maxLen {
		return text
	}

	truncated := text[:maxLen]
	if lastSpace := strings.LastIndex(truncated, " "); lastSpace > 0 {
		truncated = truncated[:lastSpace]
	}

	return truncated + "..."
}

func WordCount(text string) int {
	return len(strings.Fields(text))
}
