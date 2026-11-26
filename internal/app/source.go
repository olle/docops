package app

import (
	"context"

	"github.com/olle/docops/internal/docops"
	"github.com/olle/docops/internal/view"
)

type Source struct {
	name    string
	root    string
	finder  docops.Finder
	sniffer docops.Sniffer
	parser  docops.Parser
}

func NewSource(root string) *Source {
	return &Source{
		name:    DefaultName,
		root:    root,
		finder:  docops.GlobFinder{},
		sniffer: docops.FileSniffer{},
		parser:  docops.FileParser{},
	}
}

func (s *Source) Snapshot(ctx context.Context) (view.Snapshot, error) {
	if err := ctx.Err(); err != nil {
		return view.Snapshot{}, err
	}
	catalog, err := docops.Collect(s.root, s.finder, s.sniffer, s.parser)
	snapshot := view.NewSnapshot(s.name, catalog, err)
	if ctx.Err() != nil {
		return snapshot, ctx.Err()
	}
	return snapshot, err
}

type FixedSource struct {
	snapshot view.Snapshot
}

func NewFixedSource(name string, documents []view.Document) FixedSource {
	return FixedSource{snapshot: view.Snapshot{
		Greeting:  "Hello " + name + "!",
		Documents: documents,
	}}
}

func (s FixedSource) Snapshot(context.Context) (view.Snapshot, error) {
	return s.snapshot, nil
}
