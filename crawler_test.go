package main

import "testing"

func TestShouldCrawlChildren(t *testing.T) {
	tests := []struct {
		name     string
		maxDepth int
		depth    int
		want     bool
	}{
		{
			name:     "unlimited depth",
			maxDepth: -1,
			depth:    100,
			want:     true,
		},
		{
			name:     "root may crawl children",
			maxDepth: 0,
			depth:    0,
			want:     false,
		},
		{
			name:     "depth below limit",
			maxDepth: 2,
			depth:    1,
			want:     true,
		},
		{
			name:     "depth equals limit",
			maxDepth: 2,
			depth:    2,
			want:     false,
		},
		{
			name:     "depth exceeds limit",
			maxDepth: 2,
			depth:    3,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config{
				maxDepth: tt.maxDepth,
			}

			got := cfg.shouldCrawlChildren(tt.depth)

			if got != tt.want {
				t.Fatalf(
					"shouldCrawlChildren(%d) = %v, want %v",
					tt.depth,
					got,
					tt.want,
				)
			}
		})
	}
}
