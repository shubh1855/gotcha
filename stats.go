package main

type CrawlerStats struct {
	PagesCrawled  int
	FailedFetches int
	SkippedPages  int
}

func (cfg *config) incrementPagesCrawled() {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	cfg.stats.PagesCrawled++
}

func (cfg *config) incrementFailedFetches() {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	cfg.stats.FailedFetches++
}

func (cfg *config) incrementSkippedPages() {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	cfg.stats.SkippedPages++
}
