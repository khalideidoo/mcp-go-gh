package main

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	typeString = "string"
)

// GenerateCode generates Go code for all command definitions.
func GenerateCode(definitions []CommandDefinition, outputDir string) error {
	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0750); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate code for each command
	for _, def := range definitions {
		if err := generateCommandFile(def, outputDir); err != nil {
			return fmt.Errorf("failed to generate code for %s: %w", def.Command, err)
		}
	}

	// Generate registry file
	if err := generateRegistry(definitions, outputDir); err != nil {
		return fmt.Errorf("failed to generate registry: %w", err)
	}

	return nil
}

// generateCommandFile generates a Go file for a single command group.
func generateCommandFile(def CommandDefinition, outputDir string) error {
	tmpl, err := template.New("command").Funcs(templateFuncs()).Parse(commandTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, def); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Format the generated code
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		// If formatting fails, write unformatted code for debugging
		fmt.Fprintf(os.Stderr, "Warning: failed to format %s: %v\n", def.Command, err)
		formatted = buf.Bytes()
	}

	// Write to file
	filename := filepath.Join(outputDir, fmt.Sprintf("%s_gen.go", def.Command))
	if err := os.WriteFile(filename, formatted, 0600); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("Generated %s\n", filename)
	return nil
}

// generateRegistry generates the registry.go file that registers all tools.
func generateRegistry(definitions []CommandDefinition, outputDir string) error {
	tmpl, err := template.New("registry").Funcs(templateFuncs()).Parse(registryTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, definitions); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to format registry: %v\n", err)
		formatted = buf.Bytes()
	}

	filename := filepath.Join(outputDir, "registry_gen.go")
	if err := os.WriteFile(filename, formatted, 0600); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("Generated %s\n", filename)
	return nil
}

// templateFuncs returns custom template functions.
func templateFuncs() template.FuncMap {
	return template.FuncMap{
		"toTitle":        toTitle,
		"toCamel":        toCamel,
		"toSnake":        toSnake,
		"goType":         goType,
		"jsonTag":        jsonTag,
		"schemaTag":      schemaTag,
		"hasPositional":  hasPositional,
		"nonPositional":  nonPositional,
		"positionalArgs": positionalArgs,
	}
}

// toTitle converts string to TitleCase.
func toTitle(s string) string {
	// Replace hyphens with underscores first, then split on underscores
	s = strings.ReplaceAll(s, "-", "_")
	parts := strings.Split(s, "_")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, "")
}

// toCamel converts string to camelCase.
func toCamel(s string) string {
	title := toTitle(s)
	if len(title) > 0 {
		return strings.ToLower(title[:1]) + title[1:]
	}
	return title
}

// toSnake converts string to snake_case.
func toSnake(s string) string {
	return strings.ReplaceAll(s, "-", "_")
}

// goType returns the Go type for a parameter type.
func goType(param Parameter) string {
	switch param.Type {
	case typeString:
		return typeString
	case "integer":
		return "int"
	case "boolean":
		return "bool"
	case "array":
		itemType := typeString
		if param.ItemType != "" {
			switch param.ItemType {
			case "integer":
				itemType = "int"
			}
		}
		return "[]" + itemType
	case "map":
		return "map[string]string"
	default:
		return "string"
	}
}

// jsonTag generates the JSON struct tag.
func jsonTag(param Parameter) string {
	name := toSnake(param.Name)
	tag := fmt.Sprintf(`json:"%s,omitempty"`, name)
	return tag
}

// schemaTag generates the jsonschema struct tag
// The google/jsonschema-go package expects just the description text, not key=value format.
func schemaTag(param Parameter) string {
	// Escape quotes in description to prevent tag parsing issues
	description := strings.ReplaceAll(param.Description, `"`, `'`)

	// The jsonschema tag should contain just the description text
	// Required fields are inferred from json:"field" vs json:"field,omitempty"
	// Enum validation is not supported via struct tags in google/jsonschema-go
	return fmt.Sprintf(`jsonschema:"%s"`, description)
}

// hasPositional checks if subcommand has positional arguments.
func hasPositional(sub Subcommand) bool {
	for _, param := range sub.Parameters {
		if param.Positional {
			return true
		}
	}
	return false
}

// nonPositional returns non-positional parameters.
func nonPositional(params []Parameter) []Parameter {
	var result []Parameter
	for _, param := range params {
		if !param.Positional {
			result = append(result, param)
		}
	}
	return result
}

// positionalArgs returns positional parameters.
func positionalArgs(params []Parameter) []Parameter {
	var result []Parameter
	for _, param := range params {
		if param.Positional {
			result = append(result, param)
		}
	}
	return result
}
