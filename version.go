package main

import "fmt"

var (
	Version = "dev"
	Commit  = "unknown"
	Date    = "unknown"
)

func VersionString() string {
	return fmt.Sprintf(
		"Gotcha %s\nCommit: %s\nBuilt: %s",
		Version,
		Commit,
		Date,
	)
}
