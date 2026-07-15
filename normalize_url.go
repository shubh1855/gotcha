package main

import (
	"net/url"
	"strings"
)

func normalizeURL(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	host := strings.ToLower(u.Host)
	path := strings.ToLower(u.Path)

	return strings.TrimRight(host+path, "/"), nil
}
