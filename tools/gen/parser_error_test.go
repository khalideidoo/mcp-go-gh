package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseDefinitions_ErrorPaths(t *testing.T) {
	t.Run("nonexistent directory", func(t *testing.T) {
		nonexistentDir := "/path/that/does/not/exist/anywhere"

		definitions, err := ParseDefinitions(nonexistentDir)
		assert.NoError(t, err, "glob on nonexistent dir returns empty, not error")
		assert.Empty(t, definitions, "should return empty definitions")
	})

	t.Run("directory with no yaml files", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create some non-YAML files
		err := os.WriteFile(filepath.Join(tmpDir, "readme.txt"), []byte("text"), 0644)
		require.NoError(t, err)

		definitions, err := ParseDefinitions(tmpDir)
		assert.NoError(t, err)
		assert.Empty(t, definitions, "should return empty when no YAML files found")
	})

	t.Run("invalid YAML file", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create an invalid YAML file
		invalidYAML := filepath.Join(tmpDir, "invalid.yaml")
		err := os.WriteFile(invalidYAML, []byte("not: valid: yaml: content:"), 0644)
		require.NoError(t, err)

		definitions, err := ParseDefinitions(tmpDir)
		assert.Error(t, err, "should fail on invalid YAML")
		assert.Contains(t, err.Error(), "failed to parse")
		assert.Nil(t, definitions)
	})

	t.Run("unreadable YAML file", func(t *testing.T) {
		if os.Getuid() == 0 {
			t.Skip("skipping permission test when running as root")
		}

		tmpDir := t.TempDir()

		// Create a YAML file and make it unreadable
		unreadableFile := filepath.Join(tmpDir, "unreadable.yaml")
		err := os.WriteFile(unreadableFile, []byte("command: test"), 0644)
		require.NoError(t, err)

		// Remove read permissions
		err = os.Chmod(unreadableFile, 0000)
		require.NoError(t, err)
		defer os.Chmod(unreadableFile, 0644) // Restore for cleanup

		definitions, err := ParseDefinitions(tmpDir)
		assert.Error(t, err, "should fail on unreadable file")
		assert.Contains(t, err.Error(), "failed to parse")
		assert.Nil(t, definitions)
	})

	t.Run("valid YAML but invalid structure", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create YAML that's valid but doesn't match CommandDefinition structure
		invalidStructure := filepath.Join(tmpDir, "invalid-structure.yaml")
		yamlContent := `
command: test
subcommands:
  - name: list
    parameters:
      - name: flag
        type: unknown_type_that_is_still_valid_yaml
        nested:
          deeply:
            broken: structure
`
		err := os.WriteFile(invalidStructure, []byte(yamlContent), 0644)
		require.NoError(t, err)

		// This should parse successfully as YAML unmarshaling is lenient
		// It will just ignore unknown fields
		definitions, err := ParseDefinitions(tmpDir)
		assert.NoError(t, err, "YAML unmarshaling is lenient with extra fields")
		require.Len(t, definitions, 1)
		assert.Equal(t, "test", definitions[0].Command)
	})
}

func TestParseDefinitionFile_ErrorPaths(t *testing.T) {
	t.Run("file does not exist", func(t *testing.T) {
		nonexistentFile := "/path/to/nonexistent/file.yaml"

		def, err := parseDefinitionFile(nonexistentFile)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read file")
		assert.Equal(t, CommandDefinition{}, def)
	})

	t.Run("corrupted YAML syntax", func(t *testing.T) {
		tmpDir := t.TempDir()
		corruptedFile := filepath.Join(tmpDir, "corrupted.yaml")

		// Write completely invalid YAML syntax
		corruptedYAML := `
command: test
subcommands: [
  name: broken
  - another: broken
} this is not yaml
`
		err := os.WriteFile(corruptedFile, []byte(corruptedYAML), 0644)
		require.NoError(t, err)

		def, err := parseDefinitionFile(corruptedFile)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to unmarshal YAML")
		assert.Equal(t, CommandDefinition{}, def)
	})

	t.Run("valid YAML parses successfully", func(t *testing.T) {
		tmpDir := t.TempDir()
		validFile := filepath.Join(tmpDir, "valid.yaml")

		validYAML := `
command: test
subcommands:
  - name: list
    description: List items
    parameters:
      - name: all
        type: boolean
        description: Show all
`
		err := os.WriteFile(validFile, []byte(validYAML), 0644)
		require.NoError(t, err)

		def, err := parseDefinitionFile(validFile)
		assert.NoError(t, err)
		assert.Equal(t, "test", def.Command)
		require.Len(t, def.Subcommands, 1)
		assert.Equal(t, "list", def.Subcommands[0].Name)
		assert.Equal(t, "List items", def.Subcommands[0].Description)
	})

	t.Run("empty file", func(t *testing.T) {
		tmpDir := t.TempDir()
		emptyFile := filepath.Join(tmpDir, "empty.yaml")

		err := os.WriteFile(emptyFile, []byte(""), 0644)
		require.NoError(t, err)

		def, err := parseDefinitionFile(emptyFile)
		// Empty YAML is valid and unmarshals to zero value
		assert.NoError(t, err)
		assert.Equal(t, "", def.Command)
		assert.Empty(t, def.Subcommands)
	})
}

func TestParseDefinitions_MultipleFiles(t *testing.T) {
	t.Run("mix of valid and invalid files", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create one valid file
		validFile := filepath.Join(tmpDir, "valid.yaml")
		validYAML := `
command: valid
subcommands:
  - name: test
    description: Test
`
		err := os.WriteFile(validFile, []byte(validYAML), 0644)
		require.NoError(t, err)

		// Create one invalid file with truly broken YAML syntax
		invalidFile := filepath.Join(tmpDir, "invalid.yaml")
		// Use tabs in the wrong place and mismatched brackets to create truly invalid YAML
		err = os.WriteFile(invalidFile, []byte("command: test\nsubcommands: [\n  - name: bad\n}"), 0644)
		require.NoError(t, err)

		// Should fail on the invalid file
		definitions, err := ParseDefinitions(tmpDir)
		if assert.Error(t, err, "should fail when any file is invalid") {
			assert.Contains(t, err.Error(), "failed to parse")
		}
		assert.Nil(t, definitions)
	})

	t.Run("multiple valid files parse successfully", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create multiple valid files
		for i := 1; i <= 3; i++ {
			filename := filepath.Join(tmpDir, "cmd"+string(rune('0'+i))+".yaml")
			yamlContent := "command: cmd" + string(rune('0'+i)) + "\nsubcommands: []"
			err := os.WriteFile(filename, []byte(yamlContent), 0644)
			require.NoError(t, err)
		}

		definitions, err := ParseDefinitions(tmpDir)
		assert.NoError(t, err)
		assert.Len(t, definitions, 3, "should parse all valid files")
	})
}
