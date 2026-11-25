package main

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
		want   specificationKind
	}{
		{name: "asyncapi", prefix: "asyncapi: 3.0.0\n", want: kindAsyncAPI},
		{name: "openapi", prefix: "openapi: 3.1.0\n", want: kindOpenAPI},
		{name: "nested marker", prefix: "description:\n  asyncapi: nested\n", want: kindUnknown},
		{name: "unrelated yaml", prefix: "name: example\n", want: kindUnknown},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := detectSpecificationKind(test.prefix); got != test.want {
				t.Errorf("kind = %q, want %q", got, test.want)
			}
		})
	}
}

func TestDetectSpecificationKindJSON(t *testing.T) {
	if got := detectSpecificationKind(`{"asyncapi":"3.0.0","info":{}}`); got != kindAsyncAPI {
		t.Errorf("kind = %q, want %q", got, kindAsyncAPI)
	}
	if got := detectSpecificationKind(`{"openapi":"3.1.0","info":{}}`); got != kindOpenAPI {
		t.Errorf("kind = %q, want %q", got, kindOpenAPI)
	}
}

func TestFileDocumentSnifferReadsBoundedPrefix(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "late.yaml")
	content := strings.Repeat("x", sniffLimit) + "\nasyncapi: 3.0.0\n"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	kind, err := (fileDocumentSniffer{}).Detect(path)
	if err != nil {
		t.Fatal(err)
	}
	if kind != kindUnknown {
		t.Errorf("kind = %q, want %q for marker beyond sniff limit", kind, kindUnknown)
	}
}
