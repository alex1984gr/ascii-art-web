// Package pipeline contains the reverse decoding logic to convert ASCII art back to text
package pipeline

import (
	"bufio"  // Buffered I/O for efficient file reading line by line
	"fmt"    // Formatted I/O functions for error messages
	"os"     // Operating system functions for file operations
	"strings" // String manipulation functions like Join, Repeat
)

// ReverseFromFile converts ASCII art stored in a file back to normal text.
// It tries all supported banners and returns the first successful decode.
func ReverseFromFile(fileName string) (string, error) {
	// Read all lines from the ASCII art file
	lines, err := readAsciiFileLines(fileName)
	// If file reading failed, return the error
	if err != nil {
		return "", err
	}

	// Check if the file is empty
	if len(lines) == 0 {
		return "", fmt.Errorf("reverse file is empty")
	}

	// List of all supported banner fonts to try
	banners := []string{"standard", "shadow", "thinkertoy"}
	// Try each banner font until one successfully decodes the ASCII art
	for _, name := range banners {
		// Load the banner font
		banner, err := LoadBanner(name)
		// If banner loading failed, skip to next banner
		if err != nil {
			continue
		}
		// Attempt to decode the ASCII art lines using this banner
		decoded, err := decodeAsciiLines(lines, banner)
		// If decoding succeeded, return the decoded text
		if err == nil {
			return decoded, nil
		}
	}

	// If no banner successfully decoded the ASCII art, return an error
	return "", fmt.Errorf("failed to reverse ascii art: no matching banner")
}

// readAsciiFileLines reads a file and returns its lines as a slice of strings.
// It removes trailing empty lines from the end of the file.
func readAsciiFileLines(fileName string) ([]string, error) {
	// Open the file for reading
	f, err := os.Open(fileName)
	// If file opening failed, return the error
	if err != nil {
		return nil, err
	}
	// Ensure file is closed when function returns
	defer f.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(f)
	// Initialize slice to store all lines
	var lines []string
	// Read each line from the file
	for scanner.Scan() {
		// Append the current line to the slice
		lines = append(lines, scanner.Text())
	}
	// Check if scanner encountered any errors during reading
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Remove trailing empty lines from the end
	for len(lines) > 0 && lines[len(lines)-1] == "" {
		// Trim the last line if it's empty
		lines = lines[:len(lines)-1]
	}

	// Return the lines without trailing empty lines
	return lines, nil
}

// decodeAsciiLines converts ASCII art lines back to normal text using the provided banner.
// ASCII art must have a height that's a multiple of 8 (each character is 8 lines tall).
func decodeAsciiLines(lines []string, banner map[string][]string) (string, error) {
	// Validate that the number of lines is a multiple of 8
	if len(lines)%8 != 0 {
		return "", fmt.Errorf("invalid ascii art height")
	}

	// Initialize slice to store decoded text lines
	var result []string
	// Process the ASCII art in blocks of 8 lines (one block per text line)
	for i := 0; i < len(lines); i += 8 {
		// Extract the current 8-line block
		block := lines[i : i+8]
		// Decode this block into a line of text
		line, err := decodeBlock(block, banner)
		// If decoding failed, return the error
		if err != nil {
			return "", err
		}
		// Add the decoded line to the result
		result = append(result, line)
	}

	// Join all decoded lines with newlines and return
	return strings.Join(result, "\n"), nil
}

// decodeBlock decodes a single 8-line block of ASCII art into a line of text.
// It matches character glyphs from left to right across the block.
func decodeBlock(block []string, banner map[string][]string) (string, error) {
	// Find the maximum width among all 8 lines in the block
	maxWidth := 0
	for _, r := range block {
		// Update maxWidth if current line is wider
		if len(r) > maxWidth {
			maxWidth = len(r)
		}
	}

	// Pad all lines to the same width with spaces (for consistent matching)
	for i := range block {
		// If current line is shorter than maxWidth
		if len(block[i]) < maxWidth {
			// Add spaces to the end to reach maxWidth
			block[i] += strings.Repeat(" ", maxWidth-len(block[i]))
		}
	}

	// Build the decoded text character by character
	var b strings.Builder
	// Current position in the block (column index)
	pos := 0
	// Process the block from left to right
	for pos < maxWidth {
		// Try to match a character glyph at the current position
		matched, width, ok := matchGlyphAt(block, pos, banner)
		// If no glyph matched, return an error
		if !ok {
			return "", fmt.Errorf("cannot decode ascii art at column %d", pos)
		}
		// Add the matched character to the result
		b.WriteString(matched)
		// Move position forward by the width of the matched glyph
		pos += width
	}

	// Return the decoded line of text
	return b.String(), nil
}

// matchGlyphAt tries to match a character glyph at the specified position in the block.
// It returns the matched character, its width, and whether a match was found.
// It prefers wider matches (greedy matching) to handle overlapping patterns.
func matchGlyphAt(block []string, pos int, banner map[string][]string) (string, int, bool) {
	// Track the best (widest) match found so far
	bestChar := ""
	bestWidth := 0

	// Try all printable ASCII characters (space to tilde)
	for c := 32; c <= 126; c++ {
		// Convert ASCII code to string character
		ch := string(rune(c))
		// Look up the glyph for this character in the banner
		glyph, ok := banner[ch]
		// Skip if character not found or glyph has less than 8 lines
		if !ok || len(glyph) < 8 {
			continue
		}
		// Get the width of this glyph (width of first line)
		w := len(glyph[0])
		// Skip if glyph has zero width
		if w == 0 {
			continue
		}
		// Skip if glyph extends beyond the block width
		if pos+w > len(block[0]) {
			continue
		}

		// Check if this glyph matches at the current position
		match := true
		// Compare all 8 rows of the glyph
		for row := 0; row < 8; row++ {
			// Check if the block segment matches the glyph row
			if pos+w > len(block[row]) || block[row][pos:pos+w] != glyph[row] {
				// Mismatch found, this glyph doesn't match
				match = false
				break
			}
		}
		// If glyph didn't match, skip to next character
		if !match {
			continue
		}

		// If this glyph is wider than the best match so far, update best match
		if w > bestWidth {
			bestWidth = w
			bestChar = ch
		}
	}

	// If no character matched, return failure
	if bestChar == "" {
		return "", 0, false
	}
	// Return the best (widest) matching character and its width
	return bestChar, bestWidth, true
}
