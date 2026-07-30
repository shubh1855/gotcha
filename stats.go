package main

type CrawlerStats struct {
	PagesCrawled  int
	FailedFetches int
	SkippedPages  int

	InternalLinks int
	ExternalLinks int

	RobotsSkipped int
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

func (cfg *config) incrementInternalLinks(n int) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	cfg.stats.InternalLinks += n
}

func (cfg *config) incrementExternalLinks(n int) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	cfg.stats.ExternalLinks += n
}

func (cfg *config) incrementRobotsSkipped() {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	cfg.stats.RobotsSkipped++
}
