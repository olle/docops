package text

import (
	"bytes"
	"testing"

	"github.com/olle/docops/internal/view"
)

func TestPlainTextRender(t *testing.T) {
	var buf bytes.Buffer
	snapshot := view.Snapshot{
		Greeting: "Hello Ada!",
		Documents: []view.Document{
			{Kind: "asyncapi", KindLabel: "AsyncAPI", DisplayName: "Events API"},
		},
	}
	if err := (PlainText{}).Render(&buf, snapshot); err != nil {
		t.Fatal(err)
	}
	want := "Hello Ada!\n\nSpecifications:\n- AsyncAPI: Events API\n"
	if got := buf.String(); got != want {
		t.Errorf("plain text = %q, want %q", got, want)
	}
}

func TestPlainTextRenderWithoutSpecifications(t *testing.T) {
	var buf bytes.Buffer
	if err := (PlainText{}).Render(&buf, view.Snapshot{Greeting: "Hello World!"}); err != nil {
		t.Fatal(err)
	}
	if got, want := buf.String(), "Hello World!\n"; got != want {
		t.Errorf("plain text = %q, want %q", got, want)
	}
}

func TestFallbackRendersNoANSI(t *testing.T) {
	var buf bytes.Buffer
	snapshot := view.Snapshot{
		Greeting: "Hello Ada!",
		Documents: []view.Document{
			{Kind: "asyncapi", KindLabel: "AsyncAPI", DisplayName: "Events API"},
		},
	}
	if err := Fallback(&buf, snapshot); err != nil {
		t.Fatal(err)
	}
	want := "Hello Ada!\n\nSpecifications:\n  * AsyncAPI: Events API\n\nq: quit\n"
	if got := buf.String(); got != want {
		t.Errorf("fallback = %q, want %q", got, want)
	}
}

func TestDebugRender(t *testing.T) {
	var buf bytes.Buffer
	snapshot := view.Snapshot{Greeting: "Hello World!"}
	if err := (Debug{}).Render(&buf, snapshot); err != nil {
		t.Fatal(err)
	}
	if got, want := buf.String(), "catalog:\n  count: 0\n"; got != want {
		t.Errorf("debug = %q, want %q", got, want)
	}
}
