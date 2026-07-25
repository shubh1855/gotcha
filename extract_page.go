package main

import (
	"net/url"
)

type PageData struct {
	URL            string `json:"url"`
	Heading        string `json:"heading"`
	FirstParagraph string `json:"first_paragraph"`

	InternalLinks []string `json:"internal_links"`
	ExternalLinks []string `json:"external_links"`

	ImageURLs []string `json:"image_urls"`
}

func extractPageData(html, pageURL string) PageData {
	baseURL, err := url.Parse(pageURL)
	if err != nil {
		return PageData{
			URL: pageURL,
		}
	}

	discoveredLinks, _ := getURLsFromHTML(html, baseURL)
	imageURLs, _ := getImagesFromHTML(html, baseURL)

	return PageData{
		URL:            pageURL,
		Heading:        getHeadingFromHTML(html),
		FirstParagraph: getFirstParagraphFromHTML(html),

		// Populated with all discovered links; classified by the crawler
		InternalLinks: discoveredLinks,
		ExternalLinks: []string{},

		ImageURLs: imageURLs,
	}
}
