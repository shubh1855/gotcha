package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var httpClient = &http.Client{
	Timeout: 15 * time.Second,
}

func getHTML(rawURL string) (string, error) {
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
