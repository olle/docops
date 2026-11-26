package tui

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"charm.land/bubbletea/v2"

	"github.com/olle/docops/internal/view"
)

type fakeSource struct {
	snapshot view.Snapshot
	err      error
}

func (f fakeSource) Snapshot(context.Context) (view.Snapshot, error) {
	return f.snapshot, f.err
}

func testSnapshot() view.Snapshot {
	return view.Snapshot{
		Greeting: "Hello World!",
		Documents: []view.Document{
			{Source: "events.yaml", Kind: "asyncapi", KindLabel: "AsyncAPI", DisplayName: "Events API", Description: "Event stream."},
			{Source: "api.yaml", Kind: "openapi", KindLabel: "OpenAPI", DisplayName: "HTTP API"},
		},
	}
}

func TestModelRendersSnapshot(t *testing.T) {
	m := newModel(context.Background(), fakeSource{snapshot: testSnapshot()}, testSnapshot(), nil)
	if got := len(m.list.VisibleItems()); got != 2 {
		t.Fatalf("visible items = %d, want 2", got)
	}
	content := m.View().Content
	for _, want := range []string{"Hello World!", "AsyncAPI: Events API", "Specifications"} {
		if !strings.Contains(content, want) {
			t.Errorf("view = %q, want %q", content, want)
		}
	}
	if !strings.Contains(m.detail.View(), "Source:      events.yaml") {
		t.Errorf("detail = %q, want selected document source", m.detail.View())
	}
}

func TestModelQuitKey(t *testing.T) {
	m := newModel(context.Background(), fakeSource{snapshot: testSnapshot()}, testSnapshot(), nil)
	updated, cmd := m.Update(tea.KeyPressMsg{Code: 'q', Text: "q"})
	if updated != m {
		t.Error("expected same model after quit")
	}
	if cmd == nil {
		t.Fatal("expected quit command")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Errorf("quit message = %T, want tea.QuitMsg", cmd())
	}
}

func TestModelReload(t *testing.T) {
	reloaded := view.Snapshot{Greeting: "Hello Reloaded!"}
	m := newModel(context.Background(), fakeSource{snapshot: reloaded}, testSnapshot(), nil)
	_, cmd := m.Update(tea.KeyPressMsg{Code: 'r', Text: "r"})
	if cmd == nil {
		t.Fatal("expected reload command")
	}
	msg, ok := cmd().(loadedMsg)
	if !ok {
		t.Fatalf("reload message = %T, want loadedMsg", cmd())
	}
	if msg.snapshot.Greeting != reloaded.Greeting {
		t.Errorf("reloaded greeting = %q, want %q", msg.snapshot.Greeting, reloaded.Greeting)
	}
}

func TestModelFiltersByKind(t *testing.T) {
	m := newModel(context.Background(), fakeSource{snapshot: testSnapshot()}, testSnapshot(), nil)
	_, _ = m.Update(tea.KeyPressMsg{Code: 'f', Text: "f"})
	if m.kindFilter != "asyncapi" {
		t.Errorf("kind filter = %q, want %q", m.kindFilter, "asyncapi")
	}
	if got := len(m.list.VisibleItems()); got != 1 {
		t.Errorf("visible items = %d, want 1", got)
	}
}

func TestFallbackRendersPlainOutput(t *testing.T) {
	var out bytes.Buffer
	screen := New(fakeSource{snapshot: testSnapshot()}, Options{Out: &out, ForceFallback: true})
	if err := screen.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	got := out.String()
	for _, want := range []string{"Hello World!", "Specifications:", "  * AsyncAPI: Events API", "q: quit"} {
		if !strings.Contains(got, want) {
			t.Errorf("fallback = %q, want %q", got, want)
		}
	}
	if strings.Contains(got, "\x1b") {
		t.Errorf("fallback = %q, want no ANSI escapes", got)
	}
}

func TestProgramQuitsOnQ(t *testing.T) {
	var out bytes.Buffer
	screen := New(fakeSource{snapshot: testSnapshot()}, Options{
		In:               strings.NewReader("q"),
		Out:              &out,
		Width:            100,
		Height:           30,
		ForceInteractive: true,
	})
	if err := screen.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "Hello World!") {
		t.Errorf("program output = %q, want greeting", out.String())
	}
}
