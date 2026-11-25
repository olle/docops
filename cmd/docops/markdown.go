package main

import (
	"fmt"
	"io"
	"strings"
)

type markdownWriter struct{}

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

func (markdownWriter) ContentType() string {
	return "text/markdown; charset=utf-8"
}

func (markdownWriter) Render(w io.Writer, msg message) error {
	if _, err := fmt.Fprintf(w, "# %s\n", markdownLine(msg.greeting)); err != nil {
		return err
	}
	if len(msg.specifications) == 0 {
		return nil
	}
	if _, err := io.WriteString(w, "\n## Specifications\n\n"); err != nil {
		return err
	}
	for _, spec := range msg.specifications {
		if _, err := fmt.Fprintf(w, "- %s: %s\n", markdownLine(spec.kind.Label()), markdownLine(spec.title)); err != nil {
			return err
		}
	}
	return nil
}

func markdownLine(s string) string {
	return markdownEscaper.Replace(strings.Join(strings.Fields(s), " "))
}
