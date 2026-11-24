package main

import (
	"bytes"
	"testing"
)

func TestHTMLElementEscapesTextAndAttributes(t *testing.T) {
	var buf bytes.Buffer
	node := htmlElement{
		tag:      "p",
		attrs:    []htmlAttribute{{name: "title", value: `A & "B"`}},
		children: []htmlNode{htmlText("<safe>")},
	}
	if err := node.writeHTML(&buf); err != nil {
		t.Fatal(err)
	}

	want := `<p title="A &amp; &#34;B&#34;">&lt;safe&gt;</p>`
	if got := buf.String(); got != want {
		t.Errorf("HTML = %q, want %q", got, want)
	}
}
