package main

import "testing"

func TestGetHeadingFromHTML(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "h1 exists",
			html:     "<html><body><h1>Test Title</h1></body></html>",
			expected: "Test Title",
		},
		{
			name:     "fallback to h2",
			html:     "<html><body><h2>Heading Two</h2></body></html>",
			expected: "Heading Two",
		},
		{
			name:     "prefer h1 over h2",
			html:     "<html><body><h2>Second</h2><h1>First</h1></body></html>",
			expected: "First",
		},
		{
			name:     "no headings",
			html:     "<html><body><p>Hello</p></body></html>",
			expected: "",
		},
		{
			name:     "trim whitespace",
			html:     "<html><body><h1>   Boot.dev   </h1></body></html>",
			expected: "Boot.dev",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := getHeadingFromHTML(tt.html)
			if actual != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, actual)
			}
		})
	}
}

func TestGetFirstParagraphFromHTML(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name: "prefer paragraph inside main",
			html: `<html><body>
				<p>Outside paragraph.</p>
				<main>
					<p>Main paragraph.</p>
				</main>
			</body></html>`,
			expected: "Main paragraph.",
		},
		{
			name: "fallback to first paragraph",
			html: `<html><body>
				<p>First paragraph.</p>
				<p>Second paragraph.</p>
			</body></html>`,
			expected: "First paragraph.",
		},
		{
			name: "first paragraph inside main",
			html: `<html><body>
				<main>
					<p>First.</p>
					<p>Second.</p>
				</main>
			</body></html>`,
			expected: "First.",
		},
		{
			name:     "no paragraph",
			html:     "<html><body><h1>Title</h1></body></html>",
			expected: "",
		},
		{
			name: "trim whitespace",
			html: `<html><body>
				<p>   Hello World   </p>
			</body></html>`,
			expected: "Hello World",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := getFirstParagraphFromHTML(tt.html)
			if actual != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, actual)
			}
		})
	}
}
