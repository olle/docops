package docops

import (
	"errors"
	"fmt"
)

type Finder interface {
	Find(root string) ([]string, error)
}

type Sniffer interface {
	Detect(path string) (Kind, error)
}

type Parser interface {
	Parse(path string) (Specification, error)
}

func Collect(root string, finder Finder, sniffer Sniffer, parser Parser) (Catalog, error) {
	paths, err := finder.Find(root)
	if err != nil {
		return Catalog{}, err
	}

	var documents []Specification
	var errs []error
	for _, path := range paths {
		kind, err := sniffer.Detect(path)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", path, err))
			continue
		}
		if kind == UnknownKind {
			continue
		}
		doc, err := parser.Parse(path)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", path, err))
			continue
		}
		if doc.Kind != kind {
			continue
		}
		documents = append(documents, doc)
	}
	return NewCatalog(documents), errors.Join(errs...)
}
