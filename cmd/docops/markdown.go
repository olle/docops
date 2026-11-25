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
	if len(msg.titles) == 0 {
		return nil
	}
	if _, err := io.WriteString(w, "\n## AsyncAPI documents\n\n"); err != nil {
		return err
	}
	for _, title := range msg.titles {
		if _, err := fmt.Fprintf(w, "- %s\n", markdownLine(title)); err != nil {
			return err
		}
	}
	return nil
}

func markdownLine(s string) string {
	return markdownEscaper.Replace(strings.Join(strings.Fields(s), " "))
}
