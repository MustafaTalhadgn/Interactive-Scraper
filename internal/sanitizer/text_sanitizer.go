package sanitizer

import (
	"regexp"
	"strings"
	"unicode"
)

var (
	sqlInjectionPattern     = regexp.MustCompile(`(?i)(union\s+select|drop\s+table|insert\s+into|delete\s+from|--|;|'|\")`)
	commandInjectionPattern = regexp.MustCompile(`(?i)(&&|\|\||;|` + "`" + `|\$\(|\$\{)`)

	controlCharsPattern = regexp.MustCompile(`[\x00-\x08\x0B\x0C\x0E-\x1F\x7F]`)

	multipleSpaces   = regexp.MustCompile(`[ \t]{2,}`)
	multipleNewlines = regexp.MustCompile(`\n{3,}`)
)

type TextSanitizer struct {
	maxLength int
}

func NewTextSanitizer(maxLength int) *TextSanitizer {
	return &TextSanitizer{
		maxLength: maxLength,
	}
}

func (s *TextSanitizer) Sanitize(text string) string {

	text = s.removeControlCharacters(text)
	text = s.normalizeWhitespace(text)
	text = strings.TrimSpace(text)
	if len(text) > s.maxLength {
		text = s.truncate(text, s.maxLength)
	}

	return text
}

func (s *TextSanitizer) removeControlCharacters(text string) string {

	text = strings.ReplaceAll(text, "\x00", "")
	text = controlCharsPattern.ReplaceAllString(text, "")

	return text
}

func (s *TextSanitizer) normalizeWhitespace(text string) string {
	text = multipleSpaces.ReplaceAllString(text, " ")
	text = multipleNewlines.ReplaceAllString(text, "\n\n")
	return text
}

func (s *TextSanitizer) truncate(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}

	truncated := text[:maxLen]
	if lastSpace := strings.LastIndex(truncated, " "); lastSpace > 0 {
		truncated = truncated[:lastSpace]
	}
	return truncated + "..."
}

func (s *TextSanitizer) ContainsSQLInjection(text string) bool {
	return sqlInjectionPattern.MatchString(text)
}

func (s *TextSanitizer) ContainsCommandInjection(text string) bool {
	return commandInjectionPattern.MatchString(text)
}

func (s *TextSanitizer) EscapeForSQL(text string) string {
	return strings.ReplaceAll(text, "'", "''")
}
func (s *TextSanitizer) RemoveDangerousPatterns(text string) string {
	text = sqlInjectionPattern.ReplaceAllString(text, "[FILTERED]")
	text = commandInjectionPattern.ReplaceAllString(text, "[FILTERED]")

	return text
}
func (s *TextSanitizer) IsASCIIPrintable(text string) bool {
	for _, r := range text {
		if r > unicode.MaxASCII && !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}
