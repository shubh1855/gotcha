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
