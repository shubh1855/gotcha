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

func fetchHTML(rawURL string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return "", fmt.Errorf("create request for %q: %w", rawURL, err)
	}

	req.Header.Set("User-Agent", "Gotcha/1.0 (+https://github.com/shubh1855/Gotcha)")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("GET %q: %w", rawURL, err)
	}
	defer resp.Body.Close()

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

	return string(body), nil
}

func getHTML(rawURL string) (string, error) {
	var err error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		html, err := fetchHTML(rawURL)
		if err == nil {
			return html, nil
		}

		if !shouldRetry(err) {
			return "", err
		}

		if attempt == maxRetries {
			break
		}

		time.Sleep(backoff(attempt))
	}

	return "", fmt.Errorf(
		"GET %q: failed after %d retries: %w",
		rawURL,
		maxRetries,
		err)
}
