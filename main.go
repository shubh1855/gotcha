package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]

	switch len(args) {
	case 0:
		fmt.Println("no website provided")
		os.Exit(1)

	case 1:
		baseURL := args[0]

		fmt.Printf("starting crawl of site: %s\n", args[0])

		pages := make(map[string]int)

		crawlPage(baseURL, baseURL, pages)

		fmt.Println()
		fmt.Println("Pages crawled")

		for page, count := range pages {
			fmt.Printf("%3d %s\n", count, page)
		}

	default:
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}
}
