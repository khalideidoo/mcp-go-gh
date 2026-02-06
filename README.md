# mcp-go-gh

A lightweight, comprehensive Go-based MCP (Model Context Protocol) server that wraps the GitHub CLI (`gh`), exposing every parameter and option available in the tool.

## Features

- **Comprehensive Coverage**: 118 MCP tools covering 19 major `gh` command groups
- **Type-Safe**: Go structs with automatic JSON schema generation
- **Code Generated**: YAML definitions drive automatic Go code generation
- **Maintainable**: Easy to update when `gh` CLI evolves
- **Lightweight**: Single binary with no external dependencies (except `gh` CLI)

## Supported Commands

The server exposes the following `gh` CLI command groups:

- **project** (19 commands): GitHub Projects v2 management
- **pr** (14 commands): Pull request management
- **issue** (14 commands): Issue management
- **repo** (11 commands): Repository operations
- **extension** (8 commands): Extension management
- **release** (7 commands): Release management
- **run** (7 commands): Workflow run management
- **auth** (6 commands): Authentication
- **gist** (6 commands): Gist management
- **workflow** (5 commands): GitHub Actions workflows
- **label** (5 commands): Label management
- **variable** (4 commands): Actions variables
- **search** (3 commands): Search repos, issues, and PRs
- **secret** (3 commands): Secrets management
- **cache** (2 commands): Actions cache management
- **api** (1 command): Raw GitHub API access
- **browse** (1 command): Open in browser
- **org** (1 command): Organization management
- **status** (1 command): Show relevant items

**Total: 118 MCP tools**

## Prerequisites

- Go 1.21 or later
- GitHub CLI (`gh`) installed and authenticated
  - Install: `brew install gh` (macOS) or see [official docs](https://cli.github.com/)
  - Authenticate: `gh auth login`

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
│   │   ├── definitions/    # YAML command definitions
│   │   └── generated/      # Generated Go code
│   ├── executor/           # gh CLI executor
│   └── server/             # MCP server logic
├── tools/
│   └── gen/                # Code generator
├── Makefile
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
# Generate code and build
make

# Just generate code
make generate

# Just build
make build

# Build for all platforms
make build-all

# Clean build artifacts
make clean
```

### Running Tests

```bash
make test
```

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
- **Go**: 1.21 or later (for development)
- **OS**: macOS, Linux, or Windows

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
4. Run `make generate` and `make build`
5. Test your changes
6. Submit a pull request

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
