package main

import (
	"fmt"
	"os"
	"sort"
)

func writeMarkdownReport(pages map[string]PageData, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	urls := make([]string, 0, len(pages))

	for pageURL := range pages {
		urls = append(urls, pageURL)
	}

	sort.Strings(urls)

	if _, err := fmt.Fprintln(file, "# Crawl Report"); err != nil {
		return err
	}

	for _, pageURL := range urls {
		page := pages[pageURL]

		if _, err := fmt.Fprintf(file, "\n## %s\n\n", page.URL); err != nil {
			return err
		}

		if _, err := fmt.Fprintln(file, "### Heading"); err != nil {
			return err
		}

		if page.Heading != "" {
			if _, err := fmt.Fprintf(file, "%s\n\n", page.Heading); err != nil {
				return err
			}
		} else {
			if _, err := fmt.Fprintln(file, "_None_"); err != nil {
				return err
			}
		}

		if _, err := fmt.Fprintln(file, "### First Paragraph"); err != nil {
			return err
		}

		if page.FirstParagraph != "" {
			if _, err := fmt.Fprintf(file, "%s\n\n", page.FirstParagraph); err != nil {
				return err
			}
		} else {
			if _, err := fmt.Fprintln(file, "_None_"); err != nil {
				return err
			}
		}

		if _, err := fmt.Fprintln(file, "### Internal Links"); err != nil {
			return err
		}

		if err := writeMarkdownList(file, page.InternalLinks); err != nil {
			return err
		}

		if _, err := fmt.Fprintln(file, "### External Links"); err != nil {
			return err
		}

		if err := writeMarkdownList(file, page.ExternalLinks); err != nil {
			return err
		}

		if _, err := fmt.Fprintln(file, "### Images"); err != nil {
			return err
		}

		if err := writeMarkdownList(file, page.ImageURLs); err != nil {
			return err
		}
	}

	return nil
}

func writeMarkdownList(file *os.File, values []string) error {
	if len(values) == 0 {
		if _, err := fmt.Fprintln(file, "_None_"); err != nil {
			return err
		}

		_, err := fmt.Fprintln(file)
		return err
	}

	values = append([]string(nil), values...)
	sort.Strings(values)

	for _, value := range values {
		if _, err := fmt.Fprintf(file, "- %s\n", value); err != nil {
			return err
		}
	}

	_, err := fmt.Fprintln(file)
	return err
}
