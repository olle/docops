package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/olle/docops/internal/assets"
)

func TestHTMLWriterRender(t *testing.T) {
	var buf bytes.Buffer
	msg := newMessage("Ada & Bob", []specificationSummary{{kind: kindAsyncAPI, title: "<Events>"}})
	if err := (htmlWriter{}).Render(&buf, msg); err != nil {
		t.Fatal(err)
	}

	want := "<!doctype html>\n" +
		`<html lang="en"><head><meta charset="utf-8"><title>Hello Ada &amp; Bob!</title><style>` +
		assets.CSS +
		`</style></head>` +
		`<body><main><h1>Hello Ada &amp; Bob!</h1><h2>Specifications</h2><ul><li><strong>AsyncAPI</strong>: &lt;Events&gt;</li></ul></main></body></html>` + "\n"
	if got := buf.String(); got != want {
		t.Errorf("HTML = %q, want %q", got, want)
	}
}

func TestHTMLWriterRenderWithoutSpecifications(t *testing.T) {
	var buf bytes.Buffer
	if err := (htmlWriter{}).Render(&buf, newMessage("World", nil)); err != nil {
		t.Fatal(err)
	}

	got := buf.String()
	if !strings.Contains(got, "<h1>Hello World!</h1>") {
		t.Errorf("HTML = %q, want H1 greeting", got)
	}
	if strings.Contains(got, "<h2>") || strings.Contains(got, "<ul>") {
		t.Errorf("HTML = %q, want no list without specifications", got)
	}
	if !strings.Contains(got, "<style>") || !strings.Contains(got, ":not(pre)>code") {
		t.Errorf("HTML = %q, want inlined, unescaped CSS", got)
	}
	if strings.Contains(got, `<link rel="stylesheet"`) {
		t.Errorf("HTML = %q, want no external stylesheet link", got)
	}
}
