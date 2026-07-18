package main

import (
	"fmt"
	"net/url"
)

func isSameDomain(rawBaseURL, rawCurrentURL string) bool {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return false
	}

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		return false
	}

	return baseURL.Host == currentURL.Host
}

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
	// crawling pages on the same domain
	if !isSameDomain(rawBaseURL, rawCurrentURL) {
		return
	}

	// Normalize the URL to map pages with same key
	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("error normalizing %q: %v\n", rawCurrentURL, err)
		return
	}

	// page visited
	if count, ok := pages[normalizedURL]; ok {
		pages[normalizedURL] = count + 1
		return
	}

	// first visit of page
	pages[normalizedURL] = 1

	fmt.Printf("crawling: %s\n", rawCurrentURL)

	// download the page
	html, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("error fetching %q: %v\n", rawCurrentURL, err)
		return
	}

	baseURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("error parsing current URL %q: %v\n", rawCurrentURL, err)
		return
	}

	// extract outgoing links
	urls, err := getURLsFromHTML(html, baseURL)
	if err != nil {
		fmt.Printf("error extracting URLs from %q: %v\n", rawCurrentURL, err)
		return
	}

	// depth-first crawl
	for _, u := range urls {
		crawlPage(rawBaseURL, u, pages)
	}
}
