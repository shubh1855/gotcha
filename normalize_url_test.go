package main

import "testing"

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "remove fragment",
			input: "https://example.com/about#team",
			want:  "https://example.com/about",
		},
		{
			name:  "remove trailing slash",
			input: "https://example.com/about/",
			want:  "https://example.com/about",
		},
		{
			name:  "preserve root slash",
			input: "https://example.com/",
			want:  "https://example.com/",
		},
		{
			name:  "remove utm parameter",
			input: "https://example.com/about?utm_source=github",
			want:  "https://example.com/about",
		},
		{
			name:  "remove multiple tracking parameters",
			input: "https://example.com/about?utm_source=github&fbclid=123&gclid=456",
			want:  "https://example.com/about",
		},
		{
			name:  "preserve meaningful query parameter",
			input: "https://example.com/product?id=123",
			want:  "https://example.com/product?id=123",
		},
		{
			name:  "remove tracking while preserving meaningful query",
			input: "https://example.com/product?utm_source=github&id=123",
			want:  "https://example.com/product?id=123",
		},
		{
			name:  "sort query parameters",
			input: "https://example.com/search?z=3&a=1&m=2",
			want:  "https://example.com/search?a=1&m=2&z=3",
		},
		{
			name:  "normalize scheme and host",
			input: "HTTPS://EXAMPLE.COM/About",
			want:  "https://example.com/About",
		},
		{
			name:    "invalid URL",
			input:   "://invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizeURL(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Errorf("normalizeURL() = %q, want %q", got, tt.want)
			}
		})
	}
}
