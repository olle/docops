package text

import (
	"fmt"
	"html"
	"io"
)

type htmlNode interface {
	writeHTML(w io.Writer) error
}

type htmlAttribute struct {
	name  string
	value string
}

type htmlElement struct {
	tag      string
	attrs    []htmlAttribute
	children []htmlNode
	void     bool
}

func (e htmlElement) writeHTML(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "<%s", e.tag); err != nil {
		return err
	}
	for _, attr := range e.attrs {
		if _, err := fmt.Fprintf(w, " %s=\"%s\"", attr.name, html.EscapeString(attr.value)); err != nil {
			return err
		}
	}
	if _, err := io.WriteString(w, ">"); err != nil {
		return err
	}
	if e.void {
		return nil
	}
	for _, child := range e.children {
		if err := child.writeHTML(w); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(w, "</%s>", e.tag)
	return err
}

type htmlText string

func (t htmlText) writeHTML(w io.Writer) error {
	_, err := io.WriteString(w, html.EscapeString(string(t)))
	return err
}

type htmlRaw string

func (r htmlRaw) writeHTML(w io.Writer) error {
	_, err := io.WriteString(w, string(r))
	return err
}

func htmlH1(text string) htmlNode {
	return htmlElement{tag: "h1", children: []htmlNode{htmlText(text)}}
}
