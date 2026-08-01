# Gotcha

[![CI](https://github.com/shubh1855/Gotcha/actions/workflows/ci.yml/badge.svg)](https://github.com/shubh1855/Gotcha/actions/workflows/ci.yml)
[![Latest Release](https://img.shields.io/github/v/release/shubh1855/Gotcha)](https://github.com/shubh1855/Gotcha/releases/latest)
[![Downloads](https://img.shields.io/github/downloads/shubh1855/Gotcha/total)](https://github.com/shubh1855/Gotcha/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/shubh1855/Gotcha)](https://github.com/shubh1855/Gotcha/blob/main/go.mod)
[![License](https://img.shields.io/github/license/shubh1855/Gotcha)](LICENSE)

A fast, concurrent web crawler written in Go that recursively crawls websites, extracts structured page information, and exports the results as a JSON report.

---

## Features

- Recursive crawling within a single domain
- Configurable concurrency limit
- Configurable maximum page limit
- Configurable User-Agent
- URL normalization to avoid duplicate crawls
- Automatic retry with exponential backoff for transient HTTP failures
- robots.txt support
- Internal and external link classification
- Crawl statistics and summary reporting
- Structured logging using `log/slog`
- Thread-safe crawling using goroutines, mutexes, and wait groups
- Structured page extraction including:
  - URL
  - Heading
  - First paragraph
  - Outgoing links
  - Image URLs
- Deterministic JSON report generation

---

## Installation

### Clone the repository

```bash
git clone https://github.com/shubh1855/Gotcha.git
cd Gotcha
```

### Install dependencies

```bash
go mod download
```

### Build

```bash
go build -o gotcha .
```

---

## Usage

Run the crawler:

```bash
go run . <url> [maxConcurrency] [maxPages] [userAgent]
```

### Examples

Crawl using default settings:

```bash
go run . https://example.com
```

Specify maximum concurrency:

```bash
go run . https://example.com 10
```

Specify concurrency and maximum pages:

```bash
go run . https://example.com 10 200
```

Specify a custom User-Agent:

```bash
go run . https://example.com 10 200 "MyCrawler/1.0"
```

### Arguments

| Argument        | Description                      | Default                                             |
| --------------- | -------------------------------- | --------------------------------------------------- |
| URL             | Starting URL                     | **Required**                                        |
| Max Concurrency | Maximum concurrent crawl workers | `5`                                                 |
| Max Pages       | Maximum pages to crawl           | `100`                                               |
| User-Agent      | HTTP User-Agent                  | `Gotcha/1.0 (+https://github.com/shubh1855/Gotcha)` |

---

## Example Output

```text
Starting crawl of https://example.com

Pages crawled : 15
Pages skipped : 2
Robots skipped: 1
Failed fetches: 0
Internal links: 37
External links: 12

JSON report written to report.json
```

---

## Output

After crawling completes, a `report.json` file is generated containing the extracted page information.

Each record contains:

```json
{
  "url": "https://example.com",
  "heading": "Example Domain",
  "first_paragraph": "This domain is for use in illustrative examples in documents.",
  "outgoing_links": [],
  "image_urls": []
}
```

---

## Project Structure

```text
.
├── .github/
│   └── workflows/
│       ├── ci.yml
│       └── release.yml
├── CHANGELOG.md
├── LICENSE
├── README.md
├── config.go
├── crawler.go
├── errors.go
├── extract_content.go
├── extract_page.go
├── fetch.go
├── json_report.go
├── logger.go
├── main.go
├── normalize_url.go
├── parser.go
├── robots.go
├── stats.go
└── version.go
```

---

## How It Works

1. Start crawling from the provided URL.
2. Respect the site's `robots.txt` rules.
3. Normalize URLs to prevent duplicate visits.
4. Crawl pages recursively while remaining within the same domain.
5. Retry transient HTTP failures automatically.
6. Extract structured page information.
7. Classify internal and external links.
8. Store page data in memory.
9. Generate a deterministic JSON report.

---

## Development

Run tests:

```bash
go test ./...
```

Run static analysis:

```bash
go vet ./...
```

Build:

```bash
go build .
```

---

## Roadmap

This project is actively being improved. Planned enhancements include:

- [x] More robust HTTP error handling and retry logic
- [ ] Better handling of redirects, rate limiting, and timeouts
- [x] Tracking internal vs. external links
- [x] Crawl statistics and summary reporting
- [ ] Smarter duplicate detection
- [x] Configurable User-Agent and request headers
- [x] robots.txt support
- [ ] Sitemap generation
- [x] Structured logging
- [ ] Benchmarking and performance improvements
- [ ] Unit and integration test expansion
- [x] GitHub Actions CI
- [x] Automated release workflow

---

## Documentation

- [CHANGELOG](CHANGELOG.md)
- [LICENSE](LICENSE)

---

## License

This project is licensed under the **GNU General Public License v3.0**. See the [LICENSE](LICENSE) file for details.
