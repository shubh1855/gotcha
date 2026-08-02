# Gotcha

[![CI](https://github.com/shubh1855/Gotcha/actions/workflows/ci.yml/badge.svg)](https://github.com/shubh1855/Gotcha/actions/workflows/ci.yml)
[![Latest Release](https://img.shields.io/github/v/release/shubh1855/Gotcha)](https://github.com/shubh1855/Gotcha/releases/latest)
[![Downloads](https://img.shields.io/github/downloads/shubh1855/Gotcha/total)](https://github.com/shubh1855/Gotcha/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/shubh1855/Gotcha)](https://github.com/shubh1855/Gotcha/blob/main/go.mod)
[![License](https://img.shields.io/github/license/shubh1855/Gotcha)](LICENSE)

A fast, concurrent web crawler written in Go that recursively crawls websites, respects `robots.txt`, extracts structured page information, and exports the results as a deterministic JSON report.

---

## Motivation

To learn about concurrency in Golang and learn how websites show information to better learn about how to handle concurrent tasks.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
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
- Configurable User-Agent
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

## Usage / Quick Start

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

Format the code:

```bash
gofmt -w .
```

Build:

```bash
go build .
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

- [ ] Code is formatted with `gofmt`
- [ ] `go vet ./...` passes
- [ ] `go test ./...` passes
- [ ] Documentation has been updated if required
- [ ] `CHANGELOG.md` has been updated for user-facing changes

### Reporting Issues

If you encounter a bug or have a feature request, please open a GitHub Issue with:

- A clear description of the problem
- Steps to reproduce (for bugs)
- Expected behavior
- Relevant logs or screenshots (if applicable)

---

## Roadmap

### v0.6.0

- [ ] Better handling of redirects
- [ ] Request rate limiting
- [ ] Crawl depth limiting
- [ ] Smarter duplicate detection
- [ ] Sitemap generation
- [ ] Benchmark suite
- [ ] Additional export formats (CSV/Markdown)
- [ ] `golangci-lint` integration
- [ ] Increased unit and integration test coverage

---

## Documentation

- [CHANGELOG](CHANGELOG.md)
- [Releases](https://github.com/shubh1855/Gotcha/releases)
- [LICENSE](LICENSE)

---

## License

This project is licensed under the **GNU General Public License v3.0**. See the [LICENSE](LICENSE,,) file for details.
