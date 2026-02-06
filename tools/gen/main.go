package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// Parse command line flags
	definitionsDir := flag.String("definitions", "internal/commands/definitions", "Directory containing YAML definitions")
	outputDir := flag.String("output", "internal/commands/generated", "Output directory for generated code")
	flag.Parse()

	// Convert to absolute paths
	absDefDir, err := filepath.Abs(*definitionsDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving definitions directory: %v\n", err)
		os.Exit(1)
	}

	absOutDir, err := filepath.Abs(*outputDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving output directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Reading definitions from: %s\n", absDefDir)
	fmt.Printf("Writing generated code to: %s\n\n", absOutDir)

	// Parse YAML definitions
	definitions, err := ParseDefinitions(absDefDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing definitions: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Parsed %d command definition(s)\n", len(definitions))
	for _, def := range definitions {
		fmt.Printf("  - %s (%d subcommands)\n", def.Command, len(def.Subcommands))
	}
	fmt.Println()

	// Generate code
	if err := GenerateCode(definitions, absOutDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating code: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nCode generation completed successfully!")
}
