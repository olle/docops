package main

type specificationKind string

const (
	kindUnknown  specificationKind = ""
	kindAsyncAPI specificationKind = "asyncapi"
	kindOpenAPI  specificationKind = "openapi"
)

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

type catalog struct {
	documents []specification
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
