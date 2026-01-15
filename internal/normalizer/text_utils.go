package normalizer

import (
	"regexp"
	"strings"
)

var (
	multipleSpaces   = regexp.MustCompile(`[ \t]+`)
	multipleNewlines = regexp.MustCompile(`\n{2,}`)

	specialChars = regexp.MustCompile(`[^\w\s\-.,!?]`)
)

type TextUtils struct{}

func NewTextUtils() *TextUtils {
	return &TextUtils{}
}

func (t *TextUtils) CollapseWhitespace(text string) string {

	text = multipleSpaces.ReplaceAllString(text, " ")

	text = multipleNewlines.ReplaceAllString(text, "\n")

	return text
}

func (t *TextUtils) RemoveExtraWhitespace(text string) string {

	text = strings.TrimSpace(text)

	text = t.CollapseWhitespace(text)

	return text
}

func (t *TextUtils) ToLower(text string) string {
	return strings.ToLower(text)
}

func (t *TextUtils) RemoveSpecialChars(text string) string {

	return specialChars.ReplaceAllString(text, "")
}

func (t *TextUtils) RemovePunctuation(text string) string {
	punctuation := regexp.MustCompile(`[.,!?;:'"()\[\]{}<>]`)
	return punctuation.ReplaceAllString(text, "")
}

func (t *TextUtils) SplitIntoWords(text string) []string {
	return strings.Fields(text)
}

func (t *TextUtils) WordCount(text string) int {
	return len(strings.Fields(text))
}

func (t *TextUtils) TrimToWordBoundary(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}

	truncated := text[:maxLen]
	if lastSpace := strings.LastIndex(truncated, " "); lastSpace > 0 {
		truncated = truncated[:lastSpace]
	}

	return truncated
}
