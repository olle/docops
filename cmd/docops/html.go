package main

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

type htmlWriter struct{}

func (htmlWriter) ContentType() string {
	return "text/html; charset=utf-8"
}

func (htmlWriter) Render(w io.Writer, name string, titles []string) error {
	heading := "Hello " + name + "!"
	body := []htmlNode{htmlH1(heading)}
	if len(titles) > 0 {
		items := make([]htmlNode, 0, len(titles))
		for _, title := range titles {
			items = append(items, htmlElement{tag: "li", children: []htmlNode{htmlText(title)}})
		}
		body = append(body, htmlElement{tag: "ul", children: items})
	}

	page := htmlElement{
		tag:   "html",
		attrs: []htmlAttribute{{name: "lang", value: "en"}},
		children: []htmlNode{
			htmlElement{tag: "head", children: []htmlNode{
				htmlElement{tag: "meta", attrs: []htmlAttribute{{name: "charset", value: "utf-8"}}, void: true},
				htmlElement{tag: "title", children: []htmlNode{htmlText(heading)}},
			}},
			htmlElement{tag: "body", children: []htmlNode{
				htmlElement{tag: "main", children: body},
			}},
		},
	}

	if err := htmlRaw("<!doctype html>\n").writeHTML(w); err != nil {
		return err
	}
	if err := page.writeHTML(w); err != nil {
		return err
	}
	_, err := io.WriteString(w, "\n")
	return err
}

func htmlH1(text string) htmlNode {
	return htmlElement{tag: "h1", children: []htmlNode{htmlText(text)}}
}
