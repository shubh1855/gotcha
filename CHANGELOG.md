# Changelog

All notable changes to this project will be documented in this file.

## [0.5.0] - 2026-08-01

### Added

- HTTP retry logic with exponential backoff
- Typed HTTP and content-type errors
- Crawl statistics
- Internal and external link classification
- Configurable User-Agent
- robots.txt support
- Structured logging with slog
- Modern CLI using pflag
- Version command
- Continuous Integration with GitHub Actions
- Automated release workflow using GoReleaser

### Changed

- Improved HTTP client implementation
- Improved crawler logging
- Improved CLI argument parsing

### Fixed

- robots.txt matching logic
- URL normalization edge cases
