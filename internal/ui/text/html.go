package text

import (
	"io"

	"github.com/olle/docops/internal/assets"
	"github.com/olle/docops/internal/view"
)

type HTML struct{}

func (HTML) ContentType() string {
	return "text/html; charset=utf-8"
}

func (HTML) Render(w io.Writer, snapshot view.Snapshot) error {
	body := []htmlNode{htmlH1(snapshot.Greeting)}
	if summaries := view.Summaries(snapshot); len(summaries) > 0 {
		items := make([]htmlNode, 0, len(summaries))
		for _, doc := range summaries {
			items = append(items, htmlElement{tag: "li", children: []htmlNode{
				htmlElement{tag: "strong", children: []htmlNode{htmlText(doc.KindLabel)}},
				htmlText(": " + doc.DisplayName),
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
				htmlElement{tag: "title", children: []htmlNode{htmlText(snapshot.Greeting)}},
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
