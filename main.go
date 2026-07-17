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
		fmt.Printf("starting crawl of site: %s\n", args[0])

	default:
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}
}
