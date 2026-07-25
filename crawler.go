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

	cfg.mu.Lock()
	reachedLimit := len(cfg.pages) >= cfg.maxPages
	cfg.mu.Unlock()

	if reachedLimit {
		cfg.incrementSkippedPages()
		return
	}

	if !isSameDomain(cfg.baseURL.String(), rawCurrentURL) {
		cfg.incrementSkippedPages()
		return
	}

	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("warning: failed to normalize %q: %v\n", rawCurrentURL, err)
		cfg.incrementSkippedPages()
		return
	}

	if !cfg.addPageVisit(normalizedURL) {
		cfg.incrementSkippedPages()
		return
	}

	fmt.Printf("crawling: %s\n", rawCurrentURL)

	html, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("warning: %v\n", err)
		cfg.incrementFailedFetches()
		return
	}

	pageData := extractPageData(html, rawCurrentURL)

	pageData.InternalLinks, pageData.ExternalLinks = classifyLinks(
		cfg.baseURL,
		pageData.InternalLinks,
	)

	cfg.incrementInternalLinks(len(pageData.InternalLinks))
	cfg.incrementExternalLinks(len(pageData.ExternalLinks))

	cfg.mu.Lock()
	cfg.pages[normalizedURL] = pageData
	cfg.mu.Unlock()

	cfg.incrementPagesCrawled()

	for _, link := range pageData.InternalLinks {
		cfg.wg.Add(1)

		go func(link string) {
			cfg.concurrencyControl <- struct{}{}
			cfg.crawlPage(link)
		}(link)
	}
}
