package main

import (
	"bytes"
	"strings"
	"testing"
)

type fakeTerminal struct {
	styled   bool
	restored bool
}

func (t *fakeTerminal) Raw() (func(), bool) {
	t.restored = false
	return func() { t.restored = true }, t.styled
}

func TestRenderTUIPlain(t *testing.T) {
	var buf bytes.Buffer
	if err := renderTUI(&buf, newMessage("Ada", []string{"Events API"}), false); err != nil {
		t.Fatal(err)
	}

	got := buf.String()
	for _, want := range []string{"Hello Ada!", "AsyncAPI documents:", "  * Events API", "q: quit"} {
		if !strings.Contains(got, want) {
			t.Errorf("TUI output = %q, want %q", got, want)
		}
	}
	if strings.Contains(got, "\x1b") {
		t.Errorf("TUI output = %q, want no ANSI escapes", got)
	}
}

func TestRenderTUIStyled(t *testing.T) {
	var buf bytes.Buffer
	if err := renderTUI(&buf, newMessage("World", nil), true); err != nil {
		t.Fatal(err)
	}

	got := buf.String()
	if !strings.Contains(got, "\x1b[?1049h") || !strings.Contains(got, "\x1b[1mHello World!\x1b[0m") {
		t.Errorf("TUI output = %q, want styled screen", got)
	}
}

func TestTUIScreenRestoresTerminal(t *testing.T) {
	term := &fakeTerminal{styled: true}
	var out bytes.Buffer
	screen := tuiScreen{out: &out, in: strings.NewReader("q"), term: term}
	if err := screen.Present(newMessage("World", nil)); err != nil {
		t.Fatal(err)
	}
	if !term.restored {
		t.Error("terminal was not restored")
	}
	if !strings.Contains(out.String(), "\x1b[?25h\x1b[?1049l") {
		t.Errorf("TUI output = %q, want alternate screen exit", out.String())
	}
}

func TestWaitForQuit(t *testing.T) {
	if err := waitForQuit(strings.NewReader("xq"), true); err != nil {
		t.Errorf("raw quit: %v", err)
	}
	if err := waitForQuit(strings.NewReader("hello\nq\n"), false); err != nil {
		t.Errorf("line quit: %v", err)
	}
	if err := waitForQuit(strings.NewReader(""), false); err != nil {
		t.Errorf("EOF quit: %v", err)
	}
}

func TestParseConfigRejectsServerAndTUI(t *testing.T) {
	if _, err := parseConfig([]string{"-server", "-tui"}); err == nil {
		t.Error("expected -server and -tui to be mutually exclusive")
	}
}
