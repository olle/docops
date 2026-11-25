package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

var (
	_ fmt.Stringer = catalog{}
	_ fmt.Stringer = specification{}
	_ fmt.Stringer = specificationStructure{}
	_ fmt.Stringer = specificationSummary{}
	_ fmt.Stringer = kindAsyncAPI
)

func TestCollectDocumentsFromYAMLAndJSONCandidates(t *testing.T) {
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

	write(filepath.Join(nested, "spec.yaml"), "asyncapi: 3.0.0\ninfo:\n  title: Events API\nchannels:\n  greetings:\n    address: greetings\noperations:\n  publishGreeting:\n    action: send\n")
	write(filepath.Join(nested, "messages.json"), `{"asyncapi":"3.0.0","info":{"name":"JSON Events"}}`)
	write(filepath.Join(nested, "peek.yml"), "openapi: 3.1.0\ninfo:\n  title: Peek API\npaths:\n  /ping:\n    get: {}\n")
	write(filepath.Join(root, "api.yaml"), "openapi: 3.1.0\ninfo:\n  title: HTTP API\npaths:\n  /hello:\n    get: {}\n")
	write(filepath.Join(root, "readme.yaml"), "notes: not a specification\n")

	documents, err := collectDocuments(root, documentGlobFinder{}, fileDocumentSniffer{}, fileDocumentParser{})
	if err != nil {
		t.Fatalf("collectDocuments: %v", err)
	}

	wantTitles := []string{"Events API", "JSON Events"}
	if got := documents.Titles(kindAsyncAPI); !slices.Equal(got, wantTitles) {
		t.Errorf("AsyncAPI titles = %q, want %q", got, wantTitles)
	}

	openAPIs := documents.ByKind(kindOpenAPI)
	if len(openAPIs) != 2 {
		t.Fatalf("OpenAPI documents = %d, want 2", len(openAPIs))
	}
	byTitle := make(map[string]specification, len(openAPIs))
	for _, doc := range openAPIs {
		byTitle[doc.title] = doc
	}
	if doc, ok := byTitle["HTTP API"]; !ok || !slices.Equal(doc.structure.paths, []string{"/hello"}) {
		t.Errorf("HTTP API paths = %q, want %q", doc.structure.paths, []string{"/hello"})
	}
	if doc, ok := byTitle["Peek API"]; !ok || !slices.Equal(doc.structure.paths, []string{"/ping"}) {
		t.Errorf("Peek API paths = %q, want %q", doc.structure.paths, []string{"/ping"})
	}

	counts := documents.Counts()
	if counts[kindAsyncAPI] != 2 || counts[kindOpenAPI] != 2 {
		t.Errorf("counts = %v, want asyncapi=2 openapi=2", counts)
	}
	for _, doc := range documents.All() {
		if strings.Contains(doc.source, "readme.yaml") {
			t.Errorf("non-specification file was included: %s", doc.source)
		}
	}
}

func TestCatalogSummariesUseTitleOrNameAndKind(t *testing.T) {
	documents := catalog{documents: []specification{
		{kind: kindOpenAPI, title: "HTTP API"},
		{kind: kindAsyncAPI, name: "Events"},
		{kind: kindOpenAPI},
		{kind: kindUnknown, title: "Ignored"},
	}}
	want := []specificationSummary{
		{kind: kindAsyncAPI, title: "Events"},
		{kind: kindOpenAPI, title: "HTTP API"},
	}
	if got := documents.Summaries(); !slices.Equal(got, want) {
		t.Errorf("summaries = %+v, want %+v", got, want)
	}
}

func TestCatalogStringCascadesFullModel(t *testing.T) {
	documents := catalog{documents: []specification{{
		source:      "events.yaml",
		kind:        kindAsyncAPI,
		specVersion: "3.0.0",
		version:     "1.2.3",
		title:       "Events API",
		description: "Event stream.",
		structure: specificationStructure{
			channels:   []string{"farewells", "greetings"},
			operations: []string{"publishGreeting"},
		},
	}}}

	want := `catalog:
  count: 1
  documents:
    - specification:
      source: "events.yaml"
      kind: AsyncAPI
      specVersion: "3.0.0"
      version: "1.2.3"
      title: "Events API"
      name: ""
      description: "Event stream."
      structure:
        channels: ["farewells", "greetings"]
        operations: ["publishGreeting"]
        paths: []
`
	if got := documents.String(); got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
	if got := fmt.Sprint(documents); got != want {
		t.Errorf("fmt.Sprint = %q, want %q", got, want)
	}
}

func TestCatalogStringEmpty(t *testing.T) {
	if got, want := (catalog{}).String(), "catalog:\n  count: 0\n"; got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestSpecificationSummaryString(t *testing.T) {
	summary := specificationSummary{kind: kindOpenAPI, title: "HTTP API"}
	if got, want := summary.String(), "OpenAPI: HTTP API"; got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}
