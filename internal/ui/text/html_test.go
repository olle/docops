package text

import (
	"bytes"
	"strings"
	"testing"

	"github.com/olle/docops/internal/assets"
	"github.com/olle/docops/internal/view"
)

func TestHTMLRender(t *testing.T) {
	var buf bytes.Buffer
	snapshot := view.Snapshot{
		Greeting: "Ada & Bob",
		Documents: []view.Document{
			{Kind: "asyncapi", KindLabel: "AsyncAPI", DisplayName: "<Events>"},
		},
	}
	if err := (HTML{}).Render(&buf, snapshot); err != nil {
		t.Fatal(err)
	}

	want := "<!doctype html>\n" +
		`<html lang="en"><head><meta charset="utf-8"><title>Ada &amp; Bob</title><style>` +
		assets.CSS +
		`</style></head>` +
		`<body><main><h1>Ada &amp; Bob</h1><h2>Specifications</h2><ul><li><strong>AsyncAPI</strong>: &lt;Events&gt;</li></ul></main></body></html>` + "\n"
	if got := buf.String(); got != want {
		t.Errorf("HTML = %q, want %q", got, want)
	}
}

func TestHTMLRenderWithoutSpecifications(t *testing.T) {
	var buf bytes.Buffer
	if err := (HTML{}).Render(&buf, view.Snapshot{Greeting: "Hello World!"}); err != nil {
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
