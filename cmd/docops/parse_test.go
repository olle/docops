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

func TestParseJSONSchemaDocument(t *testing.T) {
	data := []byte(`{"$schema":"https://json-schema.org/draft/2020-12/schema","$id":"https://example.com/hello-world.schema.json","title":"Hello World Schema","description":"A minimal JSON Schema example.","properties":{"message":{"type":"string"}}}`)
	spec, err := parseJSONDocument("schema.json", data)
	if err != nil {
		t.Fatal(err)
	}
	if spec.kind != kindJSONSchema {
		t.Errorf("kind = %q, want %q", spec.kind, kindJSONSchema)
	}
	if spec.specVersion != "2020-12" {
		t.Errorf("spec version = %q, want %q", spec.specVersion, "2020-12")
	}
	if spec.id != "https://example.com/hello-world.schema.json" {
		t.Errorf("id = %q, want JSON Schema id", spec.id)
	}
	if spec.title != "Hello World Schema" || spec.description != "A minimal JSON Schema example." {
		t.Errorf("spec = %+v, want JSON Schema metadata", spec)
	}
	if !slices.Equal(spec.structure.properties, []string{"message"}) {
		t.Errorf("properties = %q, want %q", spec.structure.properties, []string{"message"})
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
	spec := parseYAMLDocument("schema.yaml", []byte(source))
	if spec.kind != kindJSONSchema || spec.specVersion != "2020-12" {
		t.Errorf("spec = %+v, want JSON Schema 2020-12", spec)
	}
	if spec.title != "YAML Schema" {
		t.Errorf("title = %q, want %q", spec.title, "YAML Schema")
	}
	if !slices.Equal(spec.structure.properties, []string{"message"}) {
		t.Errorf("properties = %q, want %q", spec.structure.properties, []string{"message"})
	}
}

func TestParseAvroDocument(t *testing.T) {
	data := []byte(`{"type":"record","name":"Greeting","namespace":"dev.docops.hello","doc":"A greeting.","fields":[{"name":"message","type":"string"},{"name":"id","type":"long"}]}`)
	spec, err := parseJSONDocument("greeting.avsc", data)
	if err != nil {
		t.Fatal(err)
	}
	if spec.kind != kindAvro {
		t.Errorf("kind = %q, want %q", spec.kind, kindAvro)
	}
	if spec.id != "dev.docops.hello.Greeting" {
		t.Errorf("id = %q, want qualified Avro name", spec.id)
	}
	if spec.title != "Greeting" || spec.description != "A greeting." {
		t.Errorf("spec = %+v, want Avro metadata", spec)
	}
	if !slices.Equal(spec.structure.fields, []string{"id", "message"}) {
		t.Errorf("fields = %q, want %q", spec.structure.fields, []string{"id", "message"})
	}
}
