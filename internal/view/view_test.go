package view

import (
	"errors"
	"slices"
	"testing"

	"github.com/olle/docops/internal/docops"
)

func TestNewSnapshotProjectsCatalogInDiscoveryOrder(t *testing.T) {
	catalog := docops.NewCatalog([]docops.Specification{
		{Source: "events.yaml", Kind: docops.AsyncAPIKind, Title: "Events API", Name: "Events Name", SpecVersion: "3.0.0"},
		{Source: "api.yaml", Kind: docops.OpenAPIKind, Name: "HTTP API"},
		{Source: "ignored.yaml", Kind: docops.UnknownKind, Title: "Ignored"},
	})

	snapshot := NewSnapshot("Ada", catalog, nil)
	if snapshot.Greeting != "Hello Ada!" {
		t.Errorf("greeting = %q, want %q", snapshot.Greeting, "Hello Ada!")
	}
	if len(snapshot.Documents) != 3 {
		t.Fatalf("documents = %d, want 3", len(snapshot.Documents))
	}
	if got, want := snapshot.Documents[0].Source, "events.yaml"; got != want {
		t.Errorf("first source = %q, want %q", got, want)
	}
	if got, want := snapshot.Documents[0].KindLabel, "AsyncAPI"; got != want {
		t.Errorf("kind label = %q, want %q", got, want)
	}
	if got, want := snapshot.Documents[1].DisplayName, "HTTP API"; got != want {
		t.Errorf("display name = %q, want %q", got, want)
	}
}

func TestSummariesFilterAndSort(t *testing.T) {
	snapshot := Snapshot{Documents: []Document{
		{Kind: "openapi", KindLabel: "OpenAPI", DisplayName: "HTTP API"},
		{Kind: "asyncapi", KindLabel: "AsyncAPI", DisplayName: "Events"},
		{Kind: "avro", KindLabel: "Avro", DisplayName: "Greeting"},
		{Kind: "jsonschema", KindLabel: "JSON Schema", DisplayName: "Schema"},
		{Kind: "openapi", KindLabel: "OpenAPI"},
		{Kind: "", KindLabel: "Unknown", DisplayName: "Ignored"},
	}}

	summaries := Summaries(snapshot)
	got := make([]string, 0, len(summaries))
	for _, doc := range summaries {
		got = append(got, doc.KindLabel+": "+doc.DisplayName)
	}
	want := []string{
		"AsyncAPI: Events",
		"Avro: Greeting",
		"JSON Schema: Schema",
		"OpenAPI: HTTP API",
	}
	if !slices.Equal(got, want) {
		t.Errorf("summaries = %q, want %q", got, want)
	}
}

func TestDebugCascadesFullModel(t *testing.T) {
	snapshot := Snapshot{Documents: []Document{{
		Source:      "events.yaml",
		KindLabel:   "AsyncAPI",
		SpecVersion: "3.0.0",
		Version:     "1.2.3",
		Title:       "Events API",
		Description: "Event stream.",
		Structure: Structure{
			Channels:   []string{"farewells", "greetings"},
			Operations: []string{"publishGreeting"},
		},
	}}}

	want := `catalog:
  count: 1
  documents:
    - specification:
      source: "events.yaml"
      kind: AsyncAPI
      specVersion: "3.0.0"
      id: ""
      version: "1.2.3"
      title: "Events API"
      name: ""
      description: "Event stream."
      structure:
        channels: ["farewells", "greetings"]
        operations: ["publishGreeting"]
        paths: []
        properties: []
        fields: []
`
	if got := Debug(snapshot); got != want {
		t.Errorf("Debug() = %q, want %q", got, want)
	}
}

func TestDebugEmpty(t *testing.T) {
	if got, want := Debug(Snapshot{}), "catalog:\n  count: 0\n"; got != want {
		t.Errorf("Debug() = %q, want %q", got, want)
	}
}

func TestNewSnapshotKeepsWarnings(t *testing.T) {
	err := errors.Join(errors.New("first"), errors.New("second"))
	snapshot := NewSnapshot("Ada", docops.Catalog{}, err)
	if !slices.Equal(snapshot.Warnings, []string{"first", "second"}) {
		t.Errorf("warnings = %q, want first and second", snapshot.Warnings)
	}
}
