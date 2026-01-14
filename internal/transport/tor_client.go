package transport

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"time"

	"InteractiveScraper/internal/config"

	"golang.org/x/net/proxy"
)

const (
	defaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"
)

type TorClient struct {
	httpClient  *http.Client
	logger      *slog.Logger
	config      *config.Config
	userAgent   string
	rateLimiter *RateLimiter
}

func NewTorClient(cfg *config.Config, logger *slog.Logger) (*TorClient, error) {
	torProxy := cfg.TorProxy
	if torProxy == "" {
		return nil, fmt.Errorf("Tor proxy adresi tanımlı değil")
	}
	dialer, err := proxy.SOCKS5("tcp", torProxy, nil, proxy.Direct)
	if err != nil {
		return nil, fmt.Errorf("Tor dialer oluşturulamadı: %v", err)
	}
	transport := &http.Transport{
		DialContext: dialer.(proxy.ContextDialer).DialContext,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: false,
		DisableKeepAlives:  false,
	}

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   cfg.RequestTimeOut,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {

			if len(via) >= 5 {
				return fmt.Errorf("stopped after 5 redirects")
			}
			return nil
		},
	}

	logger.Info("tor client initialized",
		slog.String("proxy", torProxy),
		slog.Duration("timeout", cfg.RequestTimeOut),
	)

	return &TorClient{
		config:      cfg,
		httpClient:  httpClient,
		logger:      logger,
		userAgent:   defaultUserAgent,
		rateLimiter: NewRateLimiter(cfg.RateLimitRPS),
	}, nil
}

func (t *TorClient) Fetch(ctx context.Context, url string) (*http.Response, error) {

	if err := t.rateLimiter.Wait(ctx); err != nil {
		return nil, NewTransportError(url, 0, 0,
			fmt.Errorf("rate limit wait cancelled: %w", err), false)
	}

	var lastErr error
	for attempt := 1; attempt <= t.config.MaxRetries; attempt++ {
		t.logger.Debug("istek atılıyor",
			slog.String("url", url),
			slog.Int("attempt", attempt),
			slog.Int("max_retries", t.config.MaxRetries),
		)
		resp, err := t.doRequest(ctx, url, attempt)
		if err == nil {
			t.logger.Debug("istek başarılı",
				slog.String("url", url),
				slog.Int("attempt", attempt),
				slog.Int("max_retries", t.config.MaxRetries),
			)
			return resp, nil
		}
		lastErr = err

		if transportErr, ok := err.(*TransportError); ok {
			if !transportErr.Retryable {
				t.logger.Error("non-retryable error",
					slog.String("url", url),
					slog.Int("status", transportErr.StatusCode),
					slog.String("error", err.Error()),
				)
				return nil, err
			}
		}
		t.logger.Warn("istek hatası, tekrar denenecek",
			slog.String("url", url),
			slog.Int("attempt", attempt),
			slog.String("error", err.Error()),
		)
		if attempt < t.config.MaxRetries {
			backoff := t.calculateBackoff(attempt)
			t.logger.Debug("bekleme süresi",
				slog.Duration("backoff", backoff),
				slog.Int("next_attempt", attempt+1),
			)
			select {
			case <-time.After(backoff):
			case <-ctx.Done():
				return nil, NewTransportError(url, attempt, 0, ctx.Err(), false)
			}
		}
	}

	finalErr := NewTransportError(url, t.config.MaxRetries, 0, fmt.Errorf("Maximum deneme sayısına ulaşıldı:%w", lastErr), false)

	t.logger.Error("Maximum deneme sayısına ulaşıldı",
		slog.String("url", url),
		slog.Int("max_retries", t.config.MaxRetries),
	)
	return nil, finalErr
}

func (t *TorClient) doRequest(ctx context.Context, url string, attempt int) (*http.Response, error) {

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, NewTransportError(url, attempt, 0, err, false)
	}

	req.Header.Set("User-Agent", t.userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "keep-alive")

	resp, err := t.httpClient.Do(req)
	if err != nil {

		return nil, NewTransportError(url, attempt, 0, err, true)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		retryable := IsRetryableStatus(resp.StatusCode)
		err := fmt.Errorf("HTTP %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))

		resp.Body.Close()

		return nil, NewTransportError(url, attempt, resp.StatusCode, err, retryable)
	}

	return resp, nil
}
func (t *TorClient) calculateBackoff(attempt int) time.Duration {

	backoff := time.Duration(math.Pow(2, float64(attempt))) * time.Second
	maxBackoff := 5 * time.Second
	if backoff > maxBackoff {
		backoff = maxBackoff
	}
	return backoff
}
