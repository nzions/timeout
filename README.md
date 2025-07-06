# Timeout

A GNU coreutils `timeout` compatible utility for running commands with a timeout.

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

## Building

Use the included Makefile for easy building and development:

```bash
# Build the binary
make build

# Build optimized release version
make build-release

# Build for multiple platforms
make build-all

# Show all available targets
make help
```

### Available Make Targets

- `make build` - Build the binary to `build/timeout`
- `make build-release` - Build optimized release binary
- `make build-dev` - Build with debug info for development
- `make build-all` - Cross-compile for Linux, macOS, Windows
- `make test` - Run all tests
- `make test-coverage` - Run tests with HTML coverage report
- `make bench` - Run benchmarks
- `make test-all` - Run tests, coverage, and benchmarks
- `make clean` - Clean build artifacts
- `make install` - Install to `/usr/local/bin` (requires sudo)
- `make uninstall` - Remove from system (requires sudo)
- `make fmt` - Format code
- `make lint` - Lint code (requires golint)
- `make check` - Run format, lint, and tests
- `make release` - Create release archive with binaries
