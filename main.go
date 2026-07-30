package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"sync"
)

const (
	defaultMaxConcurrency = 5
	defaultMaxPages       = 100
	defaultUserAgent      = "Gotcha/1.0 (+https://github.com/shubh1855/Gotcha)"
)

func main() {
	args := os.Args[1:]

	if len(args) < 1 || len(args) > 4 {
		fmt.Println("Usage: crawler <url> [maxConcurrency] [maxPages] [userAgent]")
		os.Exit(1)
	}

	rawURL := args[0]
	maxConcurrency := defaultMaxConcurrency
	maxPages := defaultMaxPages

	if len(args) >= 2 {
		var err error
		maxConcurrency, err = strconv.Atoi(args[1])
		if err != nil || maxConcurrency <= 0 {
			fmt.Println("maxConcurrency must be a positive integer")
			os.Exit(1)
		}
	}

	if len(args) >= 3 {
		var err error
		maxPages, err = strconv.Atoi(args[2])
		if err != nil || maxPages <= 0 {
			fmt.Println("maxPages must be a positive integer")
			os.Exit(1)
		}
	}

	userAgent := defaultUserAgent
	if len(args) >= 4 {
		userAgent = args[3]
	}

	baseURL, err := url.Parse(rawURL)
	if err != nil {
		fmt.Printf("invalid URL: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf(
		"Starting crawl of %s (max concurrency: %d, max pages: %d)\n",
		baseURL.String(),
		maxConcurrency,
		maxPages,
	)
	fmt.Printf("User-Agent: %s\n", userAgent)

	cfg := &config{
		pages:              make(map[string]PageData),
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
		maxPages:           maxPages,
		userAgent:          userAgent,
	}

	if err := cfg.loadRobotsTxt(); err != nil {
		fmt.Printf("Warning: failed to load robots.txt: %v\n", err)
	}

	cfg.wg.Add(1)

	go func() {
		cfg.concurrencyControl <- struct{}{}
		cfg.crawlPage(baseURL.String())
	}()

	cfg.wg.Wait()

	fmt.Println("\nCrawl complete")
	fmt.Printf("Pages crawled : %d\n", cfg.stats.PagesCrawled)
	fmt.Printf("Pages skipped : %d\n", cfg.stats.SkippedPages)
	fmt.Printf("Robots skipped: %d\n", cfg.stats.RobotsSkipped)
	fmt.Printf("Failed fetches: %d\n", cfg.stats.FailedFetches)
	fmt.Printf("Internal links: %d\n", cfg.stats.InternalLinks)
	fmt.Printf("External links: %d\n", cfg.stats.ExternalLinks)

	if err := writeJSONReport(
		cfg.pages,
		cfg.stats,
		"report.json",
	); err != nil {
		fmt.Printf("failed to write report: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("JSON report written to report.json")
}
