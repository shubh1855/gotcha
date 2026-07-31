package main

import (
	"bufio"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
)

// parses and extracts disallow rules for user-agent
func parseRobotsTxt(content string) []string {
	rules := []string{}

	scanner := bufio.NewScanner(strings.NewReader(content))

	inGlobalAgent := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// skip blank lines or comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		lower := strings.ToLower(line)

		switch {
		case strings.HasPrefix(lower, "user-agent:"):
			agent := strings.TrimSpace(line[len("user-agent:"):])
			inGlobalAgent = agent == "*"

		case inGlobalAgent && strings.HasPrefix(lower, "disallow:"):
			path := strings.TrimSpace(line[len("Disallow:"):])

			// empty disallow means everything is allowed
			if path == "" {
				continue
			}

			rules = append(rules, path)
		}
	}

	return rules
}

func (cfg *config) loadRobotsTxt() error {
	robotsURL := cfg.baseURL.ResolveReference(&url.URL{
		Path: "/robots.txt",
	})

	content, err := cfg.getHTML(robotsURL.String())
	if err != nil {
		var statusErr *HTTPStatusError

		if errors.As(err, &statusErr) &&
			statusErr.Code == http.StatusNotFound {
			cfg.robotsRules = []string{}
			return nil
		}

		return err
	}

	cfg.robotsRules = parseRobotsTxt(content)

	logger.Info(
		"loaded robots.txt",
		slog.Int("rules", len(cfg.robotsRules)),
	)

	return nil
}

func (cfg *config) isAllowed(rawURL string) bool {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	path := parsedURL.Path

	for _, rule := range cfg.robotsRules {
		// Disallow //
		if rule == "/" {
			return false
		}

		// Exact match
		if path == rule {
			return false
		}

		// Child paths
		if strings.HasPrefix(path, rule+"/") {
			return false
		}
	}

	return true
}
