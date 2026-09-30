# Gotcha

[![CI](https://github.com/shubh1855/Gotcha/actions/workflows/ci.yml/badge.svg)](https://github.com/shubh1855/Gotcha/actions/workflows/ci.yml)
[![Lint](https://github.com/shubh1855/Gotcha/actions/workflows/lint.yml/badge.svg)](https://github.com/shubh1855/Gotcha/actions/workflows/lint.yml)
[![Latest Release](https://img.shields.io/github/v/release/shubh1855/Gotcha)](https://github.com/shubh1855/Gotcha/releases/latest)
[![Downloads](https://img.shields.io/github/downloads/shubh1855/Gotcha/total)](https://github.com/shubh1855/Gotcha/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/shubh1855/Gotcha)](https://github.com/shubh1855/Gotcha/blob/main/go.mod)
[![License](https://img.shields.io/github/license/shubh1855/Gotcha)](LICENSE)

A fast, concurrent web crawler written in Go that recursively crawls websites, respects `robots.txt`, supports configurable redirects and rate limiting, extracts structured page information, and exports deterministic reports in JSON, CSV, Markdown, and XML sitemap formats.

---

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Example Output](#example-output)
- [Output](#output)
- [Project Structure](#project-structure)
- [How It Works](#how-it-works)
- [Development](#development)
- [Contributing](#contributing)
- [Roadmap](#roadmap)
- [Documentation](#documentation)
- [License](#license)

---

## Features

- Recursive crawling within a single domain
- Configurable concurrency limit
- Configurable maximum page limit
- Configurable crawl depth
- Configurable request delay
- Configurable User-Agent
- Configurable request rate limiting
- Configurable HTTP redirect handling
- URL normalization to avoid duplicate crawls
- Automatic retry with exponential backoff for transient HTTP failures
- `robots.txt` support
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
- XML sitemap generation with `--sitemap`
- CSV report generation with `--csv`
- Markdown report generation with `--markdown`
- Deterministic output ordering across report formats
- Benchmark suite for URL normalization and page extraction
- Integration tests for crawler behavior
- Race-condition testing with Go's race detector

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
gotcha [flags] <url>
```

Show all available options:

```bash
gotcha --help
```

### Examples

Crawl a website using the default configuration:

```bash
gotcha https://example.com
```

Limit concurrency:

```bash
gotcha --concurrency 10 https://example.com
```

Limit the number of pages:

```bash
gotcha --pages 200 https://example.com
```

Limit crawl depth:

```bash
gotcha --depth 3 https://example.com
```

Add a delay between requests:

```bash
gotcha --delay 500ms https://example.com
```

Rate limit HTTP requests:

```bash
gotcha --rate 5 https://example.com
```

Use a custom User-Agent:

```bash
gotcha --user-agent "MyCrawler/1.0" https://example.com
```

Limit HTTP redirects:

```bash
gotcha --max-redirects 5 https://example.com
```

Enable verbose logging:

```bash
gotcha --verbose https://example.com
```

Generate an XML sitemap:

```bash
gotcha --sitemap https://example.com
```

Generate a CSV report:

```bash
gotcha --csv https://example.com
```

Generate a Markdown report:

```bash
gotcha --markdown https://example.com
```

Generate all available reports:

```bash
gotcha \
    --sitemap \
    --csv \
    --markdown \
    https://example.com
```

Print version information:

```bash
gotcha --version
```

Example using multiple options:

```bash
gotcha \
    --concurrency 10 \
    --pages 200 \
    --depth 3 \
    --delay 500ms \
    --rate 5 \
    --max-redirects 5 \
    --user-agent "MyCrawler/1.0" \
    --sitemap \
    --csv \
    --markdown \
    --verbose \
    https://example.com
```

### Available Flags

| Flag                  | Description                                                   | Default            |
| --------------------- | ------------------------------------------------------------- | ------------------ |
| `-c, --concurrency`   | Maximum concurrent requests                                   | `5`                |
| `-p, --pages`         | Maximum pages to crawl                                        | `100`              |
| `--depth`             | Maximum crawl depth (`-1` = unlimited)                        | `-1`               |
| `-u, --user-agent`    | HTTP User-Agent                                               | `Gotcha/1.0 (...)` |
| `-d, --delay`         | Fixed delay between requests                                  | `0s`               |
| `-r, --max-redirects` | Maximum redirects to follow                                   | `10`               |
| `--rate`              | Maximum HTTP requests per second (`0` disables rate limiting) | `0`                |
| `--sitemap`           | Generate `sitemap.xml`                                        | `false`            |
| `--csv`               | Generate `crawl.csv`                                          | `false`            |
| `--markdown`          | Generate `crawl.md`                                           | `false`            |
| `-v, --verbose`       | Enable debug logging                                          | `false`            |
| `-V, --version`       | Print version information                                     | `false`            |

---

## Example Output

During a crawl, Gotcha reports crawl progress and statistics:

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

When optional exports are enabled, Gotcha also reports the generated files:

```text
JSON report written to report.json
Sitemap written to sitemap.xml
CSV report written to crawl.csv
Markdown report written to crawl.md
```

---

## Output

Every crawl generates a deterministic JSON report:

```text
report.json
```

Additional report formats can be generated using the corresponding flags.

| Flag         | Output        | Description                          |
| ------------ | ------------- | ------------------------------------ |
| `--sitemap`  | `sitemap.xml` | XML sitemap containing crawled URLs  |
| `--csv`      | `crawl.csv`   | Tabular crawl data                   |
| `--markdown` | `crawl.md`    | Human-readable Markdown crawl report |

For example:

```bash
gotcha \
    --sitemap \
    --csv \
    --markdown \
    https://example.com
```

Generates:

```text
report.json
sitemap.xml
crawl.csv
crawl.md
```

### JSON

The JSON report contains structured information about each crawled page:

```json
{
  "url": "https://example.com",
  "heading": "Example Domain",
  "first_paragraph": "This domain is for use in illustrative examples in documents.",
  "internal_links": [],
  "external_links": [],
  "image_urls": []
}
```

### CSV

The CSV report contains the following fields:

- URL
- Heading
- First paragraph
- Internal links
- External links
- Image URLs

Lists are stored as pipe-separated values within CSV fields.

### Markdown

The Markdown report provides a human-readable representation of the crawl:

```text
crawl.md
```

Each crawled page contains:

- URL
- Heading
- First paragraph
- Internal links
- External links
- Images

### Sitemap

The sitemap contains the URLs successfully crawled by Gotcha and is written as a standard XML sitemap:

```text
sitemap.xml
```

Output ordering is deterministic across the supported report formats.

---

## robots.txt

Gotcha automatically checks `/robots.txt` before crawling.

The current implementation supports:

- `User-agent: *`
- `Disallow`

Other directives such as `Allow`, `Crawl-delay`, and `Sitemap` are currently ignored.

If a website does not provide a `robots.txt` file, Gotcha proceeds without robots.txt restrictions.

---

## Retries

Gotcha automatically retries transient failures, including:

- Request timeouts
- HTTP `429`
- HTTP `500`
- HTTP `502`
- HTTP `503`
- HTTP `504`

Retries use exponential backoff:

```text
500ms
1s
2s
```

---

## Crawl Statistics

Gotcha tracks:

- Pages crawled
- Failed fetches
- Skipped pages
- Internal links discovered
- External links discovered
- Pages skipped by `robots.txt`

---

## Project Structure

```text
.
├── .github/
│   └── workflows/
│       ├── ci.yml
│       ├── lint.yml
│       └── release.yml
├── CHANGELOG.md
├── LICENSE
├── README.md
├── benchmark_test.go
├── config.go
├── crawler.go
├── crawler_test.go
├── csv.go
├── csv_test.go
├── errors.go
├── extract_content.go
├── extract_page.go
├── extract_page_test.go
├── fetch.go
├── fetch_test.go
├── integration_test.go
├── json_report.go
├── logger.go
├── main.go
├── markdown.go
├── markdown_test.go
├── normalize_url.go
├── normalize_url_test.go
├── parser.go
├── parser_test.go
├── robots.go
├── robots_test.go
├── sitemap.go
├── sitemap_test.go
├── stats.go
├── url_test.go
└── version.go
```

---

## How It Works

1. Start crawling from the provided URL.
2. Check the site's `robots.txt` rules.
3. Normalize URLs to prevent duplicate visits.
4. Ensure URLs belong to the same domain.
5. Crawl pages recursively while respecting concurrency, page, and depth limits.
6. Apply request rate limiting and optional request delays.
7. Retry transient HTTP failures using exponential backoff.
8. Follow HTTP redirects up to the configured redirect limit.
9. Extract structured page information.
10. Classify internal and external links.
11. Store page data in memory.
12. Generate deterministic JSON output.
13. Generate optional CSV, Markdown, and XML sitemap reports.

---

## Development

Run the test suite:

```bash
go test ./... -v
```

Run static analysis:

```bash
go vet ./...
```

Run linting:

```bash
golangci-lint run
```

Run the race detector:

```bash
go test -race ./...
```

Run benchmarks:

```bash
go test -bench=. -benchmem
```

Format the code:

```bash
gofmt -w .
```

Build:

```bash
go build .
```

### Quality Checks

Before submitting changes, run:

```bash
gofmt -w .
go test ./... -v
go vet ./...
golangci-lint run
go test -race ./...
```

---

## Contributing

Contributions are welcome! Whether it's a bug fix, a new feature, documentation improvements, or performance optimizations, your help is appreciated.

### Getting Started

Clone the repository and install dependencies:

```bash
git clone https://github.com/shubh1855/Gotcha.git
cd Gotcha
go mod download
```

### Development Workflow

1. Create a new branch from `main`.

```bash
git checkout -b feature/my-feature
```

2. Make your changes.

3. Run the quality checks.

```bash
gofmt -w .
go vet ./...
go test ./...
golangci-lint run
go test -race ./...
```

4. Commit your changes using a descriptive commit message.

Examples:

```text
feat(crawler): add crawl depth support
fix(fetch): handle redirect loops
docs: update README
```

5. Push your branch and open a Pull Request.

### Pull Request Checklist

Before opening a Pull Request, please ensure:

- [x] Code is formatted with `gofmt`
- [x] `go vet ./...` passes
- [x] `go test ./...` passes
- [x] `golangci-lint run` passes
- [x] `go test -race ./...` passes
- [x] Documentation has been updated if required
- [x] `CHANGELOG.md` has been updated for user-facing changes

### Reporting Issues

If you encounter a bug or have a feature request, please open a GitHub Issue with:

- A clear description of the problem
- Steps to reproduce (for bugs)
- Expected behavior
- Relevant logs or screenshots, if applicable

---

## Roadmap

### v0.9.0

- Improved `robots.txt` support
- Include/exclude URL patterns
- Canonical URL handling
- Richer page metadata
- Improved crawl error reporting
- Expanded crawl statistics
- Improved graceful shutdown

### v1.0.0

- Stable crawler API
- Context-aware crawling
- Graceful cancellation
- SSRF and network security hardening
- Fuzz testing
- Large-crawl testing
- Performance and memory optimization
- Stable CLI and report formats

---

## Documentation

- [CHANGELOG](CHANGELOG.md)
- [Releases](https://github.com/shubh1855/Gotcha/releases)

---

## License

This project is licensed under the **GNU General Public License v3.0**. See the [LICENSE](LICENSE) file for details.
