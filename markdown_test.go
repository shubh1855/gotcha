package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteMarkdownReport(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "crawl.md")

	pages := map[string]PageData{
		"https://example.com/": {
			URL:            "https://example.com/",
			Heading:        "Example Domain",
			FirstParagraph: "This is an example.",
			InternalLinks: []string{
				"https://example.com/about",
			},
			ExternalLinks: []string{
				"https://example.org/",
			},
			ImageURLs: []string{
				"https://example.com/image.png",
			},
		},
	}

	if err := writeMarkdownReport(pages, filename); err != nil {
		t.Fatalf("writeMarkdownReport() error = %v", err)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("os.ReadFile() error = %v", err)
	}

	content := string(data)

	expected := []string{
		"# Crawl Report",
		"## https://example.com/",
		"### Heading",
		"Example Domain",
		"### First Paragraph",
		"This is an example.",
		"### Internal Links",
		"- https://example.com/about",
		"### External Links",
		"- https://example.org/",
		"### Images",
		"- https://example.com/image.png",
	}

	for _, want := range expected {
		if !strings.Contains(content, want) {
			t.Errorf("Markdown output does not contain %q", want)
		}
	}
}

func TestWriteMarkdownReport_EmptyPages(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "crawl.md")

	if err := writeMarkdownReport(map[string]PageData{}, filename); err != nil {
		t.Fatalf("writeMarkdownReport() error = %v", err)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("os.ReadFile() error = %v", err)
	}

	want := "# Crawl Report\n"

	if string(data) != want {
		t.Errorf("output = %q, want %q", string(data), want)
	}
}

func TestWriteMarkdownReport_DeterministicOrder(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "crawl.md")

	pages := map[string]PageData{
		"https://example.com/z": {
			URL:     "https://example.com/z",
			Heading: "Z Page",
		},
		"https://example.com/a": {
			URL:     "https://example.com/a",
			Heading: "A Page",
		},
		"https://example.com/m": {
			URL:     "https://example.com/m",
			Heading: "M Page",
		},
	}

	if err := writeMarkdownReport(pages, filename); err != nil {
		t.Fatalf("writeMarkdownReport() error = %v", err)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("os.ReadFile() error = %v", err)
	}

	content := string(data)

	aIndex := strings.Index(content, "## https://example.com/a")
	mIndex := strings.Index(content, "## https://example.com/m")
	zIndex := strings.Index(content, "## https://example.com/z")

	if aIndex == -1 || mIndex == -1 || zIndex == -1 {
		t.Fatal("one or more expected URLs missing from Markdown output")
	}

	if aIndex >= mIndex || mIndex >= zIndex {
		t.Errorf(
			"URLs are not sorted: a=%d, m=%d, z=%d",
			aIndex,
			mIndex,
			zIndex,
		)
	}
}

func TestWriteMarkdownReport_EmptyFields(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "crawl.md")

	pages := map[string]PageData{
		"https://example.com/": {
			URL: "https://example.com/",
		},
	}

	if err := writeMarkdownReport(pages, filename); err != nil {
		t.Fatalf("writeMarkdownReport() error = %v", err)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("os.ReadFile() error = %v", err)
	}

	content := string(data)

	if !strings.Contains(content, "### Heading\n_None_") {
		t.Error("missing empty heading marker")
	}

	if !strings.Contains(content, "### First Paragraph\n_None_") {
		t.Error("missing empty paragraph marker")
	}

	if !strings.Contains(content, "### Internal Links\n_None_") {
		t.Error("missing empty internal links marker")
	}

	if !strings.Contains(content, "### External Links\n_None_") {
		t.Error("missing empty external links marker")
	}

	if !strings.Contains(content, "### Images\n_None_") {
		t.Error("missing empty images marker")
	}
}
