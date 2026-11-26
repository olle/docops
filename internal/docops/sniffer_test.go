package docops

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDetectSpecificationKindYAML(t *testing.T) {
	tests := []struct {
		name   string
		prefix string
		want   Kind
	}{
		{name: "asyncapi", prefix: "asyncapi: 3.0.0\n", want: AsyncAPIKind},
		{name: "openapi", prefix: "openapi: 3.1.0\n", want: OpenAPIKind},
		{name: "json schema", prefix: "$schema: https://json-schema.org/draft/2020-12/schema\n", want: JSONSchemaKind},
		{name: "nested marker", prefix: "description:\n  asyncapi: nested\n", want: UnknownKind},
		{name: "unrelated yaml", prefix: "name: example\n", want: UnknownKind},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := DetectSpecificationKind(test.prefix); got != test.want {
				t.Errorf("kind = %q, want %q", got, test.want)
			}
		})
	}
}

func TestDetectSpecificationKindJSON(t *testing.T) {
	if got := DetectSpecificationKind(`{"asyncapi":"3.0.0","info":{}}`); got != AsyncAPIKind {
		t.Errorf("kind = %q, want %q", got, AsyncAPIKind)
	}
	if got := DetectSpecificationKind(`{"openapi":"3.1.0","info":{}}`); got != OpenAPIKind {
		t.Errorf("kind = %q, want %q", got, OpenAPIKind)
	}
	if got := DetectSpecificationKind(`{"$schema":"https://json-schema.org/draft/2020-12/schema","title":"Example"}`); got != JSONSchemaKind {
		t.Errorf("kind = %q, want %q", got, JSONSchemaKind)
	}
	if got := DetectSpecificationKind(`{"type":"record","name":"Greeting","fields":[]}`); got != AvroKind {
		t.Errorf("kind = %q, want %q", got, AvroKind)
	}
}

func TestFileSnifferReadsBoundedPrefix(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "late.yaml")
	content := strings.Repeat("x", sniffLimit) + "\nasyncapi: 3.0.0\n"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	kind, err := (FileSniffer{}).Detect(path)
	if err != nil {
		t.Fatal(err)
	}
	if kind != UnknownKind {
		t.Errorf("kind = %q, want %q for marker beyond sniff limit", kind, UnknownKind)
	}
}
