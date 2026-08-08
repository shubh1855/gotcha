package main

import (
	"net/url"
	"strings"
)

func normalizeURL(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	// Fragments are client-side and do not identify a different resource.
	parsedURL.Fragment = ""

	// Normalize scheme and hostname.
	parsedURL.Scheme = strings.ToLower(parsedURL.Scheme)
	parsedURL.Host = strings.ToLower(parsedURL.Host)

	// Normalize the path.
	if parsedURL.Path == "" {
		parsedURL.Path = "/"
	}

	if parsedURL.Path != "/" {
		parsedURL.Path = strings.TrimRight(parsedURL.Path, "/")
	}

	// Remove known tracking parameters while preserving
	// meaningful query parameters.
	query := parsedURL.Query()

	for key := range query {
		lowerKey := strings.ToLower(key)

		if strings.HasPrefix(lowerKey, "utm_") ||
			lowerKey == "fbclid" ||
			lowerKey == "gclid" ||
			lowerKey == "msclkid" {
			query.Del(key)
		}
	}

	// url.Values.Encode() sorts query parameters by key
	// and produces a deterministic representation.
	parsedURL.RawQuery = query.Encode()

	return parsedURL.String(), nil
}
