package main

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestDocumentInformationBuildsMessage(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "asyncapi.yaml")
	if err := os.WriteFile(path, []byte("asyncapi: 3.0.0\ninfo:\n  title: Events API\n"), 0o600); err != nil {
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
	if !slices.Equal(msg.titles, []string{"Events API"}) {
		t.Errorf("titles = %q, want %q", msg.titles, []string{"Events API"})
	}
}

func TestFixedInformationReturnsMessage(t *testing.T) {
	want := newMessage("Ada", []string{"Events API"})
	msg, err := fixedInformation{msg: want}.Message()
	if err != nil {
		t.Fatal(err)
	}
	if msg.greeting != want.greeting || !slices.Equal(msg.titles, want.titles) {
		t.Errorf("message = %+v, want %+v", msg, want)
	}
}
