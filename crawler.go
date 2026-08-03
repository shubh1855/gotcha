package main

import (
	"log/slog"
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

func (cfg *config) crawlPage(rawCurrentURL string, depth int) {
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

	if !cfg.isAllowed(rawCurrentURL) {
		cfg.incrementRobotsSkipped()
		logger.Warn(
			"skipping page",
			slog.String("url", rawCurrentURL),
			slog.String("reason", "robots.txt"))
		return
	}

	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		logger.Warn(
			"failed to normalize URL",
			slog.String("url", rawCurrentURL),
			slog.Any("error", err),
		)
		cfg.incrementSkippedPages()
		return
	}

	if !cfg.addPageVisit(normalizedURL) {
		cfg.incrementSkippedPages()
		return
	}

	logger.Info(
		"crawling page",
		slog.String("url", rawCurrentURL),
	)

	html, err := cfg.getHTML(rawCurrentURL)
	if err != nil {
		logger.Warn(
			"failed to fetch page",
			slog.String("url", rawCurrentURL),
			slog.Any("error", err),
		)
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

	// stop here since we reached the max crawl depth
	if cfg.maxDepth >= 0 && depth >= cfg.maxDepth {
		return
	}

	for _, link := range pageData.InternalLinks {
		cfg.wg.Add(1)

		go func(link string) {
			cfg.concurrencyControl <- struct{}{}
			cfg.crawlPage(link, depth+1)
		}(link)
	}
}
