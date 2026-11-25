package main

import (
	"bytes"
	"testing"
)

func TestMarkdownWriterRender(t *testing.T) {
	var buf bytes.Buffer
	msg := newMessage("Ada & Bob", []specificationSummary{
		{kind: kindAsyncAPI, title: "<Events> *beta*"},
		{kind: kindOpenAPI, title: "Plain"},
	})
	if err := (markdownWriter{}).Render(&buf, msg); err != nil {
		t.Fatal(err)
	}

	want := `# Hello Ada & Bob!

## Specifications

- AsyncAPI: \<Events\> \*beta\*
- OpenAPI: Plain
`
	if got := buf.String(); got != want {
		t.Errorf("Markdown = %q, want %q", got, want)
	}
}

func TestMarkdownWriterRenderWithoutSpecifications(t *testing.T) {
	var buf bytes.Buffer
	if err := (markdownWriter{}).Render(&buf, newMessage("World", nil)); err != nil {
		t.Fatal(err)
	}
	if got, want := buf.String(), "# Hello World!\n"; got != want {
		t.Errorf("Markdown = %q, want %q", got, want)
	}
}
