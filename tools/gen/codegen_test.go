package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToTitle(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple lowercase",
			input:    "hello",
			expected: "Hello",
		},
		{
			name:     "with underscores",
			input:    "hello_world",
			expected: "HelloWorld",
		},
		{
			name:     "with hyphens",
			input:    "field-list",
			expected: "FieldList",
		},
		{
			name:     "mixed hyphens and underscores",
			input:    "project-field_create",
			expected: "ProjectFieldCreate",
		},
		{
			name:     "already capitalized",
			input:    "AlreadyCapitalized",
			expected: "AlreadyCapitalized",
		},
		{
			name:     "single character",
			input:    "a",
			expected: "A",
		},
		{
			name:     "multiple hyphens",
			input:    "foo-bar-baz",
			expected: "FooBarBaz",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toTitle(tt.input)
			assert.Equal(t, tt.expected, result, "toTitle(%q) should return %q", tt.input, tt.expected)
		})
	}
}

func TestToCamel(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple lowercase",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "with underscores",
			input:    "hello_world",
			expected: "helloWorld",
		},
		{
			name:     "with hyphens",
			input:    "field-list",
			expected: "fieldList",
		},
		{
			name:     "single character",
			input:    "A",
			expected: "a",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toCamel(tt.input)
			assert.Equal(t, tt.expected, result, "toCamel(%q) should return %q", tt.input, tt.expected)
		})
	}
}

func TestToSnake(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple lowercase",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "with hyphens",
			input:    "field-list",
			expected: "field_list",
		},
		{
			name:     "multiple hyphens",
			input:    "foo-bar-baz",
			expected: "foo_bar_baz",
		},
		{
			name:     "already snake_case",
			input:    "already_snake",
			expected: "already_snake",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toSnake(tt.input)
			assert.Equal(t, tt.expected, result, "toSnake(%q) should return %q", tt.input, tt.expected)
		})
	}
}

func TestGoType(t *testing.T) {
	tests := []struct {
		name     string
		expected string
		param    Parameter
	}{
		{
			name:     "string type",
			param:    Parameter{Type: "string"},
			expected: "string",
		},
		{
			name:     "integer type",
			param:    Parameter{Type: "integer"},
			expected: "int",
		},
		{
			name:     "boolean type",
			param:    Parameter{Type: "boolean"},
			expected: "bool",
		},
		{
			name:     "array of strings",
			param:    Parameter{Type: "array", ItemType: "string"},
			expected: "[]string",
		},
		{
			name:     "array of integers",
			param:    Parameter{Type: "array", ItemType: "integer"},
			expected: "[]int",
		},
		{
			name:     "array without item type defaults to string",
			param:    Parameter{Type: "array"},
			expected: "[]string",
		},
		{
			name:     "map type",
			param:    Parameter{Type: "map"},
			expected: "map[string]string",
		},
		{
			name:     "unknown type defaults to string",
			param:    Parameter{Type: "unknown"},
			expected: "string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := goType(tt.param)
			assert.Equal(t, tt.expected, result, "goType should return %q for %+v", tt.expected, tt.param)
		})
	}
}

func TestJsonTag(t *testing.T) {
	tests := []struct {
		name     string
		expected string
		param    Parameter
	}{
		{
			name:     "simple parameter",
			param:    Parameter{Name: "title"},
			expected: `json:"title,omitempty"`,
		},
		{
			name:     "parameter with hyphen",
			param:    Parameter{Name: "field-list"},
			expected: `json:"field_list,omitempty"`,
		},
		{
			name:     "parameter with underscore",
			param:    Parameter{Name: "some_param"},
			expected: `json:"some_param,omitempty"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := jsonTag(tt.param)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSchemaTag(t *testing.T) {
	tests := []struct {
		name     string
		expected string
		param    Parameter
	}{
		{
			name:     "simple parameter",
			param:    Parameter{Description: "Test description"},
			expected: `jsonschema:"Test description"`,
		},
		{
			name:     "required parameter (required inferred from json tag)",
			param:    Parameter{Description: "Test", Required: true},
			expected: `jsonschema:"Test"`,
		},
		{
			name:     "parameter with enum (enum not supported in struct tags)",
			param:    Parameter{Description: "Test", Enum: []string{"opt1", "opt2"}},
			expected: `jsonschema:"Test"`,
		},
		{
			name: "parameter with all features",
			param: Parameter{
				Description: "Test",
				Required:    true,
				Enum:        []string{"a", "b"},
			},
			expected: `jsonschema:"Test"`,
		},
		{
			name:     "parameter with quotes in description",
			param:    Parameter{Description: `Test "quoted" value`},
			expected: `jsonschema:"Test 'quoted' value"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := schemaTag(tt.param)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHasPositional(t *testing.T) {
	tests := []struct {
		name     string
		sub      Subcommand
		expected bool
	}{
		{
			name: "has positional parameter",
			sub: Subcommand{
				Parameters: []Parameter{
					{Name: "file", Positional: true},
					{Name: "flag", Positional: false},
				},
			},
			expected: true,
		},
		{
			name: "no positional parameters",
			sub: Subcommand{
				Parameters: []Parameter{
					{Name: "flag1", Positional: false},
					{Name: "flag2", Positional: false},
				},
			},
			expected: false,
		},
		{
			name: "empty parameters",
			sub: Subcommand{
				Parameters: []Parameter{},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasPositional(tt.sub)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNonPositional(t *testing.T) {
	params := []Parameter{
		{Name: "file", Positional: true},
		{Name: "flag1", Positional: false},
		{Name: "target", Positional: true},
		{Name: "flag2", Positional: false},
	}

	result := nonPositional(params)
	require.Len(t, result, 2, "should have 2 non-positional parameters")
	assert.Equal(t, "flag1", result[0].Name)
	assert.Equal(t, "flag2", result[1].Name)
}

func TestPositionalArgs(t *testing.T) {
	params := []Parameter{
		{Name: "file", Positional: true},
		{Name: "flag1", Positional: false},
		{Name: "target", Positional: true},
		{Name: "flag2", Positional: false},
	}

	result := positionalArgs(params)
	require.Len(t, result, 2, "should have 2 positional parameters")
	assert.Equal(t, "file", result[0].Name)
	assert.Equal(t, "target", result[1].Name)
}

func TestTemplateFuncs(t *testing.T) {
	t.Run("returns all required template functions", func(t *testing.T) {
		funcs := templateFuncs()

		requiredFuncs := []string{
			"toTitle",
			"toCamel",
			"toSnake",
			"goType",
			"jsonTag",
			"schemaTag",
			"hasPositional",
			"nonPositional",
			"positionalArgs",
		}

		for _, name := range requiredFuncs {
			_, exists := funcs[name]
			assert.True(t, exists, "template function %q should exist", name)
		}
	})
}
