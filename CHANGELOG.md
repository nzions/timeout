# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.0] - 2025-07-05

### Added
- Initial release of timeout utility
- GNU coreutils `timeout` compatible implementation
- Full macOS support (Intel and Apple Silicon)
- Support for all GNU timeout command-line options:
  - `--kill-after=DURATION` - Send KILL signal after timeout
  - `--signal=SIGNAL` - Custom signal to send on timeout
  - `--preserve-status` - Preserve command exit status
  - `--foreground` - Allow command to access TTY
  - `--verbose` - Verbose timeout diagnostics
  - `--help` - Display help information
  - `--version` - Show version information
- Duration format support (s, m, h, d suffixes)
- Signal name and number support
- Experimental builds for Linux and Windows (untested)
- Comprehensive test suite with unit and integration tests
- CC0 1.0 Universal Public Domain Dedication license

### Features
- 100% GNU coreutils `timeout` compatibility
- Zero external dependencies
- Signal escalation (TERM → KILL)
- Proper exit code handling (124 for timeout, command codes otherwise)
- Thread-safe implementation
- Extensive error handling and validation
- Native macOS support with comprehensive testing
- Experimental Linux and Windows builds (untested)

[Unreleased]: https://github.com/nzions/timeout/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/nzions/timeout/releases/tag/v1.0.0
