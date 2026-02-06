package executor

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
	"time"
)

// Executor handles execution of gh CLI commands
type Executor struct {
	ghPath  string
	timeout time.Duration
	logger  *slog.Logger
}

// Result contains the output of a command execution
type Result struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

// New creates a new Executor instance
func New(logger *slog.Logger) (*Executor, error) {
	// Find gh CLI in PATH
	ghPath, err := exec.LookPath("gh")
	if err != nil {
		return nil, fmt.Errorf("gh CLI not found in PATH: %w", err)
	}

	return &Executor{
		ghPath:  ghPath,
		timeout: 5 * time.Minute, // Default timeout
		logger:  logger,
	}, nil
}

// Execute runs a gh command with the given arguments
func (e *Executor) Execute(ctx context.Context, args ...string) (*Result, error) {
	// Apply timeout
	ctx, cancel := context.WithTimeout(ctx, e.timeout)
	defer cancel()

	// Build command
	cmd := exec.CommandContext(ctx, e.ghPath, args...)

	// Capture stdout and stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Log command execution
	e.logger.Info("executing gh command",
		"command", "gh",
		"args", strings.Join(args, " "))

	// Execute command
	err := cmd.Run()

	// Get exit code
	exitCode := 0
	if cmd.ProcessState != nil {
		exitCode = cmd.ProcessState.ExitCode()
	}

	result := &Result{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: exitCode,
	}

	if err != nil {
		e.logger.Error("gh command failed",
			"error", err,
			"stderr", result.Stderr,
			"exit_code", exitCode,
			"args", strings.Join(args, " "))

		return result, fmt.Errorf("gh command failed (exit %d): %s", exitCode, result.Stderr)
	}

	e.logger.Debug("gh command succeeded",
		"exit_code", exitCode,
		"args", strings.Join(args, " "))

	return result, nil
}

// SetTimeout changes the default command timeout
func (e *Executor) SetTimeout(timeout time.Duration) {
	e.timeout = timeout
}

// GetGhPath returns the path to the gh binary
func (e *Executor) GetGhPath() string {
	return e.ghPath
}
