package executor

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createTestLogger creates a logger for tests (writes to stderr).
func createTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelError, // Only show errors in tests
	}))
}

func TestNew(t *testing.T) {
	t.Run("successfully creates executor when gh is in PATH", func(t *testing.T) {
		logger := createTestLogger()
		exec, err := New(logger)

		require.NoError(t, err, "New() should not return error when gh is in PATH")
		assert.NotNil(t, exec, "executor should not be nil")
		assert.NotEmpty(t, exec.GetGhPath(), "gh path should not be empty")
		assert.Equal(t, 5*time.Minute, exec.timeout, "default timeout should be 5 minutes")
	})

	t.Run("returns error when gh is not in PATH", func(t *testing.T) {
		// Save original PATH
		originalPath := os.Getenv("PATH")
		defer os.Setenv("PATH", originalPath)

		// Set empty PATH to simulate gh not being found
		os.Setenv("PATH", "")

		logger := createTestLogger()
		exec, err := New(logger)

		assert.Error(t, err, "New() should return error when gh is not in PATH")
		assert.Nil(t, exec, "executor should be nil on error")
		assert.Contains(t, err.Error(), "gh CLI not found", "error should mention gh not found")
	})
}

func TestExecutor_GetGhPath(t *testing.T) {
	logger := createTestLogger()
	exec, err := New(logger)
	require.NoError(t, err)

	path := exec.GetGhPath()
	assert.NotEmpty(t, path, "gh path should not be empty")
	assert.Contains(t, path, "gh", "path should contain 'gh'")
}

func TestExecutor_SetTimeout(t *testing.T) {
	logger := createTestLogger()
	exec, err := New(logger)
	require.NoError(t, err)

	// Test default timeout
	assert.Equal(t, 5*time.Minute, exec.timeout, "default timeout should be 5 minutes")

	// Test setting new timeout
	newTimeout := 30 * time.Second
	exec.SetTimeout(newTimeout)
	assert.Equal(t, newTimeout, exec.timeout, "timeout should be updated")
}

func TestExecutor_Execute(t *testing.T) {
	logger := createTestLogger()
	exec, err := New(logger)
	require.NoError(t, err)

	t.Run("successfully executes simple command", func(t *testing.T) {
		ctx := context.Background()
		result, err := exec.Execute(ctx, "--version")

		require.NoError(t, err, "Execute should not return error for valid command")
		assert.NotNil(t, result, "result should not be nil")
		assert.NotEmpty(t, result.Stdout, "stdout should not be empty")
		assert.Equal(t, 0, result.ExitCode, "exit code should be 0")
		assert.Contains(t, strings.ToLower(result.Stdout), "gh version", "output should contain version info")
	})

	t.Run("captures stderr for invalid command", func(t *testing.T) {
		ctx := context.Background()
		result, err := exec.Execute(ctx, "invalid-command-that-does-not-exist")

		assert.Error(t, err, "Execute should return error for invalid command")
		assert.NotNil(t, result, "result should not be nil even on error")
		assert.NotEmpty(t, result.Stderr, "stderr should not be empty")
		assert.NotEqual(t, 0, result.ExitCode, "exit code should not be 0")
	})

	t.Run("respects context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		result, err := exec.Execute(ctx, "--version")

		assert.Error(t, err, "Execute should return error when context is canceled")
		assert.NotNil(t, result, "result should not be nil")
	})

	t.Run("respects timeout", func(t *testing.T) {
		// Set very short timeout
		exec.SetTimeout(1 * time.Nanosecond)
		defer exec.SetTimeout(5 * time.Minute) // Restore default

		ctx := context.Background()
		result, err := exec.Execute(ctx, "--version")

		assert.Error(t, err, "Execute should return error when timeout is exceeded")
		assert.NotNil(t, result, "result should not be nil")
	})

	t.Run("handles commands with multiple arguments", func(t *testing.T) {
		ctx := context.Background()
		result, _ := exec.Execute(ctx, "auth", "status")

		// This command might fail if not authenticated, but we're testing arg handling
		assert.NotNil(t, result, "result should not be nil")
		// We don't assert error here as it depends on auth state
	})

	t.Run("returns result even on command failure", func(t *testing.T) {
		ctx := context.Background()
		result, err := exec.Execute(ctx, "repo", "view", "nonexistent/repository-that-does-not-exist-12345")

		assert.Error(t, err, "Execute should return error for failed command")
		assert.NotNil(t, result, "result should not be nil even on error")
		assert.NotEqual(t, 0, result.ExitCode, "exit code should not be 0")
		// Either stdout or stderr should have content
		assert.True(t, len(result.Stdout) > 0 || len(result.Stderr) > 0, "should have output")
	})
}

func TestExecutor_Execute_Integration(t *testing.T) {
	// Skip integration tests in short mode
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	logger := createTestLogger()
	exec, err := New(logger)
	require.NoError(t, err)

	t.Run("help command works", func(t *testing.T) {
		ctx := context.Background()
		result, err := exec.Execute(ctx, "help")

		require.NoError(t, err)
		assert.Contains(t, result.Stdout, "Work seamlessly with GitHub", "help output should contain GitHub")
		assert.Equal(t, 0, result.ExitCode)
	})

	t.Run("list commands work", func(t *testing.T) {
		ctx := context.Background()

		// Test repo list (requires auth but we can test the command structure)
		result, _ := exec.Execute(ctx, "repo", "list", "--limit", "1")

		// Command might fail if not authenticated, but result should be returned
		assert.NotNil(t, result)
		// If error, stderr should contain meaningful output
		if result.ExitCode != 0 {
			assert.NotEmpty(t, result.Stderr, "stderr should not be empty on error")
		}
	})
}

func TestResult(t *testing.T) {
	t.Run("Result struct holds command output", func(t *testing.T) {
		result := &Result{
			Stdout:   "test stdout",
			Stderr:   "test stderr",
			ExitCode: 1,
		}

		assert.Equal(t, "test stdout", result.Stdout)
		assert.Equal(t, "test stderr", result.Stderr)
		assert.Equal(t, 1, result.ExitCode)
	})
}

func TestSanitizeArgs(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "empty args",
			args: []string{},
			want: "",
		},
		{
			name: "non-sensitive command passes through",
			args: []string{"issue", "list", "--repo", "owner/repo"},
			want: "issue list --repo owner/repo",
		},
		{
			name: "secret set redacts --body value",
			args: []string{"secret", "set", "MY_SECRET", "--body", "super-secret-value"},
			want: "secret set MY_SECRET --body [REDACTED]",
		},
		{
			name: "variable set redacts --body value",
			args: []string{"variable", "set", "MY_VAR", "--body", "some-value"},
			want: "variable set MY_VAR --body [REDACTED]",
		},
		{
			name: "secret set with multiple flags redacts only --body",
			args: []string{"secret", "set", "DB_PASS", "--app", "Actions", "--body", "p@ssw0rd", "--repo", "owner/repo"},
			want: "secret set DB_PASS --app Actions --body [REDACTED] --repo owner/repo",
		},
		{
			name: "secret list is not redacted",
			args: []string{"secret", "list", "--repo", "owner/repo"},
			want: "secret list --repo owner/repo",
		},
		{
			name: "non-sensitive command with --body is not redacted",
			args: []string{"issue", "create", "--title", "Bug", "--body", "Description here"},
			want: "issue create --title Bug --body Description here",
		},
		{
			name: "secret set without --body passes through",
			args: []string{"secret", "set", "MY_SECRET", "--body-file", "secret.txt"},
			want: "secret set MY_SECRET --body-file secret.txt",
		},
		{
			name: "--body at end of args with no value",
			args: []string{"secret", "set", "MY_SECRET", "--body"},
			want: "secret set MY_SECRET --body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitizeArgs(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

// Benchmark tests.
func BenchmarkExecutor_Execute(b *testing.B) {
	logger := createTestLogger()
	exec, err := New(logger)
	require.NoError(b, err)

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = exec.Execute(ctx, "--version")
	}
}
