package main

import (
	"fmt"
	"io"
)

type outputWriter interface {
	Render(w io.Writer, name string, titles []string) error
	ContentType() string
}

type plainTextWriter struct{}

func (plainTextWriter) ContentType() string {
	return "text/plain; charset=utf-8"
}

func (plainTextWriter) Render(w io.Writer, name string, titles []string) error {
	if _, err := fmt.Fprintf(w, "Hello %s!", name); err != nil {
		return err
	}
	if len(titles) == 0 {
		_, err := io.WriteString(w, "\n")
		return err
	}
	if _, err := io.WriteString(w, "\n\nAsyncAPI documents:\n"); err != nil {
		return err
	}
	for _, title := range titles {
		if _, err := fmt.Fprintf(w, "- %s\n", title); err != nil {
			return err
		}
	}
	return nil
}
