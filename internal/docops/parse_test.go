package docops

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
			spec := ParseYAMLDocument("example.yaml", []byte(test.source))
			if spec.Kind != AsyncAPIKind {
				t.Errorf("kind = %q, want %q", spec.Kind, AsyncAPIKind)
			}
			if spec.SpecVersion == "" {
				t.Error("spec version not parsed")
			}
			if spec.Title != test.wantTitle {
				t.Errorf("title = %q, want %q", spec.Title, test.wantTitle)
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
	spec := ParseYAMLDocument("events.yaml", []byte(source))
	if spec.Kind != AsyncAPIKind {
		t.Errorf("kind = %q, want %q", spec.Kind, AsyncAPIKind)
	}
	if spec.SpecVersion != "3.0.0" {
		t.Errorf("spec version = %q, want %q", spec.SpecVersion, "3.0.0")
	}
	if spec.Version != "1.2.3" {
		t.Errorf("version = %q, want %q", spec.Version, "1.2.3")
	}
	if spec.Description != "Event stream." {
		t.Errorf("description = %q, want %q", spec.Description, "Event stream.")
	}
	if !slices.Equal(spec.Structure.Channels, []string{"farewells", "greetings"}) {
		t.Errorf("channels = %q, want sorted channel names", spec.Structure.Channels)
	}
	if !slices.Equal(spec.Structure.Operations, []string{"publishGreeting"}) {
		t.Errorf("operations = %q, want sorted operation names", spec.Structure.Operations)
	}
}

func TestParseJSONDocument(t *testing.T) {
	data := []byte(`{"openapi":"3.1.0","info":{"title":"API","description":"Desc","version":"1.0.0"},"paths":{"/hello":{"get":{}}}}`)
	spec, err := ParseJSONDocument("api.json", data)
	if err != nil {
		t.Fatal(err)
	}
	if spec.Kind != OpenAPIKind {
		t.Errorf("kind = %q, want %q", spec.Kind, OpenAPIKind)
	}
	if spec.SpecVersion != "3.1.0" {
		t.Errorf("spec version = %q, want %q", spec.SpecVersion, "3.1.0")
	}
	if spec.Title != "API" || spec.Description != "Desc" || spec.Version != "1.0.0" {
		t.Errorf("spec = %+v, want parsed info metadata", spec)
	}
	if !slices.Equal(spec.Structure.Paths, []string{"/hello"}) {
		t.Errorf("paths = %q, want %q", spec.Structure.Paths, []string{"/hello"})
	}
}

func TestParseJSONSchemaDocument(t *testing.T) {
	data := []byte(`{"$schema":"https://json-schema.org/draft/2020-12/schema","$id":"https://example.com/hello-world.schema.json","title":"Hello World Schema","description":"A minimal JSON Schema example.","properties":{"message":{"type":"string"}}}`)
	spec, err := ParseJSONDocument("schema.json", data)
	if err != nil {
		t.Fatal(err)
	}
	if spec.Kind != JSONSchemaKind {
		t.Errorf("kind = %q, want %q", spec.Kind, JSONSchemaKind)
	}
	if spec.SpecVersion != "2020-12" {
		t.Errorf("spec version = %q, want %q", spec.SpecVersion, "2020-12")
	}
	if spec.ID != "https://example.com/hello-world.schema.json" {
		t.Errorf("id = %q, want JSON Schema id", spec.ID)
	}
	if spec.Title != "Hello World Schema" || spec.Description != "A minimal JSON Schema example." {
		t.Errorf("spec = %+v, want JSON Schema metadata", spec)
	}
	if !slices.Equal(spec.Structure.Properties, []string{"message"}) {
		t.Errorf("properties = %q, want %q", spec.Structure.Properties, []string{"message"})
	}
}

func TestParseYAMLJSONSchemaDocument(t *testing.T) {
	source := `$schema: https://json-schema.org/draft/2020-12/schema
$id: https://example.com/hello-world.schema.yaml
title: YAML Schema
properties:
  message:
    type: string
`
	spec := ParseYAMLDocument("schema.yaml", []byte(source))
	if spec.Kind != JSONSchemaKind || spec.SpecVersion != "2020-12" {
		t.Errorf("spec = %+v, want JSON Schema 2020-12", spec)
	}
	if spec.Title != "YAML Schema" {
		t.Errorf("title = %q, want %q", spec.Title, "YAML Schema")
	}
	if !slices.Equal(spec.Structure.Properties, []string{"message"}) {
		t.Errorf("properties = %q, want %q", spec.Structure.Properties, []string{"message"})
	}
}

func TestParseAvroDocument(t *testing.T) {
	data := []byte(`{"type":"record","name":"Greeting","namespace":"dev.docops.hello","doc":"A greeting.","fields":[{"name":"message","type":"string"},{"name":"id","type":"long"}]}`)
	spec, err := ParseJSONDocument("greeting.avsc", data)
	if err != nil {
		t.Fatal(err)
	}
	if spec.Kind != AvroKind {
		t.Errorf("kind = %q, want %q", spec.Kind, AvroKind)
	}
	if spec.ID != "dev.docops.hello.Greeting" {
		t.Errorf("id = %q, want qualified Avro name", spec.ID)
	}
	if spec.Title != "Greeting" || spec.Description != "A greeting." {
		t.Errorf("spec = %+v, want Avro metadata", spec)
	}
	if !slices.Equal(spec.Structure.Fields, []string{"id", "message"}) {
		t.Errorf("fields = %q, want %q", spec.Structure.Fields, []string{"id", "message"})
	}
}
