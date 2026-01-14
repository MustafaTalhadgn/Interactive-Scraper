package validation

import "fmt"

type ValidationError struct {
	Reason  string
	Details string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s - %s", e.Reason, e.Details)
}

var (
	ErrInvalidStatus = func(status int) *ValidationError {
		return &ValidationError{
			Reason:  "Geçersiz durum",
			Details: fmt.Sprintf("HTTP Durum Kodu %d geçersizdir", status),
		}
	}
	ErrInvalidContentType = func(contentType string) *ValidationError {
		return &ValidationError{
			Reason:  "Geçersiz içerik türü",
			Details: fmt.Sprintf("İçerik türü '%s' geçersizdir", contentType),
		}
	}

	ErrBodyTooLarge = func(size int64) *ValidationError {
		return &ValidationError{
			Reason:  "İçerik çok büyük",
			Details: fmt.Sprintf("İçerik boyutu %d bayttan büyük olamaz", size),
		}
	}

	ErrBodyEmpty = &ValidationError{
		Reason:  "Boş içerik",
		Details: "İçerik boş olamaz",
	}

	ErrInvalidEncoding = func(encoding string) *ValidationError {
		return &ValidationError{
			Reason:  "Geçersiz kodlama",
			Details: "İçerik kodlaması geçersizdir",
		}
	}
)
