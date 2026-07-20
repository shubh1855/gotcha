# Gotcha

A fast, concurrent web crawler written in Go that recursively crawls websites, extracts structured page information, and exports the results as a JSON report.

## Features

- Recursive crawling within a single domain
- Configurable concurrency limit
- Configurable maximum page limit
- URL normalization to avoid duplicate crawls
- Thread-safe crawling using goroutines, mutexes, and wait groups
- Structured page extraction including:
  - URL
  - Heading
  - First paragraph
  - Outgoing links
  - Image URLs

- Deterministic JSON report generation

## Installation

Clone the repository:

```bash
git clone https://github.com/<your-username>/Gotcha.git
cd Gotcha
```

Install dependencies:

```bash
go mod download
```

## Usage

Run the crawler:

```bash
go run . <url>
```

Example:

```bash
go run . https://learnwebscraping.dev/practice/ecommerce/
```

Specify concurrency:

```bash
go run . <url> 10
```

Specify both concurrency and maximum pages:

```bash
go run . <url> 10 200
```

| Argument        | Description                        | Default  |
| --------------- | ---------------------------------- | -------- |
| URL             | Starting URL                       | Required |
| Max Concurrency | Number of concurrent crawl workers | 5        |
| Max Pages       | Maximum number of pages to crawl   | 100      |

## Output

After crawling completes, a `report.json` file is generated containing the extracted page information.

Each record contains:

```json
{
  "url": "...",
  "heading": "...",
  "first_paragraph": "...",
  "outgoing_links": [],
  "image_urls": []
}
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

## How It Works

1. Start from the provided URL.
2. Crawl pages recursively while remaining within the same domain.
3. Normalize URLs to prevent duplicate visits.
4. Fetch and parse HTML pages.
5. Extract structured page information.
6. Store page data in memory.
7. Export a sorted JSON report.

## Roadmap

This project is actively being improved. Planned enhancements include:

- More robust HTTP error handling and retry logic
- Better handling of redirects, rate limiting, and timeouts
- Tracking internal vs. external links
- Crawl statistics and summary reporting
- Smarter duplicate detection
- Configurable User-Agent and request headers
- robots.txt support
- Sitemap generation
- Structured logging
- Benchmarking and performance improvements
- Unit and integration test expansion

## License

GPLV3
