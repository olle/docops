package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestHTMLWriterRender(t *testing.T) {
	var buf bytes.Buffer
	msg := newMessage("Ada & Bob", []string{"<Events>"})
	if err := (htmlWriter{}).Render(&buf, msg); err != nil {
		t.Fatal(err)
	}

	want := "<!doctype html>\n" +
		`<html lang="en"><head><meta charset="utf-8"><title>Hello Ada &amp; Bob!</title></head>` +
		`<body><main><h1>Hello Ada &amp; Bob!</h1><ul><li>&lt;Events&gt;</li></ul></main></body></html>` + "\n"
	if got := buf.String(); got != want {
		t.Errorf("HTML = %q, want %q", got, want)
	}
}

func TestHTMLWriterRenderWithoutTitles(t *testing.T) {
	var buf bytes.Buffer
	if err := (htmlWriter{}).Render(&buf, newMessage("World", nil)); err != nil {
		t.Fatal(err)
	}

	got := buf.String()
	if !strings.Contains(got, "<h1>Hello World!</h1>") {
		t.Errorf("HTML = %q, want H1 greeting", got)
	}
	if strings.Contains(got, "<ul>") {
		t.Errorf("HTML = %q, want no list without titles", got)
	}
}
