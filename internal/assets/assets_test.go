package assets

import (
	"os"
	"strings"
	"testing"

	"github.com/olle/docops/internal/cssmin"
)

func TestEmbeddedCSSIsMinified(t *testing.T) {
	if strings.ContainsAny(CSS, "\n\r\t") {
		t.Errorf("CSS still contains line breaks or tabs")
	}
	if strings.Contains(CSS, "  ") {
		t.Errorf("CSS still contains repeated spaces")
	}
	if strings.Contains(CSS, "/*") || strings.Contains(CSS, "*/") {
		t.Errorf("CSS still contains comments")
	}
	if !strings.Contains(CSS, "--font-sans:") || !strings.Contains(CSS, "@media (prefers-reduced-motion:reduce)") {
		t.Errorf("CSS lost recognizable content")
	}
}

func TestEmbeddedCSSMatchesMinifiedSource(t *testing.T) {
	data, err := os.ReadFile("style.css")
	if err != nil {
		t.Fatal(err)
	}

	want := cssmin.Minify(string(data))
	if CSS != want {
		t.Errorf("style.min.css is stale: got %d bytes, want %d; run go generate ./...", len(CSS), len(want))
	}
}
