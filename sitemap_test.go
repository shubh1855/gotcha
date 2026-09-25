package main

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateSitemap(t *testing.T) {
	pages := map[string]PageData{
		"https://example.com/about": {},
		"https://example.com/":      {},
		"https://example.com/blog":  {},
	}

	data, err := generateSitemap(pages)
	if err != nil {
		t.Fatalf("generateSitemap() error = %v", err)
	}

	var document struct {
		XMLName xml.Name `xml:"urlset"`
		XMLNS   string   `xml:"xmlns,attr"`
		URLs    []struct {
			Loc string `xml:"loc"`
		} `xml:"url"`
	}

	if err := xml.Unmarshal(data, &document); err != nil {
		t.Fatalf("generated sitemap is not valid XML: %v", err)
	}

	if document.XMLNS != "http://www.sitemaps.org/schemas/sitemap/0.9" {
		t.Errorf("XML namespace = %q, want sitemap namespace", document.XMLNS)
	}

	if len(document.URLs) != len(pages) {
		t.Fatalf("URL count = %d, want %d", len(document.URLs), len(pages))
	}
}

func TestGenerateSitemap_DeterministicOrder(t *testing.T) {
	pages := map[string]PageData{
		"https://example.com/z": {},
		"https://example.com/a": {},
		"https://example.com/m": {},
	}

	data, err := generateSitemap(pages)
	if err != nil {
		t.Fatalf("generateSitemap() error = %v", err)
	}

	var document struct {
		URLs []struct {
			Loc string `xml:"loc"`
		} `xml:"url"`
	}

	if err := xml.Unmarshal(data, &document); err != nil {
		t.Fatalf("generated sitemap is not valid XML: %v", err)
	}

	want := []string{
		"https://example.com/a",
		"https://example.com/m",
		"https://example.com/z",
	}

	for i, wantURL := range want {
		if document.URLs[i].Loc != wantURL {
			t.Errorf("URL[%d] = %q, want %q", i, document.URLs[i].Loc, wantURL)
		}
	}
}

func TestGenerateSitemap_EmptyPages(t *testing.T) {
	data, err := generateSitemap(map[string]PageData{})
	if err != nil {
		t.Fatalf("generateSitemap() error = %v", err)
	}

	var document struct {
		XMLName xml.Name `xml:"urlset"`
		URLs    []struct {
			Loc string `xml:"loc"`
		} `xml:"url"`
	}

	if err := xml.Unmarshal(data, &document); err != nil {
		t.Fatalf("generated sitemap is not valid XML: %v", err)
	}

	if len(document.URLs) != 0 {
		t.Errorf("URL count = %d, want 0", len(document.URLs))
	}
}

func TestGenerateSitemap_XMLEscaping(t *testing.T) {
	pages := map[string]PageData{
		"https://example.com/search?q=go&sort=desc": {},
	}

	data, err := generateSitemap(pages)
	if err != nil {
		t.Fatalf("generateSitemap() error = %v", err)
	}

	if !strings.Contains(string(data), "&amp;") {
		t.Error("generated sitemap does not XML-escape '&'")
	}
}

func TestWriteSitemap(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "sitemap.xml")

	pages := map[string]PageData{
		"https://example.com/":      {},
		"https://example.com/about": {},
	}

	if err := writeSitemap(pages, filename); err != nil {
		t.Fatalf("writeSitemap() error = %v", err)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("os.ReadFile() error = %v", err)
	}

	if len(data) == 0 {
		t.Fatal("sitemap file is empty")
	}

	var document struct {
		URLs []struct {
			Loc string `xml:"loc"`
		} `xml:"url"`
	}

	if err := xml.Unmarshal(data, &document); err != nil {
		t.Fatalf("written sitemap is not valid XML: %v", err)
	}

	if len(document.URLs) != len(pages) {
		t.Errorf("URL count = %d, want %d", len(document.URLs), len(pages))
	}
}
