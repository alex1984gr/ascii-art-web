package tests

import (
	"errors"
	"strings"
	"testing"

	"ascii-art/pipeline"
)

func TestWebRenderASCII_Success(t *testing.T) {
	out, err := pipeline.RenderASCII("hello", "standard")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if strings.TrimSpace(out) == "" {
		t.Fatal("expected non-empty ascii output")
	}
}

func TestWebRenderASCII_InvalidBanner(t *testing.T) {
	_, err := pipeline.RenderASCII("hello", "unknown")
	if !errors.Is(err, pipeline.ErrInvalidBanner) {
		t.Fatalf("expected ErrInvalidBanner, got %v", err)
	}
}

func TestWebRenderASCII_InvalidInput(t *testing.T) {
	_, err := pipeline.RenderASCII("", "standard")
	if !errors.Is(err, pipeline.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}
