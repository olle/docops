package text

import (
	"fmt"
	"io"

	"github.com/olle/docops/internal/view"
)

type PlainText struct{}

func (PlainText) Render(w io.Writer, snapshot view.Snapshot) error {
	if _, err := io.WriteString(w, snapshot.Greeting); err != nil {
		return err
	}
	summaries := view.Summaries(snapshot)
	if len(summaries) == 0 {
		_, err := io.WriteString(w, "\n")
		return err
	}
	if _, err := io.WriteString(w, "\n\nSpecifications:\n"); err != nil {
		return err
	}
	for _, doc := range summaries {
		if _, err := fmt.Fprintf(w, "- %s: %s\n", doc.KindLabel, doc.DisplayName); err != nil {
			return err
		}
	}
	return nil
}

func Fallback(w io.Writer, snapshot view.Snapshot) error {
	if _, err := fmt.Fprintf(w, "%s\n\n", snapshot.Greeting); err != nil {
		return err
	}
	summaries := view.Summaries(snapshot)
	if len(summaries) > 0 {
		if _, err := io.WriteString(w, "Specifications:\n"); err != nil {
			return err
		}
		for _, doc := range summaries {
			if _, err := fmt.Fprintf(w, "  * %s: %s\n", doc.KindLabel, doc.DisplayName); err != nil {
				return err
			}
		}
	}
	_, err := io.WriteString(w, "\nq: quit\n")
	return err
}

type Debug struct {
	Depth int
}

func (d Debug) Render(w io.Writer, snapshot view.Snapshot) error {
	_, err := io.WriteString(w, view.DebugAtDepth(snapshot, d.Depth))
	return err
}
