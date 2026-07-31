package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	maxRetries  = 3
	baseBackoff = 500 * time.Millisecond
)

var httpClient = &http.Client{
	Timeout: 15 * time.Second,
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

func (cfg *config) fetchHTML(rawURL string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return "", fmt.Errorf("create request for %q: %w", rawURL, err)
	}

	req.Header.Set("User-Agent", cfg.userAgent)

	logger.Debug(
		"sending HTTP request",
		"url", rawURL,
		"user_agent", cfg.userAgent,
	)

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("GET %q: %w", rawURL, err)
	}
	defer resp.Body.Close()

	logger.Debug(
		"received HTTP response",
		"url", rawURL,
		"status", resp.StatusCode,
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
		"url", rawURL,
		"bytes", len(body),
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
			"url", rawURL,
			"attempt", attempt+1,
			"max_attempts", maxRetries,
			"delay", delay,
			"error", err,
		)

		time.Sleep(delay)
	}

	logger.Error(
		"request failed",
		"url", rawURL,
		"retries", maxRetries,
		"error", err,
	)

	return "", fmt.Errorf(
		"GET %q: failed after %d retries: %w",
		rawURL,
		maxRetries,
		err,
	)
}
