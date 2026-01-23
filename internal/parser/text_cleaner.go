package parser

import (
	"regexp"
	"strings"
	"unicode"
)

var (
	// Markdown formatını bozmamak için sadece çoklu boşlukları hedefle
	multipleNewlines = regexp.MustCompile(`\n{3,}`)
)

func CleanText(text string) string {
	text = strings.TrimSpace(text)

	// 3'ten fazla satır boşluğunu 2'ye indir (paragraf ayrımı kalsın)
	text = multipleNewlines.ReplaceAllString(text, "\n\n")

	// Control karakterlerini temizle ama Markdown için önemli olanlara dokunma
	text = removeControlChars(text)

	return text
}

func removeControlChars(s string) string {
	var builder strings.Builder
	for _, r := range s {
		// Tab (\t) ve Newline (\n) Markdown tabloları ve listeleri için gereklidir, silme!
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
	// Kelime ortasından bölmemeye çalış
	if lastSpace := strings.LastIndex(truncated, " "); lastSpace > 0 {
		truncated = truncated[:lastSpace]
	}

	return truncated + "..."
}

func WordCount(text string) int {
	return len(strings.Fields(text))
}
