package main

import (
	"encoding/xml"
	"os"
	"sort"
)

type sitemapURL struct {
	Loc string `xml:"loc"`
}

type sitemap struct {
	XMLName xml.Name     `xml:"urlset"`
	XMLNS   string       `xml:"xmlns,attr"`
	URLs    []sitemapURL `xml:"url"`
}

func generateSitemap(pages map[string]PageData) ([]byte, error) {
	urls := make([]string, 0, len(pages))

	for pageURL := range pages {
		urls = append(urls, pageURL)
	}

	sort.Strings(urls)

	sitemapData := sitemap{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  make([]sitemapURL, 0, len(urls)),
	}

	for _, pageURL := range urls {
		sitemapData.URLs = append(sitemapData.URLs, sitemapURL{
			Loc: pageURL,
		})
	}

	data, err := xml.MarshalIndent(sitemapData, "", "  ")
	if err != nil {
		return nil, err
	}

	data = append([]byte(xml.Header), data...)

	return data, nil
}

func writeSitemap(pages map[string]PageData, filename string) error {
	data, err := generateSitemap(pages)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0o644)
}
