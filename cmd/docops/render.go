package main

import (
	"io"
)

type outputWriter interface {
	Render(w io.Writer, msg message) error
	ContentType() string
}

type plainTextWriter struct{}

func (plainTextWriter) ContentType() string {
	return "text/plain; charset=utf-8"
}

func (plainTextWriter) Render(w io.Writer, msg message) error {
	if _, err := io.WriteString(w, msg.greeting); err != nil {
		return err
	}
	if len(msg.titles) == 0 {
		_, err := io.WriteString(w, "\n")
		return err
	}
	if _, err := io.WriteString(w, "\n\nAsyncAPI documents:\n"); err != nil {
		return err
	}
	for _, title := range msg.titles {
		if _, err := io.WriteString(w, "- "+title+"\n"); err != nil {
			return err
		}
	}
	return nil
}
