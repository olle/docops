package main

import (
	"io"
)

type htmlWriter struct{}

func (htmlWriter) ContentType() string {
	return "text/html; charset=utf-8"
}

func (htmlWriter) Render(w io.Writer, msg message) error {
	body := []htmlNode{htmlH1(msg.greeting)}
	if len(msg.titles) > 0 {
		items := make([]htmlNode, 0, len(msg.titles))
		for _, title := range msg.titles {
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
				htmlElement{tag: "title", children: []htmlNode{htmlText(msg.greeting)}},
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
