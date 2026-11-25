package main

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestCollectDocumentsFromNestedGlob(t *testing.T) {
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

	write(filepath.Join(nested, "asyncapi.yaml"), "asyncapi: 3.0.0\ninfo:\n  title: Events API\nchannels:\n  greetings:\n    address: greetings\noperations:\n  publishGreeting:\n    action: send\n")
	write(filepath.Join(nested, "events.json"), `{"asyncapi":"3.0.0","info":{"name":"JSON Events"}}`)
	write(filepath.Join(root, "openapi.yaml"), "openapi: 3.1.0\ninfo:\n  title: Not AsyncAPI\npaths:\n  /hello:\n    get: {}\n")

	documents, err := collectDocuments(root, globDocumentFinder{}, fileDocumentParser{})
	if err != nil {
		t.Fatalf("collectDocuments: %v", err)
	}

	wantTitles := []string{"Events API", "JSON Events"}
	if got := documents.Titles(kindAsyncAPI); !slices.Equal(got, wantTitles) {
		t.Errorf("AsyncAPI titles = %q, want %q", got, wantTitles)
	}

	openAPIs := documents.ByKind(kindOpenAPI)
	if len(openAPIs) != 1 {
		t.Fatalf("OpenAPI documents = %d, want 1", len(openAPIs))
	}
	if openAPIs[0].title != "Not AsyncAPI" {
		t.Errorf("OpenAPI title = %q, want %q", openAPIs[0].title, "Not AsyncAPI")
	}
	if !slices.Equal(openAPIs[0].structure.paths, []string{"/hello"}) {
		t.Errorf("OpenAPI paths = %q, want %q", openAPIs[0].structure.paths, []string{"/hello"})
	}

	counts := documents.Counts()
	if counts[kindAsyncAPI] != 2 || counts[kindOpenAPI] != 1 {
		t.Errorf("counts = %v, want asyncapi=2 openapi=1", counts)
	}
}
