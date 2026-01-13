package transport

import (
	"context"
	"net/http"
)

type HTTPClient interface {
	Fetch(ctx context.Context, url string) (*http.Response, error)
}
