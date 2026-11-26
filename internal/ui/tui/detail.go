package tui

import (
	"strings"

	"github.com/olle/docops/internal/view"
)

func detailContent(doc view.Document) string {
	title := doc.DisplayName
	if title == "" {
		title = doc.Source
	}
	fields := []string{
		doc.KindLabel + ": " + title,
		"",
		"Source:      " + doc.Source,
		"Spec ver:    " + doc.SpecVersion,
		"ID:          " + doc.ID,
		"Version:     " + doc.Version,
		"Title:       " + doc.Title,
		"Name:        " + doc.Name,
		"Description: " + doc.Description,
		"",
		"Channels:   " + joinValues(doc.Structure.Channels),
		"Operations: " + joinValues(doc.Structure.Operations),
		"Paths:      " + joinValues(doc.Structure.Paths),
		"Properties: " + joinValues(doc.Structure.Properties),
		"Fields:     " + joinValues(doc.Structure.Fields),
	}
	return strings.Join(fields, "\n")
}

func joinValues(values []string) string {
	if len(values) == 0 {
		return "—"
	}
	return strings.Join(values, ", ")
}
