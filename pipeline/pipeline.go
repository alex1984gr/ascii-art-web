// Package pipeline orchestrates the ASCII art generation and reverse decoding process
package pipeline

import (
	"fmt"     // Formatted I/O functions for printing errors and output
	"io"      // Basic I/O interfaces like Writer
	"os"      // Operating system functions for stderr access
	"strings" // String manipulation functions like ReplaceAll
)

// Run executes the complete ASCII art pipeline from input to output.
// It takes command-line arguments and an output writer, returns 0 for success or 1 for error.
func Run(args []string, stdout io.Writer) int {
	// Parse command-line arguments into a configuration struct
	cfg, err := parseArgs(args)
	// If argument parsing failed, print usage message and exit with error code
	if err != nil {
		fmt.Fprintln(os.Stderr, usageMessage)
		return 1
	}

	// Check if reverse mode is enabled (converting ASCII art back to text)
	if cfg.reverseFile != "" {
		// Attempt to decode ASCII art from the specified file
		text, err := ReverseFromFile(cfg.reverseFile)
		// If reverse decoding failed, print error and exit
		if err != nil {
			fmt.Fprintln(os.Stderr, "reverse error:", err)
			return 1
		}
		// Write the decoded text to stdout
		if _, err := fmt.Fprintln(stdout, text); err != nil {
			// If writing failed, print error and exit
			fmt.Fprintln(os.Stderr, "failed writing reverse output:", err)
			return 1
		}
		// Successfully completed reverse mode, return success code
		return 0
	}

	// Replace literal \n sequences with actual newline characters in the input
	cfg.input = strings.ReplaceAll(cfg.input, "\\n", "\n")

	// Validate the input string (check for empty, length, invalid characters)
	if err := ValidateInput(cfg.input); err != nil {
		// If validation failed, print error and exit
		fmt.Fprintln(os.Stderr, "invalid input:", err)
		return 1
	}

	// Load the ASCII art font/banner file (standard, shadow, or thinkertoy)
	banner, err := LoadBanner(cfg.font)
	// If banner loading failed, print error and exit
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed loading banner:", err)
		return 1
	}

	// Get the terminal width for alignment calculations
	terminalWidth := getTerminalWidth()
	// If justify alignment is requested, pre-process the input to add spacing
	if cfg.align == "justify" {
		// Modify input text to add extra spaces for justified output
		cfg.input = justifyInputForTerminal(cfg.input, banner, terminalWidth)
	}

	// Break input string into individual character tokens
	tokens := Tokenize(cfg.input)
	// Convert tokens into ASCII art lines (8 lines per character row)
	lines := RenderLines(tokens, banner)
	// Apply horizontal alignment (left, center, right) to the rendered lines
	lines = applyAlignment(lines, cfg.align, terminalWidth)

	// If color option was specified, apply ANSI color codes to the output
	if cfg.colorName != "" {
		// Color the lines (either entire lines or specific substring)
		lines, err = ColorLinesWithBanner(lines, cfg.colorName, cfg.substring, banner)
		// If coloring failed, print error and exit
		if err != nil {
			fmt.Fprintln(os.Stderr, "color error:", err)
			return 1
		}
	}

	// Set the output writer to stdout by default
	var w io.Writer = stdout
	// If output file was specified, create the file and use it as writer
	if cfg.outFile != "" {
		// Create or overwrite the output file
		f, err := os.Create(cfg.outFile)
		// If file creation failed, print error and exit
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed creating output file:", err)
			return 1
		}
		// Ensure file is closed when function returns
		defer f.Close()
		// Use the file as the output writer
		w = f
	}

	// Write the final ASCII art lines to the output writer
	if err := WriteOutput(lines, w); err != nil {
		// If writing failed, print error and exit
		fmt.Fprintln(os.Stderr, "failed writing output:", err)
		return 1
	}

	// Successfully completed, return success code
	return 0
}
