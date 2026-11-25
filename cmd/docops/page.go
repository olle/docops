package main

import (
	"io"

	"github.com/olle/docops/internal/assets"
)

type htmlWriter struct{}

func (htmlWriter) ContentType() string {
	return "text/html; charset=utf-8"
}

func (htmlWriter) Render(w io.Writer, msg message) error {
	body := []htmlNode{htmlH1(msg.greeting)}
	if len(msg.specifications) > 0 {
		items := make([]htmlNode, 0, len(msg.specifications))
		for _, spec := range msg.specifications {
			items = append(items, htmlElement{tag: "li", children: []htmlNode{
				htmlElement{tag: "strong", children: []htmlNode{htmlText(spec.kind.Label())}},
				htmlText(": " + spec.title),
			}})
		}
		body = append(body, htmlElement{tag: "h2", children: []htmlNode{htmlText("Specifications")}})
		body = append(body, htmlElement{tag: "ul", children: items})
	}

	page := htmlElement{
		tag:   "html",
		attrs: []htmlAttribute{{name: "lang", value: "en"}},
		children: []htmlNode{
			htmlElement{tag: "head", children: []htmlNode{
				htmlElement{tag: "meta", attrs: []htmlAttribute{{name: "charset", value: "utf-8"}}, void: true},
				htmlElement{tag: "title", children: []htmlNode{htmlText(msg.greeting)}},
				htmlElement{tag: "style", children: []htmlNode{htmlRaw(assets.CSS)}},
			}},
			htmlElement{tag: "body", children: []htmlNode{
				htmlElement{tag: "main", children: body},
			}},
		},
	}

	if _, err := io.WriteString(w, "<!doctype html>\n"); err != nil {
		return err
	}
	if err := page.writeHTML(w); err != nil {
		return err
	}
	_, err := io.WriteString(w, "\n")
	return err
}
