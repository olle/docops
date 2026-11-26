package text

import (
	"bytes"
	"testing"

	"github.com/olle/docops/internal/view"
)

func TestMarkdownRender(t *testing.T) {
	var buf bytes.Buffer
	snapshot := view.Snapshot{
		Greeting: "Ada & Bob",
		Documents: []view.Document{
			{Kind: "asyncapi", KindLabel: "AsyncAPI", DisplayName: "<Events> *beta*"},
			{Kind: "openapi", KindLabel: "OpenAPI", DisplayName: "Plain"},
		},
	}
	if err := (Markdown{}).Render(&buf, snapshot); err != nil {
		t.Fatal(err)
	}

	want := `# Ada & Bob

## Specifications

- AsyncAPI: \<Events\> \*beta\*
- OpenAPI: Plain
`
	if got := buf.String(); got != want {
		t.Errorf("Markdown = %q, want %q", got, want)
	}
}

func TestMarkdownRenderWithoutSpecifications(t *testing.T) {
	var buf bytes.Buffer
	if err := (Markdown{}).Render(&buf, view.Snapshot{Greeting: "Hello World!"}); err != nil {
		t.Fatal(err)
	}
	if got, want := buf.String(), "# Hello World!\n"; got != want {
		t.Errorf("Markdown = %q, want %q", got, want)
	}
}
