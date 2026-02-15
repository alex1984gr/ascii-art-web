package tests

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"ascii-art/pipeline"
)

func TestReverse_RunFromGeneratedFile(t *testing.T) {
	tempDir := t.TempDir()
	asciiFile := filepath.Join(tempDir, "art.txt")

	var out bytes.Buffer
	code := pipeline.Run([]string{"--output=" + asciiFile, "hello", "standard"}, &out)
	if code != 0 {
		t.Fatalf("expected ascii generation success, got %d", code)
	}

	var reversed bytes.Buffer
	code = pipeline.Run([]string{"--reverse=" + asciiFile}, &reversed)
	if code != 0 {
		t.Fatalf("expected reverse success, got %d", code)
	}

	got := strings.TrimSpace(reversed.String())
	if got != "hello" {
		t.Fatalf("expected reversed text 'hello', got %q", got)
	}
}

func TestReverse_InvalidFlagFormat(t *testing.T) {
	var out bytes.Buffer
	code := pipeline.Run([]string{"--reverse", "file.txt"}, &out)
	if code == 0 {
		t.Fatal("expected failure for invalid --reverse format")
	}
}

func TestReverse_RejectsPositionalArgs(t *testing.T) {
	var out bytes.Buffer
	code := pipeline.Run([]string{"--reverse=file.txt", "standard"}, &out)
	if code == 0 {
		t.Fatal("expected failure when --reverse is combined with positional args")
	}
}
