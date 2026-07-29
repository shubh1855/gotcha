package main

import (
	"bufio"
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
