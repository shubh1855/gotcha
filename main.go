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
)

func main() {
	args := os.Args[1:]

	if len(args) < 1 || len(args) > 3 {
		fmt.Println("Usage: crawler <url> [maxConcurrency] [maxPages]")
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

	baseURL, err := url.Parse(rawURL)
	if err != nil {
		fmt.Printf("invalid URL: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf(
		"Starting crawl of %s (max concurrency: %d, max pages: %d)\n",
		rawURL,
		maxConcurrency,
		maxPages,
	)

	cfg := &config{
		pages:              make(map[string]PageData),
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
		maxPages:           maxPages,
	}

	cfg.wg.Add(1)

	go func() {
		cfg.concurrencyControl <- struct{}{}
		cfg.crawlPage(baseURL.String())
	}()

	cfg.wg.Wait()

	fmt.Println("\nPages crawled:")

	for _, page := range cfg.pages {
		fmt.Printf("%+v\n", page)
	}
}
