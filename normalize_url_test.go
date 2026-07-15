package main

import "testing"

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		expected string
		wantErr  bool
	}{
		{
			name:     "remove https scheme",
			inputURL: "https://www.boot.dev/blog/path",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "remove http scheme",
			inputURL: "http://www.boot.dev/blog/path",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "lowercase host",
			inputURL: "https://WWW.BOOT.DEV/blog/path",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "lowercase path",
			inputURL: "https://www.boot.dev/BLOG/PATH",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "remove trailing slash",
			inputURL: "https://www.boot.dev/blog/path/",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "root url trailing slash",
			inputURL: "https://www.boot.dev/",
			expected: "www.boot.dev",
		},
		{
			name:     "query parameters ignored",
			inputURL: "https://www.boot.dev/blog/path?sort=desc",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "fragment ignored",
			inputURL: "https://www.boot.dev/blog/path#section1",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "query and fragment ignored",
			inputURL: "https://www.boot.dev/blog/path/?id=1#top",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "no path",
			inputURL: "https://www.boot.dev",
			expected: "www.boot.dev",
		},
		{
			name:     "invalid url",
			inputURL: "://invalid-url",
			wantErr:  true,
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := normalizeURL(tc.inputURL)

			if tc.wantErr {
				if err == nil {
					t.Errorf("Test %d - %s FAIL: expected error", i, tc.name)
				}
				return
			}

			if err != nil {
				t.Errorf("Test %d - %s FAIL: unexpected error: %v", i, tc.name, err)
				return
			}

			if actual != tc.expected {
				t.Errorf("Test %d - %s FAIL: expected %q, got %q", i, tc.name, tc.expected, actual)
			}
		})
	}
}
