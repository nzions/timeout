# Timeout

A GNU coreutils `timeout` compatible utility for running commands with a timeout.

[![GitHub](https://img.shields.io/badge/GitHub-nzions%2Ftimeout-blue)](https://github.com/nzions/timeout)
[![Go Report Card](https://goreportcard.com/badge/github.com/nzions/timeout)](https://goreportcard.com/report/github.com/nzions/timeout)
[![codecov](https://codecov.io/gh/nzions/timeout/branch/master/graph/badge.svg)](https://codecov.io/gh/nzions/timeout)
[![License: CC0](https://img.shields.io/badge/License-CC0%201.0-lightgrey.svg)](http://creativecommons.org/publicdomain/zero/1.0/)

## Features

- **100% GNU Compatibility**: Drop-in replacement for GNU coreutils `timeout`
- **macOS Native**: Fully tested and optimized for macOS
- **Zero Dependencies**: Single binary with no external dependencies
- **Signal Handling**: Supports custom signals and signal escalation
- **Flexible Duration**: Supports seconds, minutes, hours, and days
- **Comprehensive Testing**: Extensive unit and integration test coverage on macOS

> **Note**: This utility is primarily developed and tested on macOS. It may work on Linux and Windows due to Go's cross-platform nature, but these platforms are currently untested and unsupported.

## Usage

```bash
timeout [OPTION] DURATION COMMAND [ARG]...
```

## Options

- `--kill-after=DURATION` - Also send a KILL signal if command is still running this long after the initial signal was sent
- `--signal=SIGNAL` - Specify the signal to be sent on timeout (default: TERM)
- `--preserve-status` - Exit with the same status as COMMAND, even when the command times out
- `--foreground` - When not running timeout directly from a shell prompt, allow COMMAND to read from the TTY and get TTY signals
- `--verbose` - Diagnose to stderr any signal sent upon timeout
- `--help` - Display help and exit
- `--version` - Output version information and exit

## Duration Format

DURATION is a floating point number with an optional suffix:
- `s` for seconds (the default)
- `m` for minutes  
- `h` for hours
- `d` for days

A duration of 0 disables the associated timeout.

## Examples

```bash
# Run npm test with 30 second timeout
timeout 30s npm test

# Run go test with 5 minute timeout
timeout 5m go test ./...

# Use custom signal and kill-after
timeout --signal=INT --kill-after=10s 30s ./my-script.sh

# Preserve command exit status on timeout
timeout --preserve-status 60s long-running-command

# Verbose output
timeout --verbose 30s some-command
```

## Exit Codes

- 0: Command completed successfully
- 1: Command failed or error starting command
- 124: Command timed out (standard timeout exit code)
- 125: Invalid arguments to timeout command
- 128+N: Command killed by signal N (when using KILL signal)
- 130: Command interrupted by signal (SIGINT/SIGTERM)
- Other: Exit code from the wrapped command

## Signal Names

Supports both numeric signals and named signals (with or without SIG prefix):
- TERM, KILL, INT, QUIT, HUP, USR1, USR2, PIPE, ALRM, etc.
- Numeric signals: 9, 15, 2, etc.

## GNU Compatibility

This implementation is 100% compatible with GNU coreutils `timeout`, including:
- All command-line options
- Duration parsing with suffixes
- Signal handling and escalation
- Exit codes
- Error messages and behavior

## Testing

The timeout utility includes comprehensive unit and integration tests:

```bash
# Run all tests
go test -v

# Run tests with coverage
go test -cover

# Run benchmarks
go test -bench=.

# Run specific test
go test -run TestParseDuration
```

### Test Coverage

- **Unit Tests**: Test duration parsing, signal parsing, edge cases
- **Integration Tests**: Test actual timeout behavior, command execution, signal handling
- **Benchmarks**: Performance testing for parsing functions

Test files:
- `timeout_test.go` - Unit tests for parsing functions
- `integration_test.go` - End-to-end integration tests

## Installation

### Quick Install

The easiest way to install is directly from GitHub:

```bash
go install github.com/nzions/timeout@latest
```

This will download, compile, and install the latest version to your `$GOPATH/bin` directory.

### Manual Installation

Alternatively, you can download pre-built binaries from the [releases page](https://github.com/nzions/timeout/releases).

> **Platform Support**: Pre-built binaries are provided for macOS (tested), with experimental builds for Linux and Windows (untested).

## Building from Source

You can also build the timeout utility from source using standard Go commands:

```bash
# Clone the repository
git clone https://github.com/nzions/timeout.git
cd timeout

# Build the binary (macOS)
go build -o timeout

# Build optimized release version (macOS)
go build -ldflags="-s -w" -o timeout

# Experimental builds for other platforms (untested)
GOOS=linux GOARCH=amd64 go build -o timeout-linux-amd64
GOOS=windows GOARCH=amd64 go build -o timeout-windows-amd64.exe
```

### Development Commands

- `go build` - Build the binary
- `go test -v` - Run all tests
- `go test -cover` - Run tests with coverage
- `go test -bench=.` - Run benchmarks
- `go fmt` - Format code
- `go vet` - Check for common errors

## License

This project is dedicated to the public domain under the [CC0 1.0 Universal (CC0 1.0) Public Domain Dedication](https://creativecommons.org/publicdomain/zero/1.0/).

You can copy, modify, distribute and perform the work, even for commercial purposes, all without asking permission. See the [LICENSE](LICENSE) file for the full legal text.

## Contributing

Contributions are welcome! Since this project is in the public domain under CC0, you can contribute freely without any copyright concerns.

### How to Contribute

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for any new functionality
5. Run tests (`go test ./...`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

### Development Setup

```bash
# Clone your fork
git clone https://github.com/your-username/timeout.git
cd timeout

# Install dependencies (if any)
go mod tidy

# Run tests
go test -v ./...

# Run with coverage
go test -cover ./...
```

Please ensure all tests pass and consider adding tests for new features.
