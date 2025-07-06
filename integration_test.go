package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestTimeoutIntegration(t *testing.T) {
	// Build the timeout binary first
	buildCmd := exec.Command("go", "build", "-o", "timeout_test", "timeout.go")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build timeout binary: %v", err)
	}
	defer os.Remove("timeout_test")

	tests := []struct {
		name          string
		args          []string
		expectedExit  int
		shouldTimeout bool
		maxDuration   time.Duration
	}{
		{
			name:         "successful command",
			args:         []string{"5s", "echo", "hello"},
			expectedExit: 0,
			maxDuration:  2 * time.Second,
		},
		{
			name:          "timeout occurs",
			args:          []string{"1s", "sleep", "3"},
			expectedExit:  124,
			shouldTimeout: true,
			maxDuration:   2 * time.Second,
		},
		{
			name:         "zero timeout disables timeout",
			args:         []string{"0", "echo", "test"},
			expectedExit: 0,
			maxDuration:  2 * time.Second,
		},
		{
			name:         "command with exit code",
			args:         []string{"5s", "sh", "-c", "exit 42"},
			expectedExit: 42,
			maxDuration:  2 * time.Second,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			start := time.Now()

			cmd := exec.Command("./timeout_test", test.args...)
			err := cmd.Run()

			duration := time.Since(start)

			// Check duration
			if duration > test.maxDuration {
				t.Errorf("Command took too long: %v (max: %v)", duration, test.maxDuration)
			}

			// Check exit code
			exitCode := 0
			if err != nil {
				if exitError, ok := err.(*exec.ExitError); ok {
					exitCode = exitError.ExitCode()
				} else {
					t.Fatalf("Unexpected error type: %v", err)
				}
			}

			if exitCode != test.expectedExit {
				t.Errorf("Expected exit code %d, got %d", test.expectedExit, exitCode)
			}
		})
	}
}

func TestTimeoutHelp(t *testing.T) {
	buildCmd := exec.Command("go", "build", "-o", "timeout_test", "timeout.go")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build timeout binary: %v", err)
	}
	defer os.Remove("timeout_test")

	cmd := exec.Command("./timeout_test", "--help")
	output, err := cmd.CombinedOutput()

	// Help should exit with code 0
	if err != nil {
		t.Errorf("Help command failed: %v", err)
	}

	outputStr := string(output)
	expectedStrings := []string{
		"Usage:",
		"DURATION",
		"COMMAND",
		"Options:",
		"kill-after",
		"signal",
		"preserve-status",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(outputStr, expected) {
			t.Errorf("Help output missing expected string: %q", expected)
		}
	}
}

func TestTimeoutVersion(t *testing.T) {
	buildCmd := exec.Command("go", "build", "-o", "timeout_test", "timeout.go")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build timeout binary: %v", err)
	}
	defer os.Remove("timeout_test")

	cmd := exec.Command("./timeout_test", "--version")
	output, err := cmd.CombinedOutput()

	// Version should exit with code 0
	if err != nil {
		t.Errorf("Version command failed: %v", err)
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "timeout") || !strings.Contains(outputStr, "1.0") {
		t.Errorf("Version output unexpected: %q", outputStr)
	}
}

func TestTimeoutInvalidArgs(t *testing.T) {
	buildCmd := exec.Command("go", "build", "-o", "timeout_test", "timeout.go")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build timeout binary: %v", err)
	}
	defer os.Remove("timeout_test")

	tests := []struct {
		name string
		args []string
	}{
		{"no args", []string{}},
		{"only duration", []string{"30s"}},
		{"invalid duration", []string{"invalid", "echo", "test"}},
		{"invalid signal", []string{"--signal=INVALID", "30s", "echo", "test"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cmd := exec.Command("./timeout_test", test.args...)
			err := cmd.Run()

			// Should exit with error code (125 for invalid args)
			if err == nil {
				t.Errorf("Expected command to fail, but it succeeded")
				return
			}

			if exitError, ok := err.(*exec.ExitError); ok {
				exitCode := exitError.ExitCode()
				if exitCode != 125 && exitCode != 1 {
					t.Errorf("Expected exit code 125 or 1, got %d", exitCode)
				}
			}
		})
	}
}

func TestTimeoutWithSignal(t *testing.T) {
	buildCmd := exec.Command("go", "build", "-o", "timeout_test", "timeout.go")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build timeout binary: %v", err)
	}
	defer os.Remove("timeout_test")

	// Test custom signal
	cmd := exec.Command("./timeout_test", "--signal=KILL", "1s", "sleep", "3")
	start := time.Now()
	err := cmd.Run()
	duration := time.Since(start)

	// Should timeout within reasonable time
	if duration > 2*time.Second {
		t.Errorf("Command took too long with KILL signal: %v", duration)
	}

	// Should exit with timeout code or KILL signal code
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode := exitError.ExitCode()
			// Accept either 124 (timeout) or 137 (128+9, killed by SIGKILL)
			if exitCode != 124 && exitCode != 137 {
				t.Errorf("Expected exit code 124 or 137, got %d", exitCode)
			}
		}
	} else {
		t.Errorf("Expected command to fail due to timeout")
	}
}

func TestTimeoutPreserveStatus(t *testing.T) {
	buildCmd := exec.Command("go", "build", "-o", "timeout_test", "timeout.go")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build timeout binary: %v", err)
	}
	defer os.Remove("timeout_test")

	// Test preserve-status with timeout
	cmd := exec.Command("./timeout_test", "--preserve-status", "1s", "sh", "-c", "sleep 3; exit 42")
	err := cmd.Run()

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode := exitError.ExitCode()
			// With preserve-status, should get the command's exit code even on timeout
			// But since sleep was killed, we might get different codes
			if exitCode == 124 {
				t.Errorf("Got standard timeout exit code 124, but preserve-status was set")
			}
		}
	}
}
