# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- **Retro Terminal Interface**: A Fallout-style Amber CRT web interface for managing servers (`/terminal`).
- **Web Package**: New `web` package to serve embedded static assets and handle API requests.
- **Docker Networking**: Integrated `caddy-test-net` into `docker-compose.yml` for reverse proxy support.

### Fixed

- **Bot Conflict**: Resolved "terminated by other getUpdates request" error by cleaning up zombie processes.
- **Configuration**: Fixed malformed `.env` file handling in `main.go` and script execution.

## [1.0.0] - 2026-01-30

### Added

- Initial release of Watchtower Masterbot.
- Multi-server management with secure credential storage.
