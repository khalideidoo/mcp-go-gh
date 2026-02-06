package main

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGenerateCode_StdoutPaths tests that the code generation prints to stdout.
func TestGenerateCode_StdoutPaths(t *testing.T) {
	t.Run("prints success messages to stdout", func(t *testing.T) {
		// Capture stdout
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		tmpDir := t.TempDir()
		definitions := []CommandDefinition{
			{
				Command: "test",
				Subcommands: []Subcommand{
					{
						Name:        "list",
						Description: "Test command",
						Parameters:  []Parameter{},
					},
				},
			},
		}

		err := GenerateCode(definitions, tmpDir)
		require.NoError(t, err)

		// Restore stdout
		w.Close()
		os.Stdout = oldStdout

		// Read captured output
		out, _ := io.ReadAll(r)
		output := string(out)

		// Verify output contains generation messages
		assert.Contains(t, output, "Generated", "should print generation messages")
		assert.Contains(t, output, "test_gen.go", "should mention generated file")
		assert.Contains(t, output, "registry_gen.go", "should mention registry file")
	})
}

func TestGenerateCommandFile_StdoutSuccess(t *testing.T) {
	t.Run("prints filename on success", func(t *testing.T) {
		// Capture stdout
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		tmpDir := t.TempDir()
		def := CommandDefinition{
			Command: "example",
			Subcommands: []Subcommand{
				{
					Name:        "list",
					Description: "List items",
					Parameters:  []Parameter{},
				},
			},
		}

		err := generateCommandFile(def, tmpDir)
		require.NoError(t, err)

		// Restore stdout
		w.Close()
		os.Stdout = oldStdout

		// Read captured output
		out, _ := io.ReadAll(r)
		output := string(out)

		assert.Contains(t, output, "Generated", "should print success message")
		assert.Contains(t, output, "example_gen.go", "should print filename")
	})
}

func TestGenerateRegistry_StdoutSuccess(t *testing.T) {
	t.Run("prints registry filename on success", func(t *testing.T) {
		// Capture stdout
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		tmpDir := t.TempDir()
		definitions := []CommandDefinition{
			{
				Command: "cmd1",
				Subcommands: []Subcommand{
					{Name: "sub1", Description: "Test"},
				},
			},
		}

		err := generateRegistry(definitions, tmpDir)
		require.NoError(t, err)

		// Restore stdout
		w.Close()
		os.Stdout = oldStdout

		// Read captured output
		out, _ := io.ReadAll(r)
		output := string(out)

		assert.Contains(t, output, "Generated", "should print success message")
		assert.Contains(t, output, "registry_gen.go", "should print registry filename")
	})
}

// TestGenerateCode_ComplexCommand tests more complex scenarios.
func TestGenerateCode_ComplexCommand(t *testing.T) {
	t.Run("handles command with multiple subcommands and parameter types", func(t *testing.T) {
		tmpDir := t.TempDir()

		definitions := []CommandDefinition{
			{
				Command:     "complex",
				Description: "A complex command",
				Subcommands: []Subcommand{
					{
						Name:        "create",
						Description: "Create something",
						Parameters: []Parameter{
							{Name: "name", Type: "string", Description: "Name", Required: true, Positional: true},
							{Name: "force", Type: "boolean", Description: "Force"},
							{Name: "tags", Type: "array", ItemType: "string", Description: "Tags"},
							{Name: "count", Type: "integer", Description: "Count"},
							{Name: "config", Type: "map", Description: "Config map"},
						},
					},
					{
						Name:        "delete",
						Description: "Delete something",
						Parameters: []Parameter{
							{Name: "id", Type: "string", Description: "ID", Positional: true},
							{Name: "confirm", Type: "boolean", Description: "Confirm deletion"},
						},
					},
				},
			},
		}

		err := GenerateCode(definitions, tmpDir)
		assert.NoError(t, err, "should handle complex command definitions")
	})
}
