package main

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestWriteCSVReport(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "crawl.csv")

	pages := map[string]PageData{
		"https://example.com/z": {
			URL:            "https://example.com/z",
			Heading:        "Z Page",
			FirstParagraph: "This is the Z page.",
			InternalLinks: []string{
				"https://example.com/z/about",
			},
			ExternalLinks: []string{
				"https://external.example/z",
			},
			ImageURLs: []string{
				"https://example.com/z/image.png",
			},
		},
		"https://example.com/a": {
			URL:            "https://example.com/a",
			Heading:        "A Page",
			FirstParagraph: "This is the A page.",
			InternalLinks: []string{
				"https://example.com/a/about",
				"https://example.com/a/contact",
			},
		},
	}

	if err := writeCSVReport(pages, filename); err != nil {
		t.Fatalf("writeCSVReport() error = %v", err)
	}

	file, err := os.Open(filename)
	if err != nil {
		t.Fatalf("os.Open() error = %v", err)
	}
	defer func() {
		_ = file.Close()
	}()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("csv.ReadAll() error = %v", err)
	}

	if len(records) != 3 {
		t.Fatalf("record count = %d, want 3", len(records))
	}

	wantHeader := []string{
		"url",
		"heading",
		"first_paragraph",
		"internal_links",
		"external_links",
		"image_urls",
	}

	if !reflect.DeepEqual(records[0], wantHeader) {
		t.Errorf("header = %v, want %v", records[0], wantHeader)
	}

	if records[1][0] != "https://example.com/a" {
		t.Errorf("first data URL = %q, want https://example.com/a", records[1][0])
	}

	if records[2][0] != "https://example.com/z" {
		t.Errorf("second data URL = %q, want https://example.com/z", records[2][0])
	}
}

func TestWriteCSVReport_EmptyPages(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "crawl.csv")

	if err := writeCSVReport(map[string]PageData{}, filename); err != nil {
		t.Fatalf("writeCSVReport() error = %v", err)
	}

	file, err := os.Open(filename)
	if err != nil {
		t.Fatalf("os.Open() error = %v", err)
	}
	defer func() {
		_ = file.Close()
	}()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("csv.ReadAll() error = %v", err)
	}

	if len(records) != 1 {
		t.Fatalf("record count = %d, want 1", len(records))
	}
}

func TestWriteCSVReport_EscapesCSVFields(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "crawl.csv")

	pages := map[string]PageData{
		"https://example.com/": {
			URL:            "https://example.com/",
			Heading:        `A heading, with "quotes"`,
			FirstParagraph: "First line\nSecond line",
		},
	}

	if err := writeCSVReport(pages, filename); err != nil {
		t.Fatalf("writeCSVReport() error = %v", err)
	}

	file, err := os.Open(filename)
	if err != nil {
		t.Fatalf("os.Open() error = %v", err)
	}
	defer func() {
		_ = file.Close()
	}()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("csv.ReadAll() error = %v", err)
	}

	if records[1][1] != `A heading, with "quotes"` {
		t.Errorf("heading = %q, want original value", records[1][1])
	}

	if records[1][2] != "First line\nSecond line" {
		t.Errorf("paragraph = %q, want original value", records[1][2])
	}
}
