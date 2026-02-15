// Package pipeline contains the ASCII art rendering logic for web requests
package pipeline

// Import required packages for error handling and string manipulation
import (
	"errors" // Functions to create and compare error values
	"fmt"    // Formatted I/O functions for error wrapping
	"strings" // String manipulation functions like ReplaceAll and Join
)

// Define custom error types that can be checked by the web handler
var (
	ErrInvalidInput  = errors.New("invalid input")  // Error for invalid user input
	ErrInvalidBanner = errors.New("invalid banner") // Error for invalid banner selection
)

// RenderASCII is the main entry point for web requests to generate ASCII art.
// It takes user input text and a banner name, validates them, and returns the rendered ASCII art.
func RenderASCII(input, banner string) (string, error) {
	// If no banner was specified, default to "standard"
	if banner == "" {
		banner = "standard"
	}
	// Check if the banner name is valid (standard, shadow, or thinkertoy)
	if !isBannerName(banner) {
		// Return error if banner name is not recognized
		return "", ErrInvalidBanner
	}

	// Replace literal "\\n" sequences with actual newline characters
	// This allows users to type \n in the form to create line breaks
	input = strings.ReplaceAll(input, "\\n", "\n")
	// Validate the input text (check length, empty, invalid characters)
	if err := ValidateInput(input); err != nil {
		// Wrap the validation error with ErrInvalidInput for proper error handling
		return "", fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	// Load the banner file and create a map of characters to ASCII art glyphs
	fontMap, err := LoadBanner(banner)
	// If banner loading fails (file not found, etc.), return the error
	if err != nil {
		return "", err
	}

	// Tokenize the input into individual characters and render them as ASCII art
	// Tokenize splits the string into runes, RenderLines converts them to 8-line ASCII art
	lines := RenderLines(Tokenize(input), fontMap)
	// If no lines were generated (empty input after processing), return empty string
	if len(lines) == 0 {
		return "", nil
	}

	// Join all ASCII art lines with newlines and add a trailing newline
	return strings.Join(lines, "\n") + "\n", nil
}
