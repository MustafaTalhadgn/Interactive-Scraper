package sanitizer

import (
	"strings"
)

type SQLSanitizer struct{}

func NewSQLSanitizer() *SQLSanitizer {
	return &SQLSanitizer{}
}

func (s *SQLSanitizer) SanitizeForStorage(text string) string {

	text = strings.ReplaceAll(text, "'", "''")
	text = strings.ReplaceAll(text, "\x00", "")
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")

	return text
}

func (s *SQLSanitizer) ValidateForStorage(text string) error {
	if strings.Contains(text, "\x00") {
		return &SanitizationError{
			Reason: "Null byte bulundu",
			Detail: "Metin, izin verilmeyen boş baytlar içeriyor.",
		}
	}

	return nil
}

type SanitizationError struct {
	Reason string
	Detail string
}

func (e *SanitizationError) Error() string {
	return "temizleme hatası: " + e.Reason + " - " + e.Detail
}
