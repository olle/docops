package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

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

func (k specificationKind) String() string {
	return k.Label()
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

func (s specificationStructure) String() string {
	var b strings.Builder
	b.WriteString("structure:\n")
	fmt.Fprintf(&b, "  channels: %s\n", debugStringList(s.channels))
	fmt.Fprintf(&b, "  operations: %s\n", debugStringList(s.operations))
	fmt.Fprintf(&b, "  paths: %s\n", debugStringList(s.paths))
	return b.String()
}

func (s specification) String() string {
	var b strings.Builder
	b.WriteString("specification:\n")
	fmt.Fprintf(&b, "  source: %q\n", s.source)
	fmt.Fprintf(&b, "  kind: %s\n", s.kind)
	fmt.Fprintf(&b, "  specVersion: %q\n", s.specVersion)
	fmt.Fprintf(&b, "  version: %q\n", s.version)
	fmt.Fprintf(&b, "  title: %q\n", s.title)
	fmt.Fprintf(&b, "  name: %q\n", s.name)
	fmt.Fprintf(&b, "  description: %q\n", s.description)
	b.WriteString(indentDebugBlock(s.structure.String(), "  "))
	return b.String()
}

func (s specificationSummary) String() string {
	return fmt.Sprintf("%s: %s", s.kind, s.title)
}

func (c catalog) String() string {
	var b strings.Builder
	b.WriteString("catalog:\n")
	fmt.Fprintf(&b, "  count: %d\n", len(c.documents))
	if len(c.documents) == 0 {
		return b.String()
	}
	b.WriteString("  documents:\n")
	for _, doc := range c.documents {
		b.WriteString(indentDebugBlock("- "+doc.String(), "    "))
	}
	return b.String()
}

func debugStringList(values []string) string {
	quoted := make([]string, len(values))
	for i, value := range values {
		quoted[i] = strconv.Quote(value)
	}
	return "[" + strings.Join(quoted, ", ") + "]"
}

func indentDebugBlock(text, prefix string) string {
	lines := strings.Split(strings.TrimSuffix(text, "\n"), "\n")
	for i, line := range lines {
		lines[i] = prefix + line
	}
	return strings.Join(lines, "\n") + "\n"
}
