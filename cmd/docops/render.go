package main

import (
	"fmt"
	"io"
)

type plainTextWriter struct{}

func (plainTextWriter) Render(w io.Writer, msg message) error {
	if _, err := io.WriteString(w, msg.greeting); err != nil {
		return err
	}
	if len(msg.specifications) == 0 {
		_, err := io.WriteString(w, "\n")
		return err
	}
	if _, err := io.WriteString(w, "\n\nSpecifications:\n"); err != nil {
		return err
	}
	for _, spec := range msg.specifications {
		if _, err := fmt.Fprintf(w, "- %s: %s\n", spec.kind.Label(), spec.title); err != nil {
			return err
		}
	}
	return nil
}
