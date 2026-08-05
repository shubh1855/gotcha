package main

import (
	"net/http"
	"net/url"
	"sync"
	"time"
)

type config struct {
	pages              map[string]PageData
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
	stats              CrawlerStats
	userAgent          string
	robotsRules        []string
	maxDepth           int
	requestDelay       time.Duration
	maxRedirects       int
	client             *http.Client
}

func (cfg *config) addPageVisit(normalizedURL string) bool {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	if _, ok := cfg.pages[normalizedURL]; ok {
		return false
	}

	cfg.pages[normalizedURL] = PageData{}
	return true
}
