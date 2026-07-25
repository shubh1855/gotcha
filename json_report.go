package main

import (
	"encoding/json"
	"os"
	"sort"
)

type CrawlReport struct {
	Summary CrawlerStats `json:"summary"`
	Pages   []PageData   `json:"pages"`
}

func writeJSONReport(
	pages map[string]PageData,
	stats CrawlerStats,
	filename string,
) error {
	keys := make([]string, 0, len(pages))

	for key := range pages {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	report := make([]PageData, 0, len(keys))
	for _, key := range keys {
		report = append(report, pages[key])
	}

	crawlReport := CrawlReport{
		Summary: stats,
		Pages:   report,
	}

	data, err := json.MarshalIndent(crawlReport, "", "")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}
