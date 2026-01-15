package sanitizer

import (
	"log/slog"
)

type Sanitizer struct {
	htmlSanitizer *HTMLSanitizer
	textSanitizer *TextSanitizer
	sqlSanitizer  *SQLSanitizer
	logger        *slog.Logger
}

type Config struct {
	MaxTextLength int
	StripHTML     bool
	LogDangerous  bool
}

func DefaultConfig() *Config {
	return &Config{
		MaxTextLength: 50000,
		StripHTML:     true,
		LogDangerous:  true,
	}
}

func NewSanitizer(config *Config, logger *slog.Logger) *Sanitizer {
	if config == nil {
		config = DefaultConfig()
	}
	return &Sanitizer{
		htmlSanitizer: NewHTMLSanitizer(),
		textSanitizer: NewTextSanitizer(config.MaxTextLength),
		sqlSanitizer:  NewSQLSanitizer(),
		logger:        logger,
	}
}

type SanitizedContent struct {
	PlainText    string
	SafeHTML     string
	SQLSafe      string
	WasDangerous bool
}

func (s *Sanitizer) Sanitize(html string) *SanitizedContent {
	result := &SanitizedContent{}

	result.PlainText = s.htmlSanitizer.SanitizeToPlainText(html)
	result.SafeHTML = s.htmlSanitizer.SanitizeToSafeHTML(html)
	result.PlainText = s.textSanitizer.Sanitize(result.PlainText)

	if s.textSanitizer.ContainsSQLInjection(result.PlainText) {
		s.logger.Warn("SQL injection algılandı",
			slog.String("pattern", "sql_injection"),
		)
		result.WasDangerous = true
		result.PlainText = s.textSanitizer.RemoveDangerousPatterns(result.PlainText)
	}

	if s.textSanitizer.ContainsCommandInjection(result.PlainText) {
		s.logger.Warn("Command injection algılandı",
			slog.String("pattern", "command_injection"),
		)
		result.WasDangerous = true
		result.PlainText = s.textSanitizer.RemoveDangerousPatterns(result.PlainText)
	}

	result.SQLSafe = s.sqlSanitizer.SanitizeForStorage(result.PlainText)

	return result
}

func (s *Sanitizer) SanitizeTitle(title string) string {

	title = s.htmlSanitizer.SanitizeToPlainText(title)

	title = s.textSanitizer.normalizeWhitespace(title)

	if len(title) > 200 {
		title = s.textSanitizer.truncate(title, 200)
	}
	return title
}

func (s *Sanitizer) SanitizeContent(content string) string {
	sanitized := s.Sanitize(content)
	return sanitized.PlainText
}

func (s *Sanitizer) ValidateForStorage(content string) error {
	return s.sqlSanitizer.ValidateForStorage(content)
}
