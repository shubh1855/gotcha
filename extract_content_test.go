package main

import "testing"

func TestGetHeadingFromHTML(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "single h1",
			html:     "<h1>Hello</h1>",
			expected: "Hello",
		},
		{
			name:     "fallback h2",
			html:     "<h2>Hello</h2>",
			expected: "Hello",
		},
		{
			name:     "prefer h1 over h2",
			html:     "<h2>Second</h2><h1>First</h1>",
			expected: "First",
		}, {name: "multiple h1",
			html:     "<h1>One</h1><h1>Two</h1>",
			expected: "One",
		},
		{
			name:     "nested formatting",
			html:     "<h1>Hello <strong>World</strong></h1>",
			expected: "Hello World",
		},
		{
			name:     "trim whitespace",
			html:     "<h1>   Hello   </h1>",
			expected: "Hello",
		},
		{
			name:     "empty heading",
			html:     "<h1></h1>",
			expected: "",
		},
		{
			name:     "uppercase tag",
			html:     "<H1>Hello</H1>",
			expected: "Hello",
		},
		{
			name:     "malformed html",
			html:     "<h1>Hello",
			expected: "Hello",
		},
		{
			name:     "no headings",
			html:     "<p>Hello</p>",
			expected: "",
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
			name:     "single paragraph",
			html:     "<p>Hello</p>",
			expected: "Hello",
		},
		{
			name:     "multiple paragraphs",
			html:     "<p>One</p><p>Two</p>",
			expected: "One",
		},
		{
			name: "prefer paragraph inside main",
			html: `
				<p>Outside</p>
				<main>
					<p>Inside</p>
				</main>
			`,
			expected: "Inside",
		},
		{
			name: "nested paragraph inside article",
			html: `
				<main>
					<article>
						<p>Nested</p>
					</article>
				</main>
			`,
			expected: "Nested",
		},
		{
			name:     "nested formatting",
			html:     "<p>Hello <em>World</em></p>",
			expected: "Hello World",
		},
		{
			name:     "trim whitespace",
			html:     "<p>   Hello World   </p>",
			expected: "Hello World",
		},
		{
			name:     "empty paragraph",
			html:     "<p></p>",
			expected: "",
		},
		{
			name:     "uppercase tag",
			html:     "<P>Hello</P>",
			expected: "Hello",
		},
		{
			name:     "malformed html",
			html:     "<p>Hello",
			expected: "Hello",
		},
		{
			name:     "no paragraph",
			html:     "<div>Hello</div>",
			expected: "",
		},
		{
			name: "main exists but no paragraph",
			html: `
				<p>Outside</p>
				<main>
					<div>No paragraph</div>
				</main>
			`,
			expected: "Outside",
		},
		{
			name: "multiple main tags",
			html: `
				<main><p>First</p></main>
				<main><p>Second</p></main>
			`,
			expected: "First",
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
