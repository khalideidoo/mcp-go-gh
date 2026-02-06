# mcp-go-gh

A lightweight, comprehensive Go-based MCP (Model Context Protocol) server that wraps the GitHub CLI (`gh`), exposing every parameter and option available in the tool.

## Features

- **Complete Coverage**: 152 MCP tools covering 27 `gh` command groups - **100% of stable gh CLI commands**
- **Type-Safe**: Go structs with automatic JSON schema generation
- **Code Generated**: YAML definitions drive automatic Go code generation
- **Maintainable**: Easy to update when `gh` CLI evolves
- **Lightweight**: Single binary with no external dependencies (except `gh` CLI)

## Supported Commands

The server exposes **100% of stable gh CLI commands** across 27 command groups:

### Core Commands
- **project** (19): GitHub Projects v2 - full CRUD, fields, and items
- **pr** (14): Pull request management
- **issue** (14): Issue tracking and management
- **codespace** (13): Codespace creation and management
- **repo** (11): Repository operations
- **extension** (8): Extension installation and management
- **release** (7): Release management
- **run** (7): Workflow run management
- **auth** (6): Authentication and setup
- **gist** (6): Gist management

### Actions & Workflow Commands
- **workflow** (5): GitHub Actions workflow management
- **variable** (4): Actions variables
- **secret** (3): Secrets management
- **cache** (2): Actions cache operations

### Additional Commands
- **label** (5): Label management
- **alias** (4): Command shortcuts
- **config** (4): Configuration management
- **attestation** (3): Artifact attestations
- **gpg-key** (3): GPG key management
- **ruleset** (3): Repository rulesets
- **search** (3): Search repos, issues, and PRs
- **ssh-key** (3): SSH key management
- **cache** (2): Actions cache management
- **api** (1): Raw GitHub API access
- **browse** (1): Open resources in browser
- **completion** (1): Shell completion
- **org** (1): Organization operations
- **status** (1): Status overview

**Total: 152 MCP tools = 100% stable command coverage** ✅

## Prerequisites

- **Go 1.25.6 or later**
- **GitHub CLI (`gh`)** installed and authenticated
  - Install: `brew install gh` (macOS) or see [official docs](https://cli.github.com/)
  - Authenticate: `gh auth login`
- **(Optional)** golangci-lint v2 for development - [installation guide](https://golangci-lint.run/welcome/install/)

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/khalideidoo/mcp-go-gh.git
cd mcp-go-gh

# Build the server
make build

# Or install to GOPATH/bin
make install
```

### Using Go Install

```bash
go install github.com/khalideidoo/mcp-go-gh/cmd/mcp-go-gh@latest
```

## Usage

### With Claude Desktop

Add to your Claude Desktop configuration (`~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):

```json
{
  "mcpServers": {
    "gh": {
      "command": "/path/to/mcp-go-gh"
    }
  }
}
```

### With Other MCP Clients

The server uses stdio transport and follows the MCP protocol specification. Configure your client to launch the `mcp-go-gh` binary.

### Environment Variables

The server respects all `gh` CLI environment variables:

- `GH_TOKEN` / `GITHUB_TOKEN`: Authentication token
- `GH_HOST`: GitHub hostname (for Enterprise)
- `GH_REPO`: Default repository
- `GH_EDITOR`: Text editor preference
- And more (see `gh` documentation)

## Example Tools

### Create a Pull Request

```json
{
  "name": "gh_pr_create",
  "arguments": {
    "title": "Fix bug in authentication",
    "body": "This PR fixes the authentication bug",
    "base": "main",
    "draft": false,
    "assignee": ["@me"]
  }
}
```

### List Issues

```json
{
  "name": "gh_issue_list",
  "arguments": {
    "state": "open",
    "label": ["bug", "priority"],
    "limit": 10,
    "json": ["number", "title", "author"]
  }
}
```

### Search Repositories

```json
{
  "name": "gh_search_repos",
  "arguments": {
    "query": "language:go stars:>1000",
    "limit": 20
  }
}
```

### GitHub API Access

```json
{
  "name": "gh_api_request",
  "arguments": {
    "endpoint": "repos/{owner}/{repo}/issues",
    "method": "GET",
    "jq": ".[] | {number, title}"
  }
}
```

## Development

### Project Structure

```
mcp-go-gh/
├── cmd/
│   └── mcp-go-gh/          # Server entry point
├── internal/
│   ├── commands/
│   │   ├── definitions/    # YAML command definitions (27 files)
│   │   └── generated/      # Generated Go code (152 tools)
│   ├── executor/           # gh CLI executor
│   └── server/             # MCP server logic
├── tools/
│   └── gen/                # Code generator
├── .golangci.yml           # golangci-lint v2 configuration
├── Makefile                # Build automation
└── README.md
```

### Adding New Commands

1. Create or update YAML definition in `internal/commands/definitions/`
2. Run `make generate` to generate Go code
3. Build: `make build`

Example YAML definition:

```yaml
command: example
description: Example command
subcommands:
  - name: create
    description: Create something
    parameters:
      - name: title
        type: string
        flag: --title
        short: -t
        description: Title for the item
      - name: draft
        type: boolean
        flag: --draft
        description: Create as draft
```

### Building

```bash
# Generate code and build (default)
make

# Just generate code
make generate

# Just build
make build

# Build for all platforms
make build-all

# Clean build artifacts
make clean

# Install to GOPATH/bin
make install
```

### Code Quality

This project uses **golangci-lint v2** for comprehensive code quality checks:

```bash
# Run all linters
make lint

# Auto-fix issues (formatting, imports, etc.)
make lint-fix

# Format code
make fmt

# Or use golangci-lint v2 formatter directly
golangci-lint fmt
```

**Linters enabled**: 25+ including errcheck, govet, staticcheck, gosec, revive, and more. See [.golangci.yml](.golangci.yml) for full configuration.

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
go test -v -coverprofile=coverage.out ./...

# View coverage report
go tool cover -html=coverage.out
```

### Available Make Targets

| Target | Description |
|--------|-------------|
| `make` or `make all` | Generate code and build (default) |
| `make generate` | Generate Go code from YAML definitions |
| `make build` | Build the MCP server binary |
| `make test` | Run all tests |
| `make lint` | Run golangci-lint v2 |
| `make lint-fix` | Run golangci-lint v2 with auto-fix |
| `make fmt` | Format code with go fmt |
| `make install` | Install binary to GOPATH/bin |
| `make clean` | Remove build artifacts |
| `make build-all` | Build for multiple platforms |
| `make deps` | Install and tidy dependencies |
| `make help` | Show available targets |

## Architecture

### Code Generation

The project uses a YAML-driven code generation approach:

1. **YAML Definitions**: Command structures defined in `internal/commands/definitions/*.yaml`
2. **Code Generator**: `tools/gen/` reads YAML and generates Go code
3. **Generated Code**: Type-safe structs and registration functions in `internal/commands/generated/`

This approach ensures:
- Consistency across all commands
- Easy maintenance and updates
- Comprehensive parameter coverage
- Automatic JSON schema generation

### Executor

The `internal/executor` package handles `gh` CLI execution:
- Finds `gh` binary in PATH
- Executes commands with proper timeout handling
- Captures stdout/stderr
- Logs all operations to stderr (stdout reserved for MCP protocol)

## Minimum Requirements

- **gh CLI**: Version 2.30.0 or later
- **Go**: 1.25.6 or later (for development)
- **OS**: macOS, Linux, or Windows
- **golangci-lint**: v2.8.0 or later (optional, for development)

## Troubleshooting

### gh CLI Not Found

Ensure `gh` is in your PATH:
```bash
which gh
# Should output: /usr/local/bin/gh or similar
```

### Authentication Issues

Check `gh` authentication status:
```bash
gh auth status
```

If not authenticated:
```bash
gh auth login
```

### Debugging

The server logs to stderr. Enable debug logging:
```bash
# Run directly to see logs
./bin/mcp-go-gh

# Or check your MCP client's logs
```

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Add/update YAML definitions for new commands
4. Run `make generate` to generate Go code
5. Run `make lint` to ensure code quality
6. Run `make test` to verify tests pass
7. Run `make build` to verify it compiles
8. Submit a pull request

### Development Workflow

```bash
# 1. Make changes to YAML definitions
vim internal/commands/definitions/example.yaml

# 2. Generate code
make generate

# 3. Run quality checks
make lint-fix  # Auto-fix issues
make lint      # Verify all checks pass

# 4. Run tests
make test

# 5. Build
make build
```

## License

MIT License - see [LICENSE](LICENSE) file for details

## Acknowledgments

- Built with the [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk)
- Wraps the [GitHub CLI](https://cli.github.com/)
- Inspired by the Model Context Protocol specification

## Links

- [Model Context Protocol](https://modelcontextprotocol.io/)
- [GitHub CLI Documentation](https://cli.github.com/manual/)
- [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk)
