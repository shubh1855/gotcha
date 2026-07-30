package main

import (
	"reflect"
	"testing"
)

func TestParseRobotsTxt(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name: "single disallow",
			input: `
User-agent: *
Disallow: /admin
`,
			expected: []string{"/admin"},
		},
		{
			name: "multiple disallow",
			input: `
User-agent: *
Disallow: /admin
Disallow: /private
`,
			expected: []string{
				"/admin",
				"/private",
			},
		},
		{
			name: "ignore other user agents",
			input: `
User-agent: Googlebot
Disallow: /google

User-agent: *
Disallow: /admin
`,
			expected: []string{
				"/admin",
			},
		},
		{
			name: "empty disallow",
			input: `
User-agent: *
Disallow:
`,
			expected: []string{},
		},
		{
			name: "comments ignored",
			input: `
# robots
User-agent: *
# hidden
Disallow: /secret
`,
			expected: []string{
				"/secret",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := parseRobotsTxt(tt.input)

			if !reflect.DeepEqual(actual, tt.expected) {
				t.Fatalf("expected %v, got %v", tt.expected, actual)
			}
		})
	}
}

func TestIsAllowed(t *testing.T) {
	tests := []struct {
		name    string
		rules   []string
		url     string
		allowed bool
	}{
		{
			name:    "allowed page",
			rules:   []string{"/admin", "/private"},
			url:     "https://example.com/about",
			allowed: true,
		},
		{
			name:    "exact match",
			rules:   []string{"/admin", "/private"},
			url:     "https://example.com/admin",
			allowed: false,
		},
		{
			name:    "child path",
			rules:   []string{"/admin", "/private"},
			url:     "https://example.com/admin/users",
			allowed: false,
		},
		{
			name:    "private child",
			rules:   []string{"/admin", "/private"},
			url:     "https://example.com/private/data",
			allowed: false,
		},
		{
			name:    "administrator should be allowed",
			rules:   []string{"/admin", "/private"},
			url:     "https://example.com/administrator",
			allowed: true,
		},
		{
			name:    "root disallow blocks everything",
			rules:   []string{"/"},
			url:     "https://example.com/anything",
			allowed: false,
		},
		{
			name:    "root itself blocked",
			rules:   []string{"/"},
			url:     "https://example.com/",
			allowed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config{
				robotsRules: tt.rules,
			}

			if got := cfg.isAllowed(tt.url); got != tt.allowed {
				t.Errorf("isAllowed(%q) = %v, want %v", tt.url, got, tt.allowed)
			}
		})
	}
}
