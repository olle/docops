package main

import (
	"errors"
	"fmt"
)

type documentFinder interface {
	Find(root string) ([]string, error)
}

type documentParser interface {
	Parse(path string) (specification, error)
}

func collectDocuments(root string, finder documentFinder, sniffer documentSniffer, parser documentParser) (catalog, error) {
	paths, err := finder.Find(root)
	if err != nil {
		return catalog{}, err
	}

	var documents []specification
	var errs []error
	for _, path := range paths {
		kind, err := sniffer.Detect(path)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", path, err))
			continue
		}
		if kind == kindUnknown {
			continue
		}
		doc, err := parser.Parse(path)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", path, err))
			continue
		}
		if doc.kind != kind {
			continue
		}
		documents = append(documents, doc)
	}
	return catalog{documents: documents}, errors.Join(errs...)
}
