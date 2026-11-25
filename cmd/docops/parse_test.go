package main

import (
	"slices"
	"testing"
)

func TestParseYAMLDocument(t *testing.T) {
	tests := []struct {
		name      string
		source    string
		wantTitle string
	}{
		{
			name:      "block scalar with comment",
			source:    "asyncapi: 3.0.0\ninfo:\n  title: Block Title # comment\n",
			wantTitle: "Block Title",
		},
		{
			name:      "quoted scalar",
			source:    "asyncapi: '3.0.0'\ninfo:\n  title: \"Quoted Title\"\n",
			wantTitle: "Quoted Title",
		},
		{
			name:      "flow mapping",
			source:    "asyncapi: 3.0.0\ninfo: {title: Flow Title, version: 1.0.0}\n",
			wantTitle: "Flow Title",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			spec := parseYAMLDocument("example.yaml", []byte(test.source))
			if spec.kind != kindAsyncAPI {
				t.Errorf("kind = %q, want %q", spec.kind, kindAsyncAPI)
			}
			if spec.specVersion == "" {
				t.Error("spec version not parsed")
			}
			if spec.title != test.wantTitle {
				t.Errorf("title = %q, want %q", spec.title, test.wantTitle)
			}
		})
	}
}

func TestParseYAMLDocumentStructure(t *testing.T) {
	source := `asyncapi: 3.0.0
info:
  title: Events API
  version: 1.2.3
  description: Event stream.
channels:
  greetings:
    address: greetings
  farewells:
    address: farewells
operations:
  publishGreeting:
    action: send
`
	spec := parseYAMLDocument("events.yaml", []byte(source))
	if spec.kind != kindAsyncAPI {
		t.Errorf("kind = %q, want %q", spec.kind, kindAsyncAPI)
	}
	if spec.specVersion != "3.0.0" {
		t.Errorf("spec version = %q, want %q", spec.specVersion, "3.0.0")
	}
	if spec.version != "1.2.3" {
		t.Errorf("version = %q, want %q", spec.version, "1.2.3")
	}
	if spec.description != "Event stream." {
		t.Errorf("description = %q, want %q", spec.description, "Event stream.")
	}
	if !slices.Equal(spec.structure.channels, []string{"farewells", "greetings"}) {
		t.Errorf("channels = %q, want sorted channel names", spec.structure.channels)
	}
	if !slices.Equal(spec.structure.operations, []string{"publishGreeting"}) {
		t.Errorf("operations = %q, want sorted operation names", spec.structure.operations)
	}
}

func TestParseJSONDocument(t *testing.T) {
	data := []byte(`{"openapi":"3.1.0","info":{"title":"API","description":"Desc","version":"1.0.0"},"paths":{"/hello":{"get":{}}}}`)
	spec, err := parseJSONDocument("api.json", data)
	if err != nil {
		t.Fatal(err)
	}
	if spec.kind != kindOpenAPI {
		t.Errorf("kind = %q, want %q", spec.kind, kindOpenAPI)
	}
	if spec.specVersion != "3.1.0" {
		t.Errorf("spec version = %q, want %q", spec.specVersion, "3.1.0")
	}
	if spec.title != "API" || spec.description != "Desc" || spec.version != "1.0.0" {
		t.Errorf("spec = %+v, want parsed info metadata", spec)
	}
	if !slices.Equal(spec.structure.paths, []string{"/hello"}) {
		t.Errorf("paths = %q, want %q", spec.structure.paths, []string{"/hello"})
	}
}
