package tui

import (
	"github.com/olle/docops/internal/view"
)

type docItem struct {
	doc view.Document
}

func (i docItem) FilterValue() string {
	return i.doc.DisplayName + " " + i.doc.KindLabel + " " + i.doc.Source
}

func (i docItem) Title() string {
	return i.doc.KindLabel + ": " + i.doc.DisplayName
}

func (i docItem) Description() string {
	if i.doc.Description != "" {
		return i.doc.Description
	}
	return i.doc.Source
}
