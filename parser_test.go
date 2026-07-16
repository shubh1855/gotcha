package main

import (
	"net/url"
	"reflect"
	"testing"
)

func TestGetURLsFromHTML(t *testing.T) {
	base, _ := url.Parse("https://crawler-test.com")

	tests := []struct {
		name     string
		html     string
		expected []string
	}{
		{
			name: "absolute url",
			html: `<a href="https://google.com">Google</a>`,
			expected: []string{
				"https://google.com",
			},
		},
		{
			name: "relative url",
			html: `<a href="/about">About</a>`,
			expected: []string{
				"https://crawler-test.com/about",
			},
		},
		{
			name: "multiple urls",
			html: `
				<a href="/one">One</a>
				<a href="https://google.com">Google</a>
				<a href="about">About</a>
			`,
			expected: []string{
				"https://crawler-test.com/one",
				"https://google.com",
				"https://crawler-test.com/about",
			},
		},
		{
			name:     "missing href",
			html:     `<a>Broken</a>`,
			expected: []string{},
		},
		{
			name:     "empty href",
			html:     `<a href="">Empty</a>`,
			expected: []string{},
		},
		{
			name:     "no links",
			html:     `<p>Hello</p>`,
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := getURLsFromHTML(tt.html, base)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("expected %v got %v", tt.expected, actual)
			}
		})
	}
}

func TestGetImagesFromHTML(t *testing.T) {
	base, _ := url.Parse("https://crawler-test.com")

	tests := []struct {
		name     string
		html     string
		expected []string
	}{
		{
			name: "absolute image",
			html: `<img src="https://cdn.com/logo.png">`,
			expected: []string{
				"https://cdn.com/logo.png",
			},
		},
		{
			name: "relative image",
			html: `<img src="/logo.png">`,
			expected: []string{
				"https://crawler-test.com/logo.png",
			},
		},
		{
			name: "multiple images",
			html: `
				<img src="/one.png">
				<img src="two.png">
			`,
			expected: []string{
				"https://crawler-test.com/one.png",
				"https://crawler-test.com/two.png",
			},
		},
		{
			name:     "missing src",
			html:     `<img alt="logo">`,
			expected: []string{},
		},
		{
			name:     "empty src",
			html:     `<img src="">`,
			expected: []string{},
		},
		{
			name:     "no images",
			html:     `<p>Hello</p>`,
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := getImagesFromHTML(tt.html, base)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("expected %v got %v", tt.expected, actual)
			}
		})
	}
}
