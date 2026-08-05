package main

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	maxRetries  = 3
	baseBackoff = 500 * time.Millisecond
)

func (cfg *config) httpClient() *http.Client {
	return &http.Client{
		Timeout: 15 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= cfg.maxRedirects {
				return fmt.Errorf(
					"stopped after %d redirects",
					cfg.maxRedirects,
				)
			}

			logger.Debug(
				"following redirect",
				slog.String("from", via[len(via)-1].URL.String()),
				slog.String("to", req.URL.String()),
			)

			return nil
		},
	}
}

func backoff(attempt int) time.Duration {
	return baseBackoff * time.Duration(1<<attempt)
}

func shouldRetry(err error) bool {
	if statusErr, ok := errors.AsType[*HTTPStatusError](err); ok {
		switch statusErr.Code {
		case http.StatusTooManyRequests,
			http.StatusInternalServerError,
			http.StatusBadGateway,
			http.StatusServiceUnavailable,
			http.StatusGatewayTimeout:
			return true
		}
	}

	if netErr, ok := errors.AsType[net.Error](err); ok {
		return netErr.Timeout()
	}

	return false
}

// waitBeforeRequest applies the configured delay before issuing
// each HTTP request.
func (cfg *config) waitBeforeRequest() {
	if cfg.requestDelay <= 0 {
		return
	}

	logger.Debug(
		"waiting before request",
		slog.Duration("delay", cfg.requestDelay),
	)

	time.Sleep(cfg.requestDelay)
}

func (cfg *config) fetchHTML(rawURL string) (string, error) {
	cfg.waitBeforeRequest()

	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return "", fmt.Errorf("create request for %q: %w", rawURL, err)
	}

	req.Header.Set("User-Agent", cfg.userAgent)

	logger.Debug(
		"sending HTTP request",
		slog.String("url", rawURL),
		slog.String("user_agent", cfg.userAgent),
	)

	client := cfg.httpClient()
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("GET %q: %w", rawURL, err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			logger.Warn(
				"failed to close response body",
				slog.String("url", rawURL),
				slog.Any("error", err),
			)
		}
	}()

	logger.Debug(
		"received HTTP response",
		slog.String("url", rawURL),
		slog.Int("status", resp.StatusCode),
	)

	if resp.StatusCode != http.StatusOK {
		return "", &HTTPStatusError{
			URL:    rawURL,
			Code:   resp.StatusCode,
			Status: resp.Status,
		}
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/html") {
		return "", &ContentTypeError{
			URL:         rawURL,
			ContentType: contentType,
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response body for %q: %w", rawURL, err)
	}

	logger.Debug(
		"downloaded page",
		slog.String("url", rawURL),
		slog.Int("bytes", len(body)),
	)

	return string(body), nil
}

func (cfg *config) getHTML(rawURL string) (string, error) {
	var err error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		html, err := cfg.fetchHTML(rawURL)
		if err == nil {
			return html, nil
		}

		if !shouldRetry(err) {
			return "", err
		}

		if attempt == maxRetries {
			break
		}

		delay := backoff(attempt)

		logger.Warn(
			"retrying request",
			slog.String("url", rawURL),
			slog.Int("attempt", attempt+1),
			slog.Int("max_attempts", maxRetries),
			slog.Duration("delay", delay),
			slog.Any("error", err),
		)

		time.Sleep(delay)
	}

	logger.Error(
		"request failed",
		slog.String("url", rawURL),
		slog.Int("retries", maxRetries),
		slog.Any("error", err),
	)

	return "", fmt.Errorf(
		"GET %q: failed after %d retries: %w",
		rawURL,
		maxRetries,
		err,
	)
}
