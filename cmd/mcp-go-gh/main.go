package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/khalideidoo/mcp-go-gh/internal/commands/generated"
	"github.com/khalideidoo/mcp-go-gh/internal/executor"
)

func main() {
	// Set up structured logging to stderr only (stdout is reserved for MCP protocol)
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	logger.Info("starting mcp-go-gh server")

	// Create executor for running gh CLI commands
	exec, err := executor.New(logger)
	if err != nil {
		logger.Error("failed to create executor", "error", err)
		os.Exit(1)
	}

	logger.Info("initialized gh CLI executor", "gh_path", exec.GetGhPath())

	// Create MCP server
	impl := &mcp.Implementation{
		Name:    "mcp-go-gh",
		Title:   "GitHub CLI MCP Server",
		Version: "1.0.0",
	}

	server := mcp.NewServer(impl, &mcp.ServerOptions{})

	logger.Info("created MCP server", "name", "mcp-go-gh", "version", "1.0.0")

	// Register all generated gh command tools
	generated.RegisterAllTools(server, exec)
	logger.Info("registered all tools successfully")

	// Create stdio transport for communication
	transport := &mcp.StdioTransport{}

	// Start the server
	logger.Info("starting MCP server with stdio transport")
	if err := server.Run(context.Background(), transport); err != nil {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}

	logger.Info("server stopped")
}
