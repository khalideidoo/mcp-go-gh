package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ParseDefinitions reads all YAML files from a directory.
func ParseDefinitions(dir string) ([]CommandDefinition, error) {
	var definitions []CommandDefinition

	files, err := filepath.Glob(filepath.Join(dir, "*.yaml"))
	if err != nil {
		return nil, fmt.Errorf("failed to glob YAML files: %w", err)
	}

	for _, file := range files {
		def, err := parseDefinitionFile(file)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s: %w", file, err)
		}
		definitions = append(definitions, def)
	}

	return definitions, nil
}

// parseDefinitionFile reads and parses a single YAML file.
func parseDefinitionFile(path string) (CommandDefinition, error) {
	// #nosec G304 -- path is from filepath.Glob, which is safe
	data, err := os.ReadFile(path)
	if err != nil {
		return CommandDefinition{}, fmt.Errorf("failed to read file: %w", err)
	}

	var def CommandDefinition
	if err := yaml.Unmarshal(data, &def); err != nil {
		return CommandDefinition{}, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return def, nil
}
