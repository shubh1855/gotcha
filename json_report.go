package main

import (
	"encoding/json"
	"os"
	"sort"
)

func writeJSONReport(pages map[string]PageData, filename string) error {
	keys := make([]string, 0, len(pages))

	for key := range pages {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	report := make([]PageData, 0, len(keys))
	for _, key := range keys {
		report = append(report, pages[key])
	}

	data, err := json.MarshalIndent(report, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}
