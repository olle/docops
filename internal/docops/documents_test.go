package docops

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
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
	write(filepath.Join(root, "jsonschema.json"), `{"$schema":"https://json-schema.org/draft/2020-12/schema","title":"Schema API"}`)
	write(filepath.Join(root, "greeting.avsc"), `{"type":"record","name":"Greeting","fields":[{"name":"message","type":"string"}]}`)
	write(filepath.Join(root, "readme.yaml"), "notes: not a specification\n")

	catalog, err := Collect(root, GlobFinder{}, FileSniffer{}, FileParser{})
	if err != nil {
		t.Fatalf("Collect: %v", err)
	}

	wantTitles := []string{"Events API", "JSON Events"}
	if got := catalog.Titles(AsyncAPIKind); !slices.Equal(got, wantTitles) {
		t.Errorf("AsyncAPI titles = %q, want %q", got, wantTitles)
	}

	openAPIs := catalog.ByKind(OpenAPIKind)
	if len(openAPIs) != 2 {
		t.Fatalf("OpenAPI documents = %d, want 2", len(openAPIs))
	}
	byTitle := make(map[string]Specification, len(openAPIs))
	for _, doc := range openAPIs {
		byTitle[doc.Title] = doc
	}
	if doc, ok := byTitle["HTTP API"]; !ok || !slices.Equal(doc.Structure.Paths, []string{"/hello"}) {
		t.Errorf("HTTP API paths = %q, want %q", doc.Structure.Paths, []string{"/hello"})
	}
	if doc, ok := byTitle["Peek API"]; !ok || !slices.Equal(doc.Structure.Paths, []string{"/ping"}) {
		t.Errorf("Peek API paths = %q, want %q", doc.Structure.Paths, []string{"/ping"})
	}

	counts := catalog.Counts()
	if counts[AsyncAPIKind] != 2 || counts[OpenAPIKind] != 2 || counts[JSONSchemaKind] != 1 || counts[AvroKind] != 1 {
		t.Errorf("counts = %v, want asyncapi=2 openapi=2 jsonschema=1 avro=1", counts)
	}
	for _, doc := range catalog.All() {
		if strings.Contains(doc.Source, "readme.yaml") {
			t.Errorf("non-specification file was included: %s", doc.Source)
		}
	}
}
