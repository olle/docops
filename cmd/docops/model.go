package main

import "sort"

type specificationKind string

const (
	kindUnknown  specificationKind = ""
	kindAsyncAPI specificationKind = "asyncapi"
	kindOpenAPI  specificationKind = "openapi"
)

func (k specificationKind) Label() string {
	switch k {
	case kindAsyncAPI:
		return "AsyncAPI"
	case kindOpenAPI:
		return "OpenAPI"
	default:
		return "Unknown"
	}
}

type specificationStructure struct {
	channels   []string
	operations []string
	paths      []string
}

type specification struct {
	source      string
	kind        specificationKind
	specVersion string
	version     string
	title       string
	name        string
	description string
	structure   specificationStructure
}

type specificationSummary struct {
	kind  specificationKind
	title string
}

type catalog struct {
	documents []specification
}

func (c catalog) Summaries() []specificationSummary {
	summaries := make([]specificationSummary, 0, len(c.documents))
	for _, doc := range c.documents {
		title := preferredTitle(doc.title, doc.name)
		if doc.kind == kindUnknown || title == "" {
			continue
		}
		summaries = append(summaries, specificationSummary{kind: doc.kind, title: title})
	}
	sort.Slice(summaries, func(i, j int) bool {
		if summaries[i].kind != summaries[j].kind {
			return summaries[i].kind.Label() < summaries[j].kind.Label()
		}
		return summaries[i].title < summaries[j].title
	})
	return summaries
}

func (c catalog) All() []specification {
	return append([]specification(nil), c.documents...)
}

func (c catalog) ByKind(kind specificationKind) []specification {
	var documents []specification
	for _, doc := range c.documents {
		if doc.kind == kind {
			documents = append(documents, doc)
		}
	}
	return documents
}

func (c catalog) Titles(kind specificationKind) []string {
	var titles []string
	for _, doc := range c.documents {
		if doc.kind == kind && doc.title != "" {
			titles = append(titles, doc.title)
		}
	}
	return sortedCopy(titles)
}

func (c catalog) Counts() map[specificationKind]int {
	counts := make(map[specificationKind]int)
	for _, doc := range c.documents {
		counts[doc.kind]++
	}
	return counts
}
