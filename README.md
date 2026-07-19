# Gotcha

A concurrent web crawler written in Go that recursively crawls websites, extracts structured page data, and generates a JSON report.

## Features

- Concurrent crawling with configurable worker limits
- Recursive traversal of pages within the same domain
- URL normalization to avoid duplicate crawls
- Configurable maximum page limit
- Extracts:
  - Page URL
  - Main heading
  - First paragraph
  - Outgoing links
  - Image URLs
- Generates a structured `report.json`
- Thread-safe crawling using goroutines, mutexes, and wait groups

## Requirements

- Go 1.26 or newer

## Installation

Clone the repository:

```bash
git clone https://github.com/<your-username>/Gotcha.git
cd Gotcha
```

Download dependencies:

```bash
go mod download
```

## Usage

Basic usage:

```bash
go run . <url>
```

Example:

```bash
go run . https://learnwebscraping.dev/practice/ecommerce/
```

Specify maximum concurrency:

```bash
go run . <url> 10
```

Specify both maximum concurrency and maximum pages:

```bash
go run . <url> 10 200
```

You can also build it using the:

```bash
go build -o gotcha
```

Arguments:

| Argument        | Description                        | Default  |
| --------------- | ---------------------------------- | -------- |
| URL             | Website to crawl                   | Required |
| Max Concurrency | Number of concurrent crawl workers | 5        |
| Max Pages       | Maximum pages to crawl             | 100      |

## Output

After the crawl completes, a `report.json` file is generated in the project directory.

Example:

```json
[
  {
    "url": "https://example.com/",
    "heading": "Example Domain",
    "first_paragraph": "This domain is for use in illustrative examples.",
    "outgoing_links": ["https://www.iana.org/domains/example"],
    "image_urls": []
  }
]
```

## Project Structure

```
.
├── config.go
├── crawler.go
├── extract_content.go
├── extract_page.go
├── fetch.go
├── json_report.go
├── main.go
├── normalize_url.go
└── parser.go
```

## How it Works

1. Starts from the provided URL.
2. Crawls pages concurrently while respecting the configured concurrency limit.
3. Ignores pages outside the starting domain.
4. Normalizes URLs to avoid revisiting the same page.
5. Extracts structured page information.
6. Stores results in memory.
7. Exports the crawl as a sorted JSON report.

## Notes

- Only pages within the starting domain are crawled.
- Duplicate URLs are skipped.
- Crawling stops when the configured page limit is reached.
