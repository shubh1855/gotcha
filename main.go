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

		html, err := getHTML(baseURL)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(html)

	default:
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}
}
