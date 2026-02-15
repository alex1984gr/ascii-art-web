// Package pipeline handles command-line argument parsing and validation for the ASCII art application
package pipeline

import (
	"fmt"     // Formatted I/O functions for creating error messages
	"strings" // String manipulation functions like HasPrefix, TrimPrefix, ToLower
)

// usageMessage is the help text displayed when arguments are invalid
const usageMessage = "Usage: go run . [OPTION]\n\nEX: go run . --reverse=<fileName>"

// runConfig holds all configuration options parsed from command-line arguments
type runConfig struct {
	font        string // Banner font name (standard, shadow, thinkertoy)
	outFile     string // Output file path (if --out or --output specified)
	input       string // Input text to convert to ASCII art
	colorName   string // Color name for ANSI coloring (if --color specified)
	substring   string // Substring to color (if specified with --color)
	align       string // Alignment type (left, center, right, justify)
	reverseFile string // File path for reverse mode (if --reverse specified)
	fontSet     bool   // Flag indicating if --font was explicitly set
}

// parseArgs parses command-line arguments and returns a runConfig struct or an error.
// It handles both flag-style arguments (--flag=value) and positional arguments.
func parseArgs(args []string) (runConfig, error) {
	// Initialize config with default values: standard font and left alignment
	cfg := runConfig{font: "standard", align: "left"}
	// Create a slice to collect non-flag positional arguments
	positionals := make([]string, 0, len(args))

	// Iterate through all command-line arguments
	for _, arg := range args {
		// Use a switch statement to handle different argument types
		switch {
		// Handle --font=<value> flag
		case strings.HasPrefix(arg, "--font="):
			// Extract the font name after the equals sign
			cfg.font = strings.TrimPrefix(arg, "--font=")
			// Mark that font was explicitly set (to distinguish from default)
			cfg.fontSet = true
			// Validate that font name is not empty
			if cfg.font == "" {
				return runConfig{}, fmt.Errorf("empty font")
			}
		// Handle invalid --font flag without equals sign
		case arg == "--font":
			return runConfig{}, fmt.Errorf("invalid font format")
		// Handle --out=<value> flag for output file
		case strings.HasPrefix(arg, "--out="):
			// Extract the output file path
			cfg.outFile = strings.TrimPrefix(arg, "--out=")
			// Validate that output file path is not empty
			if cfg.outFile == "" {
				return runConfig{}, fmt.Errorf("empty out file")
			}
		// Handle --output=<value> flag (alternative to --out)
		case strings.HasPrefix(arg, "--output="):
			// Extract the output file path
			cfg.outFile = strings.TrimPrefix(arg, "--output=")
			// Validate that output file path is not empty
			if cfg.outFile == "" {
				return runConfig{}, fmt.Errorf("empty output file")
			}
		// Handle invalid --out or --output flags without equals sign
		case arg == "--out" || arg == "--output":
			return runConfig{}, fmt.Errorf("invalid output format")
		// Handle --color=<value> flag for ANSI coloring
		case strings.HasPrefix(arg, "--color="):
			// Extract the color name
			cfg.colorName = strings.TrimPrefix(arg, "--color=")
			// Validate that color name is not empty
			if cfg.colorName == "" {
				return runConfig{}, fmt.Errorf("empty color")
			}
		// Handle invalid --color flag without equals sign
		case arg == "--color":
			return runConfig{}, fmt.Errorf("invalid color format")
		// Handle --align=<value> flag for text alignment
		case strings.HasPrefix(arg, "--align="):
			// Extract and convert alignment type to lowercase
			cfg.align = strings.ToLower(strings.TrimPrefix(arg, "--align="))
			// Validate that alignment type is one of the supported types
			if !isAlignType(cfg.align) {
				return runConfig{}, fmt.Errorf("invalid align type")
			}
		// Handle invalid --align flag without equals sign
		case arg == "--align":
			return runConfig{}, fmt.Errorf("invalid align format")
		// Handle --reverse=<value> flag for reverse mode (ASCII art to text)
		case strings.HasPrefix(arg, "--reverse="):
			// Extract the reverse file path
			cfg.reverseFile = strings.TrimPrefix(arg, "--reverse=")
			// Validate that reverse file path is not empty
			if cfg.reverseFile == "" {
				return runConfig{}, fmt.Errorf("empty reverse file")
			}
		// Handle invalid --reverse flag without equals sign
		case arg == "--reverse":
			return runConfig{}, fmt.Errorf("invalid reverse format")
		// Handle unknown flags (anything starting with --)
		case strings.HasPrefix(arg, "--"):
			return runConfig{}, fmt.Errorf("unknown option")
		// Handle positional arguments (non-flag arguments)
		default:
			// Add to positionals slice for later processing
			positionals = append(positionals, arg)
		}
	}

	// If reverse mode is enabled, validate that no positional arguments were provided
	if cfg.reverseFile != "" {
		// Reverse mode doesn't accept positional arguments
		if len(positionals) > 0 {
			return runConfig{}, fmt.Errorf("reverse mode takes no positional args")
		}
		// Return the config for reverse mode
		return cfg, nil
	}

	// If no color option was specified, handle positional arguments for normal mode
	if cfg.colorName == "" {
		// Check number of positional arguments
		switch len(positionals) {
		// Single argument: it's the input text
		case 1:
			cfg.input = positionals[0]
		// Two arguments: input text and banner name
		case 2:
			// If font was already set via --font flag, this is an error
			if cfg.fontSet {
				return runConfig{}, fmt.Errorf("too many args")
			}
			// Validate that second argument is a valid banner name
			if !isBannerName(positionals[1]) {
				return runConfig{}, fmt.Errorf("invalid banner")
			}
			// First argument is input text
			cfg.input = positionals[0]
			// Second argument is banner name
			cfg.font = positionals[1]
		// Any other number of arguments is invalid
		default:
			return runConfig{}, fmt.Errorf("invalid args")
		}
		// Return the config for normal mode without color
		return cfg, nil
	}

	// Color option was specified, handle positional arguments with color logic
	switch len(positionals) {
	// Single argument: it's the input text (color entire output)
	case 1:
		cfg.input = positionals[0]
	// Two arguments: could be (input, banner) or (substring, input)
	case 2:
		// If font wasn't explicitly set and second arg is a banner name
		if !cfg.fontSet && isBannerName(positionals[1]) {
			// First argument is input text
			cfg.input = positionals[0]
			// Second argument is banner name
			cfg.font = positionals[1]
		} else {
			// First argument is substring to color
			cfg.substring = positionals[0]
			// Second argument is input text
			cfg.input = positionals[1]
		}
	// Three arguments: substring, input, and banner name
	case 3:
		// If font was already set via --font flag, this is an error
		if cfg.fontSet {
			return runConfig{}, fmt.Errorf("too many args")
		}
		// Validate that third argument is a valid banner name
		if !isBannerName(positionals[2]) {
			return runConfig{}, fmt.Errorf("invalid banner")
		}
		// First argument is substring to color
		cfg.substring = positionals[0]
		// Second argument is input text
		cfg.input = positionals[1]
		// Third argument is banner name
		cfg.font = positionals[2]
	// Any other number of arguments is invalid
	default:
		return runConfig{}, fmt.Errorf("invalid args")
	}

	// Return the config for color mode
	return cfg, nil
}

// isBannerName checks if the given name is one of the supported banner fonts.
// Returns true if name is "standard", "shadow", or "thinkertoy".
func isBannerName(name string) bool {
	return name == "standard" || name == "shadow" || name == "thinkertoy"
}

// isAlignType checks if the given alignment is one of the supported types.
// Returns true if align is "left", "center", "right", or "justify".
func isAlignType(align string) bool {
	return align == "left" || align == "center" || align == "right" || align == "justify"
}
