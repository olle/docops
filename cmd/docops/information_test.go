package main

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestDocumentInformationBuildsMessage(t *testing.T) {
	root := t.TempDir()
	asyncPath := filepath.Join(root, "asyncapi.yaml")
	if err := os.WriteFile(asyncPath, []byte("asyncapi: 3.0.0\ninfo:\n  title: Events API\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	openPath := filepath.Join(root, "openapi.json")
	if err := os.WriteFile(openPath, []byte(`{"openapi":"3.1.0","info":{"title":"Hello World API"}}`), 0o600); err != nil {
		t.Fatal(err)
	}

	source := documentInformation{
		name:    "Ada",
		root:    root,
		finder:  documentGlobFinder{},
		sniffer: fileDocumentSniffer{},
		parser:  fileDocumentParser{},
	}
	msg, err := source.Message()
	if err != nil {
		t.Fatal(err)
	}
	if msg.greeting != "Hello Ada!" {
		t.Errorf("greeting = %q, want %q", msg.greeting, "Hello Ada!")
	}
	want := []specificationSummary{
		{kind: kindAsyncAPI, title: "Events API"},
		{kind: kindOpenAPI, title: "Hello World API"},
	}
	if !slices.Equal(msg.specifications, want) {
		t.Errorf("specifications = %+v, want %+v", msg.specifications, want)
	}
}

func TestFixedInformationReturnsMessage(t *testing.T) {
	want := newMessage("Ada", []specificationSummary{{kind: kindOpenAPI, title: "Events API"}})
	msg, err := fixedInformation{msg: want}.Message()
	if err != nil {
		t.Fatal(err)
	}
	if msg.greeting != want.greeting || !slices.Equal(msg.specifications, want.specifications) {
		t.Errorf("message = %+v, want %+v", msg, want)
	}
}
