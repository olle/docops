package pdf

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/olle/docops/internal/view"
)

func TestPresenterWritesPDFEnvelope(t *testing.T) {
	var out bytes.Buffer
	snapshot := view.Snapshot{
		Greeting: "Hello Ada!",
		Documents: []view.Document{{
			Source:      "events.yaml",
			Kind:        "asyncapi",
			KindLabel:   "AsyncAPI",
			DisplayName: "Events API",
			SpecVersion: "3.0.0",
			Description: "Event stream.",
			Structure:   view.Structure{Channels: []string{"greetings"}},
		}},
	}
	if err := (Presenter{}).Render(&out, snapshot); err != nil {
		t.Fatal(err)
	}

	got := out.String()
	for _, want := range []string{
		"%PDF-1.4\n",
		"/Type /Catalog",
		"/Count 1",
		"/BaseFont /Helvetica",
		"/BaseFont /Helvetica-Bold",
		"/Title (Hello Ada!)",
		"(Hello Ada!) Tj",
		"(AsyncAPI: Events API) Tj",
		"(Source: events.yaml) Tj",
		"(Channels: greetings) Tj",
		"trailer",
		"%%EOF\n",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("PDF output missing %q", want)
		}
	}
}

func TestPresenterEscapesText(t *testing.T) {
	var out bytes.Buffer
	if err := (Presenter{}).Render(&out, view.Snapshot{Greeting: `A (B) \ C`}); err != nil {
		t.Fatal(err)
	}
	if got, want := out.String(), `(A \(B\) \\ C) Tj`; !strings.Contains(got, want) {
		t.Errorf("PDF output = %q, want %q", got, want)
	}
}

func TestPresenterReplacesUnsupportedRunes(t *testing.T) {
	var out bytes.Buffer
	if err := (Presenter{}).Render(&out, view.Snapshot{Greeting: "Hello 😀"}); err != nil {
		t.Fatal(err)
	}
	if got, want := out.String(), `(Hello ?) Tj`; !strings.Contains(got, want) {
		t.Errorf("PDF output = %q, want %q", got, want)
	}
}

func TestPresenterPaginates(t *testing.T) {
	var out bytes.Buffer
	documents := make([]view.Document, 80)
	for i := range documents {
		documents[i] = view.Document{
			Source:      "spec-" + strconv.Itoa(i) + ".yaml",
			Kind:        "openapi",
			KindLabel:   "OpenAPI",
			DisplayName: "API " + strconv.Itoa(i),
			SpecVersion: "3.1.0",
			Description: strings.Repeat("description ", 12),
		}
	}
	if err := (Presenter{}).Render(&out, view.Snapshot{Greeting: "Hello World!", Documents: documents}); err != nil {
		t.Fatal(err)
	}

	count := pageCount(t, out.String())
	if count < 2 {
		t.Errorf("page count = %d, want more than 1", count)
	}
}

func TestPresenterWritesValidXrefOffsets(t *testing.T) {
	var out bytes.Buffer
	snapshot := view.Snapshot{
		Greeting: "Hello Ada!",
		Documents: []view.Document{{
			Source:      "api.yaml",
			Kind:        "openapi",
			KindLabel:   "OpenAPI",
			DisplayName: "HTTP API",
		}},
	}
	if err := (Presenter{}).Render(&out, snapshot); err != nil {
		t.Fatal(err)
	}

	data := out.Bytes()
	text := string(data)
	marker := "startxref\n"
	index := strings.LastIndex(text, marker)
	if index < 0 {
		t.Fatal("PDF output missing startxref")
	}
	offsetLine := strings.SplitN(text[index+len(marker):], "\n", 2)[0]
	xrefOffset, err := strconv.Atoi(strings.TrimSpace(offsetLine))
	if err != nil {
		t.Fatalf("parsing startxref: %v", err)
	}
	if !strings.HasPrefix(text[xrefOffset:], "xref\n") {
		t.Fatalf("startxref %d does not point to xref", xrefOffset)
	}

	lines := strings.Split(text[xrefOffset:], "\n")
	var size int
	if _, err := fmt.Sscanf(lines[1], "0 %d", &size); err != nil {
		t.Fatalf("parsing xref size: %v", err)
	}
	for id := 1; id < size; id++ {
		fields := strings.Fields(lines[2+id])
		if len(fields) != 3 || fields[2] != "n" {
			t.Fatalf("xref entry %d = %q, want in-use entry", id, lines[2+id])
		}
		offset, err := strconv.Atoi(fields[0])
		if err != nil {
			t.Fatalf("parsing xref offset for object %d: %v", id, err)
		}
		want := strconv.Itoa(id) + " 0 obj"
		if !strings.HasPrefix(text[offset:], want) {
			t.Errorf("xref offset for object %d points to %q, want %q", id, text[offset:offset+len(want)], want)
		}
	}
}

func TestWrapTextSplitsLongWords(t *testing.T) {
	lines := wrapText("supercalifragilisticexpialidocious end", 12)
	want := []string{"supercalifra", "gilisticexpi", "alidocious", "end"}
	if strings.Join(lines, "|") != strings.Join(want, "|") {
		t.Errorf("lines = %q, want %q", lines, want)
	}
}

func pageCount(t *testing.T, pdf string) int {
	t.Helper()
	const prefix = "/Count "
	index := strings.Index(pdf, prefix)
	if index < 0 {
		t.Fatalf("PDF output missing %q", prefix)
	}
	rest := pdf[index+len(prefix):]
	end := strings.IndexAny(rest, " >")
	if end < 0 {
		t.Fatalf("PDF output has malformed page count")
	}
	count, err := strconv.Atoi(rest[:end])
	if err != nil {
		t.Fatalf("parsing page count: %v", err)
	}
	return count
}
