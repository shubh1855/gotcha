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

func (cfg *config) crawlPage(rawCurrentURL string) {
	defer cfg.wg.Done()
	defer func() {
		<-cfg.concurrencyControl
	}()

	if !isSameDomain(cfg.baseURL.String(), rawCurrentURL) {
		return
	}

	// Normalize the URL to map pages with same key
	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("error normalizing %q: %v\n", rawCurrentURL, err)
		return
	}

	if !cfg.addPageVisit(normalizedURL) {
		return
	}

	fmt.Printf("crawling: %s\n", rawCurrentURL)

	html, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Println(err)
		return
	}

	pageData := extractPageData(html, rawCurrentURL)

	cfg.mu.Lock()
	cfg.pages[normalizedURL] = pageData
	cfg.mu.Unlock()

	for _, link := range pageData.OutgoingLinks {
		cfg.wg.Add(1)

		go func(link string) {
			cfg.concurrencyControl <- struct{}{}
			cfg.crawlPage(link)
		}(link)
	}
}
