package main

import (
	"testing"
)

func BenchmarkNormalizeURL(b *testing.B) {
	rawURL := "HTTPS://Example.COM/blog/article/?utm_source=test&foo=bar#section"

	b.ResetTimer()

	for i := 0; b.Loop(); i++ {
		_, _ = normalizeURL(rawURL)
	}
}

func BenchmarkExtractPageData(b *testing.B) {
	html := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Example Page</title>
	</head>
	<body>
		<h1>Example Heading</h1>
		<p>This is an example paragraph containing some text.</p>
		<a href="/about">About</a>
		<a href="/contact">Contact</a>
		<a href="https://external.example/">External</a>
		<img src="/image.png">
	</body>
	</html>
	`

	pageURL := "https://example.com/"

	b.ResetTimer()

	for i := 0; b.Loop(); i++ {
		_ = extractPageData(html, pageURL)
	}
}
