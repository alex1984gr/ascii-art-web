// Package pipeline contains all the processing stages for ASCII art generation
package pipeline

import (
	"errors" // Functions to create simple error values
	"fmt"    // Formatted I/O functions for error messages
)

// ValidateInput checks that the input string meets the following requirements:
// 1. The input is not empty.
// 2. The input does not exceed 10,000 runes (characters) in length.
// 3. The input contains only printable ASCII characters (32..126), plus tab/newline/carriage return.
// If any condition is violated, it returns a non-nil error describing the problem.
func ValidateInput(input string) error {
	// Count the total number of runes (Unicode characters) in the input string.
	// This accounts for multi-byte UTF-8 sequences as single runes.
	runeCount := 0
	// Loop through each rune in the input (range on string iterates over runes, not bytes)
	for range input {
		runeCount++ // Increment counter for each character
	}

	// Check if the input is empty (zero runes)
	if runeCount == 0 {
		return errors.New("input is empty")
	}

	// Check if the input exceeds the maximum allowed length of 10,000 runes
	if runeCount > 10000 {
		return errors.New("input too long")
	}

	// Iterate through each rune to validate individual characters
	for _, r := range input {
		// Allow tab/newline/carriage return as line/whitespace controls.
		if r == '\t' || r == '\n' || r == '\r' {
			continue
		}
		// Reject anything outside printable ASCII to preserve banner lookup behavior.
		if r < 32 || r > 126 {
			return fmt.Errorf("invalid character: 0x%x", r)
		}
	}

	// If all validations passed, return nil to indicate the input is valid
	return nil
}
