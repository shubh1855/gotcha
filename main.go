package main

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"sync"
	"time"

	flag "github.com/spf13/pflag"
)

const (
	defaultMaxConcurrency = 5
	defaultMaxPages       = 100
	defaultUserAgent      = "Gotcha/1.0 (+https://github.com/shubh1855/Gotcha)"
	defaultMaxDepth       = -1
	defaultRequestDelay   = 0 * time.Second
	defaultMaxRedirects   = 10
)

func main() {
	concurrency := flag.IntP(
		"concurrency",
		"c",
		defaultMaxConcurrency,
		"Maximum number of concurrent requests",
	)

	pages := flag.IntP(
		"pages",
		"p",
		defaultMaxPages,
		"Maximum number of pages to crawl",
	)

	userAgent := flag.StringP(
		"user-agent",
		"u",
		defaultUserAgent,
		"HTTP User-Agent",
	)

	verbose := flag.BoolP(
		"verbose",
		"v",
		false,
		"Enable debug logging",
	)

	version := flag.BoolP(
		"version",
		"V",
		false,
		"Print version information",
	)

	depth := flag.Int(
		"depth",
		defaultMaxDepth,
		"Maximum crawl depth (-1 for unlimited)",
	)

	delay := flag.DurationP(
		"delay",
		"d",
		defaultRequestDelay,
		"Delay between HTTP requests",
	)

	redirects := flag.IntP(
		"max-redirects",
		"r",
		defaultMaxRedirects,
		"Maximum number of HTTP redirects")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] <url>\n\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if *version {
		fmt.Println(VersionString())
		return
	}

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	initLogger(*verbose)

	rawURL := flag.Arg(0)

	baseURL, err := url.Parse(rawURL)
	if err != nil {
		logger.Error(
			"invalid URL",
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	logger.Info(
		"starting crawl",
		slog.String("url", baseURL.String()),
		slog.Int("concurrency", *concurrency),
		slog.Int("max_pages", *pages),
		slog.Int("max_depth", *depth),
		slog.Duration("request_delay", *delay),
		slog.Int("max_redirects", *redirects),
		slog.String("user_agent", *userAgent),
	)

	cfg := &config{
		pages:              make(map[string]PageData),
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, *concurrency),
		wg:                 &sync.WaitGroup{},
		maxPages:           *pages,
		userAgent:          *userAgent,
		maxDepth:           *depth,
		requestDelay:       *delay,
		maxRedirects:       *redirects,
	}

	if err := cfg.loadRobotsTxt(); err != nil {
		logger.Warn(
			"failed to load robots.txt",
			slog.Any("error", err),
		)
	}

	cfg.wg.Add(1)

	go func() {
		cfg.concurrencyControl <- struct{}{}
		cfg.crawlPage(baseURL.String(), 0)
	}()

	cfg.wg.Wait()

	logger.Info(
		"crawl complete",
		slog.Int("pages_crawled", cfg.stats.PagesCrawled),
		slog.Int("pages_skipped", cfg.stats.SkippedPages),
		slog.Int("robots_skipped", cfg.stats.RobotsSkipped),
		slog.Int("failed_fetches", cfg.stats.FailedFetches),
		slog.Int("internal_links", cfg.stats.InternalLinks),
		slog.Int("external_links", cfg.stats.ExternalLinks),
	)

	if err := writeJSONReport(cfg.pages, cfg.stats, "report.json"); err != nil {
		logger.Error(
			"failed to write JSON report",
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	logger.Info(
		"JSON report written",
		slog.String("path", "report.json"),
	)
}
