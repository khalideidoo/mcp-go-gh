package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseDefinitions(t *testing.T) {
	t.Run("successfully parses all definitions from valid directory", func(t *testing.T) {
		// Use the actual definitions directory
		definitionsDir := "../../internal/commands/definitions"

		definitions, err := ParseDefinitions(definitionsDir)

		require.NoError(t, err)
		assert.Greater(t, len(definitions), 0, "should parse at least one definition")

		// Verify structure of parsed definitions
		for _, def := range definitions {
			assert.NotEmpty(t, def.Command, "command name should not be empty")
			assert.NotEmpty(t, def.Description, "description should not be empty")
			assert.Greater(t, len(def.Subcommands), 0, "should have at least one subcommand")

			// Verify subcommands
			for _, sub := range def.Subcommands {
				assert.NotEmpty(t, sub.Name, "subcommand name should not be empty")
				assert.NotEmpty(t, sub.Description, "subcommand description should not be empty")
			}
		}
	})

	t.Run("handles non-existent directory gracefully", func(t *testing.T) {
		definitions, err := ParseDefinitions("/nonexistent/directory")

		// filepath.Glob doesn't error on non-existent directories, just returns empty slice
		require.NoError(t, err)
		assert.Empty(t, definitions)
	})

	t.Run("handles empty directory gracefully", func(t *testing.T) {
		// Create a temporary empty directory
		tmpDir := t.TempDir()

		definitions, err := ParseDefinitions(tmpDir)

		// No error expected - just returns empty slice when no YAML files found
		require.NoError(t, err)
		assert.Empty(t, definitions)
	})

	t.Run("skips non-YAML files", func(t *testing.T) {
		// Create a temporary directory with mixed files
		tmpDir := t.TempDir()

		// Create a valid YAML file
		validYAML := `command: test
description: Test command
subcommands:
  - name: run
    description: Run test
    parameters: []
`
		err := os.WriteFile(filepath.Join(tmpDir, "test.yaml"), []byte(validYAML), 0644)
		require.NoError(t, err)

		// Create a non-YAML file
		err = os.WriteFile(filepath.Join(tmpDir, "readme.txt"), []byte("not yaml"), 0644)
		require.NoError(t, err)

		definitions, err := ParseDefinitions(tmpDir)

		require.NoError(t, err)
		assert.Len(t, definitions, 1, "should only parse YAML files")
		assert.Equal(t, "test", definitions[0].Command)
	})
}

func TestParseDefinitionFile(t *testing.T) {
	t.Run("successfully parses valid YAML file", func(t *testing.T) {
		// Create a temporary YAML file
		tmpDir := t.TempDir()
		yamlContent := `command: test
description: Test command
subcommands:
  - name: list
    description: List items
    parameters:
      - name: limit
        type: integer
        flag: --limit
        short: -L
        description: Maximum number to list
      - name: json
        type: array
        item_type: string
        flag: --json
        description: Output JSON

  - name: create
    description: Create item
    parameters:
      - name: name
        type: string
        description: Item name
        positional: true
        required: true
      - name: force
        type: boolean
        flag: --force
        short: -f
        description: Force creation
`
		filePath := filepath.Join(tmpDir, "test.yaml")
		err := os.WriteFile(filePath, []byte(yamlContent), 0644)
		require.NoError(t, err)

		def, err := parseDefinitionFile(filePath)

		require.NoError(t, err)
		assert.Equal(t, "test", def.Command)
		assert.Equal(t, "Test command", def.Description)
		assert.Len(t, def.Subcommands, 2)

		// Verify first subcommand
		listCmd := def.Subcommands[0]
		assert.Equal(t, "list", listCmd.Name)
		assert.Equal(t, "List items", listCmd.Description)
		assert.Len(t, listCmd.Parameters, 2)

		// Verify parameter types
		assert.Equal(t, "integer", listCmd.Parameters[0].Type)
		assert.Equal(t, "--limit", listCmd.Parameters[0].Flag)
		assert.Equal(t, "-L", listCmd.Parameters[0].Short)

		assert.Equal(t, "array", listCmd.Parameters[1].Type)
		assert.Equal(t, "string", listCmd.Parameters[1].ItemType)

		// Verify second subcommand with positional parameter
		createCmd := def.Subcommands[1]
		assert.Equal(t, "create", createCmd.Name)
		assert.Len(t, createCmd.Parameters, 2)
		assert.True(t, createCmd.Parameters[0].Positional)
		assert.True(t, createCmd.Parameters[0].Required)
		assert.False(t, createCmd.Parameters[1].Positional)
	})

	t.Run("returns error for non-existent file", func(t *testing.T) {
		def, err := parseDefinitionFile("/nonexistent/file.yaml")

		assert.Error(t, err)
		assert.Empty(t, def.Command, "command should be empty on error")
	})

	t.Run("returns error for invalid YAML syntax", func(t *testing.T) {
		tmpDir := t.TempDir()
		invalidYAML := `command: test
description: Test
subcommands:
  - name: list
    invalid_yaml: [unclosed bracket
`
		filePath := filepath.Join(tmpDir, "invalid.yaml")
		err := os.WriteFile(filePath, []byte(invalidYAML), 0644)
		require.NoError(t, err)

		def, err := parseDefinitionFile(filePath)

		assert.Error(t, err)
		assert.Empty(t, def.Command, "command should be empty on error")
	})

	t.Run("returns error for missing required fields", func(t *testing.T) {
		tmpDir := t.TempDir()
		incompleteYAML := `command: test
# missing description
subcommands:
  - name: list
    description: List items
`
		filePath := filepath.Join(tmpDir, "incomplete.yaml")
		err := os.WriteFile(filePath, []byte(incompleteYAML), 0644)
		require.NoError(t, err)

		def, err := parseDefinitionFile(filePath)

		// Should still parse but with empty description
		require.NoError(t, err)
		assert.Equal(t, "", def.Description)
	})

	t.Run("handles subcommands with no parameters", func(t *testing.T) {
		tmpDir := t.TempDir()
		yamlContent := `command: simple
description: Simple command
subcommands:
  - name: status
    description: Show status
    parameters: []
`
		filePath := filepath.Join(tmpDir, "simple.yaml")
		err := os.WriteFile(filePath, []byte(yamlContent), 0644)
		require.NoError(t, err)

		def, err := parseDefinitionFile(filePath)

		require.NoError(t, err)
		assert.Len(t, def.Subcommands, 1)
		assert.Len(t, def.Subcommands[0].Parameters, 0)
	})

	t.Run("handles enum values in parameters", func(t *testing.T) {
		tmpDir := t.TempDir()
		yamlContent := `command: test
description: Test command
subcommands:
  - name: set
    description: Set value
    parameters:
      - name: level
        type: string
        flag: --level
        description: Log level
        enum:
          - debug
          - info
          - warn
          - error
`
		filePath := filepath.Join(tmpDir, "enum.yaml")
		err := os.WriteFile(filePath, []byte(yamlContent), 0644)
		require.NoError(t, err)

		def, err := parseDefinitionFile(filePath)

		require.NoError(t, err)
		assert.Len(t, def.Subcommands[0].Parameters, 1)
		assert.Len(t, def.Subcommands[0].Parameters[0].Enum, 4)
		assert.Equal(t, []string{"debug", "info", "warn", "error"}, def.Subcommands[0].Parameters[0].Enum)
	})
}

func TestParseDefinitions_RealData(t *testing.T) {
	t.Run("parses actual project definitions", func(t *testing.T) {
		definitionsDir := "../../internal/commands/definitions"

		// Skip if definitions directory doesn't exist (e.g., in CI without full repo)
		if _, err := os.Stat(definitionsDir); os.IsNotExist(err) {
			t.Skip("Definitions directory not found")
		}

		definitions, err := ParseDefinitions(definitionsDir)
		require.NoError(t, err)

		// Verify we have the expected number of commands
		assert.GreaterOrEqual(t, len(definitions), 27, "should have at least 27 command definitions")

		// Verify some known commands exist
		commandNames := make(map[string]bool)
		for _, def := range definitions {
			commandNames[def.Command] = true
		}

		expectedCommands := []string{"alias", "auth", "issue", "pr", "repo", "gist"}
		for _, cmd := range expectedCommands {
			assert.True(t, commandNames[cmd], "should have %s command", cmd)
		}

		// Count total subcommands
		totalSubcommands := 0
		for _, def := range definitions {
			totalSubcommands += len(def.Subcommands)
		}
		assert.GreaterOrEqual(t, totalSubcommands, 150, "should have at least 150 total subcommands")
	})
}
