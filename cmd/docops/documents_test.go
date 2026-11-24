package main

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestCollectTitlesFromNestedGlob(t *testing.T) {
	root := t.TempDir()
	nested := filepath.Join(root, "nested")
	if err := os.Mkdir(nested, 0o755); err != nil {
		t.Fatal(err)
	}

	write := func(path, content string) {
		t.Helper()
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	write(filepath.Join(nested, "asyncapi.yaml"), "asyncapi: 3.0.0\ninfo:\n  title: Events API\n")
	write(filepath.Join(nested, "events.json"), `{"asyncapi":"3.0.0","info":{"name":"JSON Events"}}`)
	write(filepath.Join(root, "openapi.yaml"), "openapi: 3.1.0\ninfo:\n  title: Not AsyncAPI\n")

	titles, err := collectTitles(root, globDocumentFinder{}, fileDocumentParser{})
	if err != nil {
		t.Fatalf("collectTitles: %v", err)
	}
	want := []string{"Events API", "JSON Events"}
	if !slices.Equal(titles, want) {
		t.Errorf("titles = %q, want %q", titles, want)
	}
}
