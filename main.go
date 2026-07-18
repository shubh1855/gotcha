package main

import (
	"fmt"
	"net/url"
	"os"
	"sync"
)

func main() {
	args := os.Args[1:]

	switch len(args) {
	case 0:
		fmt.Println("no website provided")
		os.Exit(1)

	case 1:
		baseURL, err := url.Parse(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		cfg := &config{
			pages:              make(map[string]PageData),
			baseURL:            baseURL,
			mu:                 &sync.Mutex{},
			concurrencyControl: make(chan struct{}, 5),
			wg:                 &sync.WaitGroup{},
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

	default:
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}
}
