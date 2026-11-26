package view

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/olle/docops/internal/docops"
)

type Structure struct {
	Channels   []string
	Operations []string
	Paths      []string
	Properties []string
	Fields     []string
}

type Document struct {
	Source      string
	Kind        string
	KindLabel   string
	SpecVersion string
	ID          string
	Version     string
	Title       string
	Name        string
	DisplayName string
	Description string
	Structure   Structure
}

type Snapshot struct {
	Greeting  string
	Documents []Document
	Warnings  []string
}

func NewSnapshot(name string, catalog docops.Catalog, err error) Snapshot {
	documents := make([]Document, 0, len(catalog.Documents))
	for _, spec := range catalog.Documents {
		documents = append(documents, Document{
			Source:      spec.Source,
			Kind:        string(spec.Kind),
			KindLabel:   KindLabel(spec.Kind),
			SpecVersion: spec.SpecVersion,
			ID:          spec.ID,
			Version:     spec.Version,
			Title:       spec.Title,
			Name:        spec.Name,
			DisplayName: preferredTitle(spec.Title, spec.Name),
			Description: spec.Description,
			Structure: Structure{
				Channels:   append([]string(nil), spec.Structure.Channels...),
				Operations: append([]string(nil), spec.Structure.Operations...),
				Paths:      append([]string(nil), spec.Structure.Paths...),
				Properties: append([]string(nil), spec.Structure.Properties...),
				Fields:     append([]string(nil), spec.Structure.Fields...),
			},
		})
	}
	return Snapshot{
		Greeting:  "Hello " + name + "!",
		Documents: documents,
		Warnings:  warningLines(err),
	}
}

func KindLabel(kind docops.Kind) string {
	switch kind {
	case docops.AsyncAPIKind:
		return "AsyncAPI"
	case docops.AvroKind:
		return "Avro"
	case docops.JSONSchemaKind:
		return "JSON Schema"
	case docops.OpenAPIKind:
		return "OpenAPI"
	default:
		return "Unknown"
	}
}

func Summaries(snapshot Snapshot) []Document {
	summaries := make([]Document, 0, len(snapshot.Documents))
	for _, doc := range snapshot.Documents {
		if doc.Kind == "" || doc.DisplayName == "" {
			continue
		}
		summaries = append(summaries, doc)
	}
	sort.Slice(summaries, func(i, j int) bool {
		if summaries[i].KindLabel != summaries[j].KindLabel {
			return summaries[i].KindLabel < summaries[j].KindLabel
		}
		return summaries[i].DisplayName < summaries[j].DisplayName
	})
	return summaries
}

func Debug(snapshot Snapshot) string {
	var b strings.Builder
	b.WriteString("catalog:\n")
	fmt.Fprintf(&b, "  count: %d\n", len(snapshot.Documents))
	if len(snapshot.Documents) == 0 {
		return b.String()
	}
	b.WriteString("  documents:\n")
	for _, doc := range snapshot.Documents {
		b.WriteString(indentDebugBlock("- "+debugDocument(doc), "    "))
	}
	return b.String()
}

func debugDocument(doc Document) string {
	var b strings.Builder
	b.WriteString("specification:\n")
	fmt.Fprintf(&b, "  source: %q\n", doc.Source)
	fmt.Fprintf(&b, "  kind: %s\n", doc.KindLabel)
	fmt.Fprintf(&b, "  specVersion: %q\n", doc.SpecVersion)
	fmt.Fprintf(&b, "  id: %q\n", doc.ID)
	fmt.Fprintf(&b, "  version: %q\n", doc.Version)
	fmt.Fprintf(&b, "  title: %q\n", doc.Title)
	fmt.Fprintf(&b, "  name: %q\n", doc.Name)
	fmt.Fprintf(&b, "  description: %q\n", doc.Description)
	b.WriteString(indentDebugBlock(debugStructure(doc.Structure), "  "))
	return b.String()
}

func debugStructure(structure Structure) string {
	var b strings.Builder
	b.WriteString("structure:\n")
	fmt.Fprintf(&b, "  channels: %s\n", debugStringList(structure.Channels))
	fmt.Fprintf(&b, "  operations: %s\n", debugStringList(structure.Operations))
	fmt.Fprintf(&b, "  paths: %s\n", debugStringList(structure.Paths))
	fmt.Fprintf(&b, "  properties: %s\n", debugStringList(structure.Properties))
	fmt.Fprintf(&b, "  fields: %s\n", debugStringList(structure.Fields))
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

func preferredTitle(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func warningLines(err error) []string {
	if err == nil {
		return nil
	}
	return strings.Split(err.Error(), "\n")
}
