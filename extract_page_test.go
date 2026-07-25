package main

import (
	"reflect"
	"testing"
)

func TestExtractPageData(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		html     string
		expected PageData
	}{
		{
			name: "complete page",
			url:  "https://crawler-test.com",
			html: `<html><body>
				<h1>Test Title</h1>
				<p>This is the first paragraph.</p>
				<a href="/link1">Link 1</a>
				<img src="/image1.jpg">
			</body></html>`,
			expected: PageData{
				URL:            "https://crawler-test.com",
				Heading:        "Test Title",
				FirstParagraph: "This is the first paragraph.",
				InternalLinks: []string{
					"https://crawler-test.com/link1",
				},
				ExternalLinks: []string{},
				ImageURLs: []string{
					"https://crawler-test.com/image1.jpg",
				},
			},
		},
		{
			name: "relative and absolute resources",
			url:  "https://crawler-test.com",
			html: `<html><body>
				<h2>Heading</h2>
				<main>
					<p>Main paragraph.</p>
				</main>

				<a href="about">About</a>
				<a href="https://google.com">Google</a>

				<img src="/logo.png">
				<img src="https://cdn.com/banner.png">
			</body></html>`,
			expected: PageData{
				URL:            "https://crawler-test.com",
				Heading:        "Heading",
				FirstParagraph: "Main paragraph.",
				InternalLinks: []string{
					"https://crawler-test.com/about",
					"https://google.com",
				},
				ExternalLinks: []string{},
				ImageURLs: []string{
					"https://crawler-test.com/logo.png",
					"https://cdn.com/banner.png",
				},
			},
		},
		{
			name: "minimal page",
			url:  "https://crawler-test.com",
			html: `<html><body></body></html>`,
			expected: PageData{
				URL:            "https://crawler-test.com",
				Heading:        "",
				FirstParagraph: "",
				InternalLinks:  []string{},
				ExternalLinks:  []string{},
				ImageURLs:      []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := extractPageData(tt.html, tt.url)

			if !reflect.DeepEqual(actual, tt.expected) {
				// t.Errorf("\nexpected:\n%+v\n\nactual:\n%+v",
				// 	tt.expected,
				// 	actual,
				//)
				t.Errorf(
					"expected InternalLinks nil? %v len=%d\nactual InternalLinks nil? %v len=%d\n"+
						"expected ExternalLinks nil? %v len=%d\nactual ExternalLinks nil? %v len=%d\n"+
						"expected ImageURLs nil? %v len=%d\nactual ImageURLs nil? %v len=%d",
					tt.expected.InternalLinks == nil, len(tt.expected.InternalLinks),
					actual.InternalLinks == nil, len(actual.InternalLinks),
					tt.expected.ExternalLinks == nil, len(tt.expected.ExternalLinks),
					actual.ExternalLinks == nil, len(actual.ExternalLinks),
					tt.expected.ImageURLs == nil, len(tt.expected.ImageURLs),
					actual.ImageURLs == nil, len(actual.ImageURLs),
				)
			}
		})
	}
}
