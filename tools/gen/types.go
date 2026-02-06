package main

// CommandDefinition represents a top-level gh command group.
type CommandDefinition struct {
	Command     string       `yaml:"command"`
	Description string       `yaml:"description"`
	Subcommands []Subcommand `yaml:"subcommands"`
}

// Subcommand represents a specific gh subcommand.
type Subcommand struct {
	Name        string      `yaml:"name"`
	Description string      `yaml:"description"`
	Parameters  []Parameter `yaml:"parameters"`
}

// Parameter represents a command parameter/flag.
type Parameter struct {
	Name        string   `yaml:"name"`
	Type        string   `yaml:"type"`
	ItemType    string   `yaml:"item_type"`
	Flag        string   `yaml:"flag"`
	Short       string   `yaml:"short"`
	Description string   `yaml:"description"`
	Enum        []string `yaml:"enum"`
	Required    bool     `yaml:"required"`
	Positional  bool     `yaml:"positional"`
}
