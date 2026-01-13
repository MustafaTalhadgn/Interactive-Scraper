package transport

import "fmt"

type TransportError struct {
	URL        string
	Attempt    int
	StatusCode int
	Err        error
	Retryable  bool
}

func (e *TransportError) Error() string {
	if e.StatusCode > 0 {
		return fmt.Sprintf("transport error: HTTP %d for %s (attempt %d/%d): %v",
			e.StatusCode, e.URL, e.Attempt, e.Attempt, e.Err)
	}
	return fmt.Sprintf("transport error: %s (attempt %d): %v",
		e.URL, e.Attempt, e.Err)
}

func NewTransportError(url string, attempt int, statusCode int, err error, retryable bool) *TransportError {
	return &TransportError{
		URL:        url,
		Attempt:    attempt,
		StatusCode: statusCode,
		Err:        err,
		Retryable:  retryable,
	}
}

func IsRetryableStatus(statusCode int) bool {
	switch statusCode {
	case 408, 429, 500, 502, 503, 504:
		return true
	default:
		return false
	}
}
