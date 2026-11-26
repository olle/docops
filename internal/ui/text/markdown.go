package text

import (
	"fmt"
	"io"
	"strings"

	"github.com/olle/docops/internal/view"
)

type Markdown struct{}

var markdownEscaper = strings.NewReplacer(
	"\\", "\\\\",
	"`", "\\`",
	"*", "\\*",
	"_", "\\_",
	"[", "\\[",
	"]", "\\]",
	"<", "\\<",
	">", "\\>",
)

func (Markdown) ContentType() string {
	return "text/markdown; charset=utf-8"
}

func (Markdown) Render(w io.Writer, snapshot view.Snapshot) error {
	if _, err := fmt.Fprintf(w, "# %s\n", markdownLine(snapshot.Greeting)); err != nil {
		return err
	}
	summaries := view.Summaries(snapshot)
	if len(summaries) == 0 {
		return nil
	}
	if _, err := io.WriteString(w, "\n## Specifications\n\n"); err != nil {
		return err
	}
	for _, doc := range summaries {
		if _, err := fmt.Fprintf(w, "- %s: %s\n", markdownLine(doc.KindLabel), markdownLine(doc.DisplayName)); err != nil {
			return err
		}
	}
	return nil
}

func markdownLine(s string) string {
	return markdownEscaper.Replace(strings.Join(strings.Fields(s), " "))
}
