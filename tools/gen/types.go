package main

// CommandDefinition represents a top-level gh command group
type CommandDefinition struct {
	Command     string       `yaml:"command"`
	Description string       `yaml:"description"`
	Subcommands []Subcommand `yaml:"subcommands"`
}

// Subcommand represents a specific gh subcommand
type Subcommand struct {
	Name        string      `yaml:"name"`
	Description string      `yaml:"description"`
	Parameters  []Parameter `yaml:"parameters"`
}

// Parameter represents a command parameter/flag
type Parameter struct {
	Name        string   `yaml:"name"`
	Type        string   `yaml:"type"`        // string, integer, boolean, array, map
	ItemType    string   `yaml:"item_type"`   // for array types
	Flag        string   `yaml:"flag"`        // --flag-name
	Short       string   `yaml:"short"`       // -f
	Description string   `yaml:"description"`
	Required    bool     `yaml:"required"`
	Positional  bool     `yaml:"positional"`  // positional argument
	Enum        []string `yaml:"enum"`        // valid values
}
