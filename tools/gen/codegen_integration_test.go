package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateCode(t *testing.T) {
	t.Run("successfully generates code from valid definitions", func(t *testing.T) {
		// Create test definitions
		definitions := []CommandDefinition{
			{
				Command:     "test",
				Description: "Test command",
				Subcommands: []Subcommand{
					{
						Name:        "list",
						Description: "List items",
						Parameters: []Parameter{
							{
								Name:        "limit",
								Type:        "integer",
								Flag:        "--limit",
								Short:       "-L",
								Description: "Maximum number",
							},
						},
					},
					{
						Name:        "create",
						Description: "Create item",
						Parameters: []Parameter{
							{
								Name:        "name",
								Type:        "string",
								Description: "Item name",
								Positional:  true,
								Required:    true,
							},
						},
					},
				},
			},
		}

		// Create temporary output directory
		tmpDir := t.TempDir()

		err := GenerateCode(definitions, tmpDir)
		require.NoError(t, err)

		// Verify files were created
		commandFile := filepath.Join(tmpDir, "test_gen.go")
		registryFile := filepath.Join(tmpDir, "registry_gen.go")

		assert.FileExists(t, commandFile)
		assert.FileExists(t, registryFile)

		// Read and verify command file content
		commandContent, err := os.ReadFile(commandFile)
		require.NoError(t, err)

		contentStr := string(commandContent)
		assert.Contains(t, contentStr, "package generated")
		assert.Contains(t, contentStr, "TestListArgs")
		assert.Contains(t, contentStr, "TestCreateArgs")
		assert.Contains(t, contentStr, "RegisterTestListTool")
		assert.Contains(t, contentStr, "RegisterTestCreateTool")
		assert.Contains(t, contentStr, `Limit int`)
		assert.Contains(t, contentStr, `Name string`)
		assert.Contains(t, contentStr, `"test", "list"`)
		assert.Contains(t, contentStr, `"test", "create"`)

		// Verify registry file content
		registryContent, err := os.ReadFile(registryFile)
		require.NoError(t, err)

		registryStr := string(registryContent)
		assert.Contains(t, registryStr, "package generated")
		assert.Contains(t, registryStr, "RegisterAllTools")
		assert.Contains(t, registryStr, "RegisterTestListTool")
		assert.Contains(t, registryStr, "RegisterTestCreateTool")
	})

	t.Run("handles multiple command definitions", func(t *testing.T) {
		definitions := []CommandDefinition{
			{
				Command:     "cmd1",
				Description: "First command",
				Subcommands: []Subcommand{
					{Name: "sub1", Description: "Subcommand 1", Parameters: []Parameter{}},
				},
			},
			{
				Command:     "cmd2",
				Description: "Second command",
				Subcommands: []Subcommand{
					{Name: "sub2", Description: "Subcommand 2", Parameters: []Parameter{}},
				},
			},
		}

		tmpDir := t.TempDir()

		err := GenerateCode(definitions, tmpDir)
		require.NoError(t, err)

		// Verify both command files were created
		assert.FileExists(t, filepath.Join(tmpDir, "cmd1_gen.go"))
		assert.FileExists(t, filepath.Join(tmpDir, "cmd2_gen.go"))

		// Verify registry includes both commands
		registryContent, err := os.ReadFile(filepath.Join(tmpDir, "registry_gen.go"))
		require.NoError(t, err)

		registryStr := string(registryContent)
		assert.Contains(t, registryStr, "RegisterCmd1Sub1Tool")
		assert.Contains(t, registryStr, "RegisterCmd2Sub2Tool")
	})

	t.Run("handles parameters with different types", func(t *testing.T) {
		definitions := []CommandDefinition{
			{
				Command:     "types",
				Description: "Test types",
				Subcommands: []Subcommand{
					{
						Name:        "run",
						Description: "Run with various types",
						Parameters: []Parameter{
							{Name: "str", Type: "string", Description: "String param"},
							{Name: "num", Type: "integer", Description: "Integer param"},
							{Name: "flag", Type: "boolean", Description: "Boolean param"},
							{Name: "list", Type: "array", ItemType: "string", Description: "Array param"},
						},
					},
				},
			},
		}

		tmpDir := t.TempDir()

		err := GenerateCode(definitions, tmpDir)
		require.NoError(t, err)

		content, err := os.ReadFile(filepath.Join(tmpDir, "types_gen.go"))
		require.NoError(t, err)

		contentStr := string(content)
		// Check for field names and types (allowing for spacing/alignment)
		assert.Regexp(t, `Str\s+string`, contentStr, "Should have Str string field")
		assert.Regexp(t, `Num\s+int`, contentStr, "Should have Num int field")
		assert.Regexp(t, `Flag\s+bool`, contentStr, "Should have Flag bool field")
		assert.Regexp(t, `List\s+\[\]string`, contentStr, "Should have List []string field")
	})

	t.Run("returns error for invalid output directory", func(t *testing.T) {
		definitions := []CommandDefinition{
			{
				Command:     "test",
				Description: "Test",
				Subcommands: []Subcommand{
					{Name: "run", Description: "Run", Parameters: []Parameter{}},
				},
			},
		}

		// Use a file as output directory (should fail)
		tmpFile, err := os.CreateTemp("", "not-a-dir")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())
		tmpFile.Close()

		err = GenerateCode(definitions, tmpFile.Name())
		assert.Error(t, err)
	})
}

func TestGenerateCommandFile(t *testing.T) {
	t.Run("generates valid Go code", func(t *testing.T) {
		def := CommandDefinition{
			Command:     "example",
			Description: "Example command",
			Subcommands: []Subcommand{
				{
					Name:        "run",
					Description: "Run example",
					Parameters: []Parameter{
						{
							Name:        "verbose",
							Type:        "boolean",
							Flag:        "--verbose",
							Short:       "-v",
							Description: "Enable verbose output",
						},
					},
				},
			},
		}

		tmpDir := t.TempDir()

		err := generateCommandFile(def, tmpDir)
		require.NoError(t, err)

		// Verify file exists and content is valid
		filePath := filepath.Join(tmpDir, "example_gen.go")
		assert.FileExists(t, filePath)

		content, err := os.ReadFile(filePath)
		require.NoError(t, err)

		// Check for required code elements
		contentStr := string(content)
		assert.Contains(t, contentStr, "// Code generated by tools/gen. DO NOT EDIT.")
		assert.Contains(t, contentStr, "package generated")
		assert.Contains(t, contentStr, "func RegisterExampleRunTool")
		assert.Contains(t, contentStr, "mcp.AddTool")
		assert.Contains(t, contentStr, "exec.Execute")

		// Verify code compiles (basic syntax check)
		assert.NotContains(t, contentStr, "Warning: failed to format")
	})

	t.Run("handles positional parameters correctly", func(t *testing.T) {
		def := CommandDefinition{
			Command:     "pos",
			Description: "Positional test",
			Subcommands: []Subcommand{
				{
					Name:        "add",
					Description: "Add item",
					Parameters: []Parameter{
						{
							Name:        "name",
							Type:        "string",
							Description: "Item name (positional)",
							Positional:  true,
							Required:    true,
						},
						{
							Name:        "value",
							Type:        "string",
							Description: "Item value (positional)",
							Positional:  true,
							Required:    false,
						},
						{
							Name:        "force",
							Type:        "boolean",
							Flag:        "--force",
							Description: "Force operation",
						},
					},
				},
			},
		}

		tmpDir := t.TempDir()

		err := generateCommandFile(def, tmpDir)
		require.NoError(t, err)

		content, err := os.ReadFile(filepath.Join(tmpDir, "pos_gen.go"))
		require.NoError(t, err)

		contentStr := string(content)
		// Positional args should be appended directly to cmd
		assert.Contains(t, contentStr, "// Add positional argument: name")
		assert.Contains(t, contentStr, "// Add positional argument: value")
		// Flags should use if statements
		assert.Contains(t, contentStr, "if args.Force {")
	})
}

func TestGenerateRegistry(t *testing.T) {
	t.Run("generates registry with all tools", func(t *testing.T) {
		definitions := []CommandDefinition{
			{
				Command:     "cmd1",
				Description: "Command 1",
				Subcommands: []Subcommand{
					{Name: "list", Description: "List", Parameters: []Parameter{}},
					{Name: "create", Description: "Create", Parameters: []Parameter{}},
				},
			},
			{
				Command:     "cmd2",
				Description: "Command 2",
				Subcommands: []Subcommand{
					{Name: "run", Description: "Run", Parameters: []Parameter{}},
				},
			},
		}

		tmpDir := t.TempDir()

		err := generateRegistry(definitions, tmpDir)
		require.NoError(t, err)

		// Verify registry file
		registryPath := filepath.Join(tmpDir, "registry_gen.go")
		assert.FileExists(t, registryPath)

		content, err := os.ReadFile(registryPath)
		require.NoError(t, err)

		contentStr := string(content)
		assert.Contains(t, contentStr, "package generated")
		assert.Contains(t, contentStr, "func RegisterAllTools(")
		assert.Contains(t, contentStr, "server *mcp.Server")
		assert.Contains(t, contentStr, "exec *executor.Executor")

		// Verify all tools are registered
		assert.Contains(t, contentStr, "RegisterCmd1ListTool(server, exec)")
		assert.Contains(t, contentStr, "RegisterCmd1CreateTool(server, exec)")
		assert.Contains(t, contentStr, "RegisterCmd2RunTool(server, exec)")

		// Count the number of Register*Tool function calls (not just any "Register")
		registrationCount := strings.Count(contentStr, "RegisterCmd1ListTool(server, exec)") +
			strings.Count(contentStr, "RegisterCmd1CreateTool(server, exec)") +
			strings.Count(contentStr, "RegisterCmd2RunTool(server, exec)")
		assert.Equal(t, 3, registrationCount, "Should have 3 tool registration calls")
	})

	t.Run("generates empty registry for no definitions", func(t *testing.T) {
		definitions := []CommandDefinition{}

		tmpDir := t.TempDir()

		err := generateRegistry(definitions, tmpDir)
		require.NoError(t, err)

		content, err := os.ReadFile(filepath.Join(tmpDir, "registry_gen.go"))
		require.NoError(t, err)

		contentStr := string(content)
		assert.Contains(t, contentStr, "func RegisterAllTools(")
		// Should have no Register*Tool calls, only RegisterAllTools function definition
		// Count occurrences of "Tool(server, exec)" which indicates tool registration calls
		toolCallCount := strings.Count(contentStr, "Tool(server, exec)")
		assert.Equal(t, 0, toolCallCount, "Should have no tool registration calls for empty definitions")
	})
}

func TestGenerateCode_RealDefinitions(t *testing.T) {
	t.Run("generates code for actual project definitions", func(t *testing.T) {
		definitionsDir := "../../internal/commands/definitions"

		// Skip if definitions directory doesn't exist
		if _, err := os.Stat(definitionsDir); os.IsNotExist(err) {
			t.Skip("Definitions directory not found")
		}

		// Parse real definitions
		definitions, err := ParseDefinitions(definitionsDir)
		require.NoError(t, err)
		require.Greater(t, len(definitions), 0)

		// Generate code in temp directory
		tmpDir := t.TempDir()

		err = GenerateCode(definitions, tmpDir)
		require.NoError(t, err)

		// Verify registry was created
		registryPath := filepath.Join(tmpDir, "registry_gen.go")
		require.FileExists(t, registryPath)

		// Count generated files (should be one per command + registry)
		files, err := filepath.Glob(filepath.Join(tmpDir, "*_gen.go"))
		require.NoError(t, err)
		assert.Equal(t, len(definitions)+1, len(files), "Should generate one file per command plus registry")

		// Verify each generated file is valid Go code
		for _, file := range files {
			content, err := os.ReadFile(file)
			require.NoError(t, err)

			contentStr := string(content)
			assert.Contains(t, contentStr, "package generated")
			assert.NotContains(t, contentStr, "Warning: failed to format")

			// Check for required imports
			if !strings.Contains(file, "registry_gen.go") {
				assert.Contains(t, contentStr, `"context"`)
				assert.Contains(t, contentStr, `"github.com/khalideidoo/mcp-go-gh/internal/executor"`)
				assert.Contains(t, contentStr, `"github.com/modelcontextprotocol/go-sdk/mcp"`)
			}
		}
	})
}
