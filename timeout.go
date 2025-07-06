// Package main implements a GNU coreutils timeout compatible utility.
//
// This utility runs commands with a timeout, sending signals when the timeout
// is exceeded and optionally escalating to KILL if the command doesn't respond.
// It's designed to be 100% compatible with GNU coreutils timeout.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// Version information
const (
	Version = "1.0.0"
	Author  = "github.com/nzions/timeout"
)

// Config holds all the configuration for the timeout command
type Config struct {
	KillAfter      string
	SignalName     string
	PreserveStatus bool
	Foreground     bool
	Verbose        bool
	Help           bool
	Version        bool

	// For testing
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader
}

// Result holds the result of running a command
type Result struct {
	ExitCode int
	Error    error
}

func usage(w io.Writer, progName string) {
	fmt.Fprintf(w, "Usage: %s [OPTION] DURATION COMMAND [ARG]...\n", progName)
	fmt.Fprintf(w, "  or:  %s [OPTION]\n", progName)
	fmt.Fprintf(w, "Start COMMAND, and kill it if still running after DURATION.\n\n")
	fmt.Fprintf(w, "Options:\n")
	flag.CommandLine.SetOutput(w)
	flag.PrintDefaults()
	fmt.Fprintf(w, "\nDURATION is a floating point number with an optional suffix:\n")
	fmt.Fprintf(w, "'s' for seconds (the default), 'm' for minutes, 'h' for hours or 'd' for days.\n")
	fmt.Fprintf(w, "A duration of 0 disables the associated timeout.\n\n")
	fmt.Fprintf(w, "If the command times out, and --preserve-status is not set, then exit with\n")
	fmt.Fprintf(w, "status 124.  Otherwise, exit with the status of COMMAND.  If no signal\n")
	fmt.Fprintf(w, "is specified, send the TERM signal upon timeout.  The TERM signal kills\n")
	fmt.Fprintf(w, "any process that does not block or catch that signal.  It may be necessary\n")
	fmt.Fprintf(w, "to use the KILL (9) signal, since this signal cannot be caught, in which\n")
	fmt.Fprintf(w, "case the exit status is 128+9 rather than 124.\n")
}

func parseDuration(s string) (time.Duration, error) {
	if s == "" {
		return 0, fmt.Errorf("empty duration")
	}

	// Handle suffixes
	var multiplier time.Duration = time.Second
	suffix := s[len(s)-1:]

	switch suffix {
	case "s":
		s = s[:len(s)-1]
		multiplier = time.Second
	case "m":
		s = s[:len(s)-1]
		multiplier = time.Minute
	case "h":
		s = s[:len(s)-1]
		multiplier = time.Hour
	case "d":
		s = s[:len(s)-1]
		multiplier = 24 * time.Hour
	default:
		// No suffix, assume seconds
		multiplier = time.Second
	}

	// Parse the numeric part
	if f, err := strconv.ParseFloat(s, 64); err != nil {
		return 0, err
	} else {
		return time.Duration(f * float64(multiplier)), nil
	}
}

func parseSignal(s string) (syscall.Signal, error) {
	// Handle numeric signals
	if num, err := strconv.Atoi(s); err == nil {
		return syscall.Signal(num), nil
	}

	// Handle named signals (with or without SIG prefix)
	s = strings.ToUpper(s)
	if !strings.HasPrefix(s, "SIG") {
		s = "SIG" + s
	}

	signals := map[string]syscall.Signal{
		"SIGTERM": syscall.SIGTERM,
		"SIGKILL": syscall.SIGKILL,
		"SIGINT":  syscall.SIGINT,
		"SIGQUIT": syscall.SIGQUIT,
		"SIGHUP":  syscall.SIGHUP,
		"SIGUSR1": syscall.SIGUSR1,
		"SIGUSR2": syscall.SIGUSR2,
		"SIGPIPE": syscall.SIGPIPE,
		"SIGALRM": syscall.SIGALRM,
		"SIGCHLD": syscall.SIGCHLD,
		"SIGCONT": syscall.SIGCONT,
		"SIGSTOP": syscall.SIGSTOP,
		"SIGTSTP": syscall.SIGTSTP,
		"SIGTTIN": syscall.SIGTTIN,
		"SIGTTOU": syscall.SIGTTOU,
	}

	if sig, ok := signals[s]; ok {
		return sig, nil
	}

	return 0, fmt.Errorf("invalid signal: %s", s)
}

// runTimeout executes the timeout logic and returns the result
func runTimeout(config Config, args []string) Result {
	if config.Help {
		usage(config.Stderr, "timeout")
		return Result{ExitCode: 0}
	}

	if config.Version {
		fmt.Fprintf(config.Stdout, "timeout (GNU coreutils compatible) %s\n", Version)
		fmt.Fprintf(config.Stdout, "Source: %s\n", Author)
		fmt.Fprintf(config.Stdout, "License: CC0 1.0 Universal (Public Domain)\n")
		return Result{ExitCode: 0}
	}

	if len(args) < 2 {
		fmt.Fprintf(config.Stderr, "timeout: missing operand\n")
		fmt.Fprintf(config.Stderr, "Try 'timeout --help' for more information.\n")
		return Result{ExitCode: 125}
	}

	// Parse timeout
	timeoutDuration, err := parseDuration(args[0])
	if err != nil {
		fmt.Fprintf(config.Stderr, "timeout: invalid time interval '%s'\n", args[0])
		return Result{ExitCode: 125}
	}

	// Get command and args
	command := args[1]
	cmdArgs := args[2:]

	// Parse signal
	timeoutSignal, err := parseSignal(config.SignalName)
	if err != nil {
		fmt.Fprintf(config.Stderr, "timeout: %v\n", err)
		return Result{ExitCode: 125}
	}

	// Parse kill-after duration
	var killAfterDuration time.Duration
	if config.KillAfter != "" {
		killAfterDuration, err = parseDuration(config.KillAfter)
		if err != nil {
			fmt.Fprintf(config.Stderr, "timeout: invalid time interval '%s'\n", config.KillAfter)
			return Result{ExitCode: 125}
		}
	}

	// Create context with timeout (0 duration means no timeout)
	var ctx context.Context
	var cancel context.CancelFunc
	if timeoutDuration > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), timeoutDuration)
		defer cancel()
	} else {
		ctx = context.Background()
	}

	// Create command
	cmd := exec.CommandContext(ctx, command, cmdArgs...)
	cmd.Stdout = config.Stdout
	cmd.Stderr = config.Stderr
	cmd.Stdin = config.Stdin

	// Handle interrupt signals to clean up properly
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigChan)

	// Start the command
	if err := cmd.Start(); err != nil {
		fmt.Fprintf(config.Stderr, "Error starting command: %v\n", err)
		return Result{ExitCode: 1, Error: err}
	}

	// Wait for either completion or signal
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-ctx.Done():
		// Timeout occurred
		if config.Verbose {
			fmt.Fprintf(config.Stderr, "timeout: sending signal %s to command '%s'\n", config.SignalName, command)
		}

		if cmd.Process != nil {
			// Send the specified signal
			if err := cmd.Process.Signal(timeoutSignal); err != nil && config.Verbose {
				fmt.Fprintf(config.Stderr, "timeout: failed to send signal: %v\n", err)
			}

			// If kill-after is specified, wait and then send KILL
			if config.KillAfter != "" && killAfterDuration > 0 {
				select {
				case <-time.After(killAfterDuration):
					if config.Verbose {
						fmt.Fprintf(config.Stderr, "timeout: sending signal KILL to command '%s'\n", command)
					}
					cmd.Process.Signal(syscall.SIGKILL)
				case <-done:
					// Process exited before kill-after timeout
				}
			}
		}

		// Wait for process to finish
		<-done

		if config.PreserveStatus {
			// Exit with command's status (if available)
			if cmd.ProcessState != nil {
				return Result{ExitCode: cmd.ProcessState.ExitCode()}
			}
			return Result{ExitCode: 1}
		} else {
			// Standard timeout exit code
			if timeoutSignal == syscall.SIGKILL {
				return Result{ExitCode: 128 + 9} // 128 + SIGKILL
			}
			return Result{ExitCode: 124}
		}
	case sig := <-sigChan:
		// Signal received
		if cmd.Process != nil {
			cmd.Process.Signal(sig)
		}
		<-done                       // Wait for process to finish
		return Result{ExitCode: 130} // Standard interrupt exit code
	case err := <-done:
		// Command completed
		if err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				return Result{ExitCode: exitError.ExitCode()}
			}
			fmt.Fprintf(config.Stderr, "timeout: %v\n", err)
			return Result{ExitCode: 1, Error: err}
		}
		return Result{ExitCode: 0}
	}
}

var (
	killAfter      = flag.String("kill-after", "", "also send a KILL signal if command is still running this long after the initial signal was sent")
	signalName     = flag.String("signal", "TERM", "specify the signal to be sent on timeout")
	preserveStatus = flag.Bool("preserve-status", false, "exit with the same status as COMMAND, even when the command times out")
	foreground     = flag.Bool("foreground", false, "when not running timeout directly from a shell prompt, allow COMMAND to read from the TTY and get TTY signals")
	verbose        = flag.Bool("verbose", false, "diagnose to stderr any signal sent upon timeout")
	help           = flag.Bool("help", false, "display this help and exit")
	version        = flag.Bool("version", false, "output version information and exit")
)

func main() {
	flag.Usage = func() { usage(os.Stderr, os.Args[0]) }
	flag.Parse()

	config := Config{
		KillAfter:      *killAfter,
		SignalName:     *signalName,
		PreserveStatus: *preserveStatus,
		Foreground:     *foreground,
		Verbose:        *verbose,
		Help:           *help,
		Version:        *version,
		Stdout:         os.Stdout,
		Stderr:         os.Stderr,
		Stdin:          os.Stdin,
	}

	result := runTimeout(config, flag.Args())
	os.Exit(result.ExitCode)
}
