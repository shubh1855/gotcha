package main

import (
	"encoding/csv"
	"os"
	"sort"
	"strings"
)

const (
	csvColumnURL            = "url"
	csvColumnHeading        = "heading"
	csvColumnFirstParagraph = "first_paragraph"
	csvColumnInternalLinks  = "internal_links"
	csvColumnExternalLinks  = "external_links"
	csvColumnImageURLs      = "image_urls"
)

func writeCSVReport(pages map[string]PageData, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	writer := csv.NewWriter(file)

	if err := writer.Write([]string{
		csvColumnURL,
		csvColumnHeading,
		csvColumnFirstParagraph,
		csvColumnInternalLinks,
		csvColumnExternalLinks,
		csvColumnImageURLs,
	}); err != nil {
		return err
	}

	urls := make([]string, 0, len(pages))

	for pageURL := range pages {
		urls = append(urls, pageURL)
	}

	sort.Strings(urls)

	for _, pageURL := range urls {
		page := pages[pageURL]

		record := []string{
			page.URL,
			page.Heading,
			page.FirstParagraph,
			strings.Join(page.InternalLinks, "|"),
			strings.Join(page.ExternalLinks, "|"),
			strings.Join(page.ImageURLs, "|"),
		}

		if err := writer.Write(record); err != nil {
			return err
		}
	}

	writer.Flush()

	return writer.Error()
}
