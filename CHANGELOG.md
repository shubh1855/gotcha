# Changelog

## [0.8.1] - 2026-10-01

### Fixed

- Set the canonical Go module path to `github.com/shubh1855/gotcha`
- Enable installation with `go install github.com/shubh1855/gotcha@version`

## [0.8.0] - 2026-09-30

### Added

- XML sitemap generation with `--sitemap`
- CSV report generation with `--csv`
- Markdown report generation with `--markdown`
- Benchmark suite for URL normalization and page extraction
- Integration tests for multi-page crawling
- Integration tests for same-domain restrictions
- Integration tests for crawl depth limits
- Integration tests for maximum page limits

### Changed

- Improved URL normalization and duplicate URL detection
- Enforced maximum page limits atomically during concurrent crawling
- Expanded crawl report formats
- Expanded CLI documentation

### Testing

- Added crawler integration test coverage
- Added race-condition testing
- Added benchmark coverage

## [0.7.0] - 2026-08-05

### Added

- Configurable HTTP redirect handling
- Configurable request rate limiting
- `--max-redirects` CLI flag
- `--rate` CLI flag
- Integration tests for HTTP redirects

## [0.6.0] - 2026-08-03

### Added

- Configurable crawl depth using the `--depth` flag
- Configurable request delay using the `--delay` flag
- Unit tests for crawl depth behavior
- Dedicated `golangci-lint` GitHub Actions workflow

### Changed

- Improved HTTP request logging with typed `slog` attributes
- Updated CLI to support additional crawler configuration flags
- Updated project documentation and usage examples

### Fixed

- Improved logging consistency throughout the HTTP fetch pipeline

---

## [0.5.0] - 2026-08-01

### Added

- HTTP retry logic with exponential backoff
- Typed HTTP and content-type errors
- Crawl statistics
- Internal and external link classification
- Configurable User-Agent
- `robots.txt` support
- Structured logging with `log/slog`
- Modern CLI using `pflag`
- Version command
- Continuous Integration with GitHub Actions
- Automated release workflow using GoReleaser

### Changed

- Improved HTTP client implementation
- Improved crawler logging
- Improved CLI argument parsing

### Fixed

- `robots.txt` matching logic
- URL normalization edge cases
