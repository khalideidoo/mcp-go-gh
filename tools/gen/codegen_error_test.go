package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateCode_ErrorPaths(t *testing.T) {
	t.Run("invalid output directory", func(t *testing.T) {
		// Create a file where we expect a directory
		tmpFile := filepath.Join(t.TempDir(), "not-a-dir")
		err := os.WriteFile(tmpFile, []byte("content"), 0644)
		require.NoError(t, err)

		// Try to use the file as a directory - should fail to create subdirectories
		invalidDir := filepath.Join(tmpFile, "subdir")
		definitions := []CommandDefinition{
			{Command: "test", Subcommands: []Subcommand{}},
		}

		err = GenerateCode(definitions, invalidDir)
		assert.Error(t, err, "should fail when output directory cannot be created")
		assert.Contains(t, err.Error(), "failed to create output directory")
	})

	t.Run("write permission denied", func(t *testing.T) {
		if os.Getuid() == 0 {
			t.Skip("skipping permission test when running as root")
		}

		// Create a read-only directory
		tmpDir := t.TempDir()
		readOnlyDir := filepath.Join(tmpDir, "readonly")
		err := os.Mkdir(readOnlyDir, 0555) // r-x r-x r-x (no write)
		require.NoError(t, err)

		// Make sure to restore permissions for cleanup
		defer os.Chmod(readOnlyDir, 0755)

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

		err = GenerateCode(definitions, readOnlyDir)
		assert.Error(t, err, "should fail when cannot write to directory")
		assert.Contains(t, err.Error(), "failed to generate code for test")
	})
}

func TestGenerateCommandFile_ErrorPaths(t *testing.T) {
	t.Run("write file error - read-only directory", func(t *testing.T) {
		if os.Getuid() == 0 {
			t.Skip("skipping permission test when running as root")
		}

		// Create a read-only directory
		tmpDir := t.TempDir()
		readOnlyDir := filepath.Join(tmpDir, "readonly")
		err := os.Mkdir(readOnlyDir, 0555)
		require.NoError(t, err)
		defer os.Chmod(readOnlyDir, 0755)

		def := CommandDefinition{
			Command: "test",
			Subcommands: []Subcommand{
				{
					Name:        "list",
					Description: "Test",
					Parameters:  []Parameter{},
				},
			},
		}

		err = generateCommandFile(def, readOnlyDir)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to write file")
	})

	t.Run("malformed template execution", func(t *testing.T) {
		// This tests the template execution path, though in practice
		// our templates should always execute successfully with valid data
		tmpDir := t.TempDir()

		// Create definition with fields that might cause template issues
		def := CommandDefinition{
			Command:     "test",
			Subcommands: []Subcommand{},
		}

		// With empty subcommands, the template should still execute
		// This ensures we're testing the template execution code path
		err := generateCommandFile(def, tmpDir)
		assert.NoError(t, err, "should handle empty subcommands gracefully")
	})
}

func TestGenerateRegistry_ErrorPaths(t *testing.T) {
	t.Run("write file error - read-only directory", func(t *testing.T) {
		if os.Getuid() == 0 {
			t.Skip("skipping permission test when running as root")
		}

		// Create a read-only directory
		tmpDir := t.TempDir()
		readOnlyDir := filepath.Join(tmpDir, "readonly")
		err := os.Mkdir(readOnlyDir, 0555)
		require.NoError(t, err)
		defer os.Chmod(readOnlyDir, 0755)

		definitions := []CommandDefinition{
			{
				Command: "test",
				Subcommands: []Subcommand{
					{Name: "list", Description: "Test"},
				},
			},
		}

		err = generateRegistry(definitions, readOnlyDir)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to write file")
	})

	t.Run("empty definitions list", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Test with empty definitions - should still generate valid registry
		err := generateRegistry([]CommandDefinition{}, tmpDir)
		assert.NoError(t, err, "should handle empty definitions gracefully")

		// Verify the file was created
		registryPath := filepath.Join(tmpDir, "registry_gen.go")
		_, err = os.Stat(registryPath)
		assert.NoError(t, err, "registry file should exist")
	})
}

func TestGenerateCode_IntegrationErrorCases(t *testing.T) {
	t.Run("nonexistent parent directory", func(t *testing.T) {
		// Use a path with nonexistent parent
		tmpDir := t.TempDir()
		deepPath := filepath.Join(tmpDir, "does", "not", "exist", "yet")

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

		// Should succeed - MkdirAll creates parent directories
		err := GenerateCode(definitions, deepPath)
		assert.NoError(t, err, "should create nested directories")

		// Verify files were created
		files, err := filepath.Glob(filepath.Join(deepPath, "*_gen.go"))
		require.NoError(t, err)
		assert.NotEmpty(t, files, "should have generated files")
	})
}
