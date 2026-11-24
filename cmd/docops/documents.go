package main

import (
	"errors"
	"fmt"
)

type document struct {
	asyncAPI string
	title    string
}

type documentFinder interface {
	Find(root string) ([]string, error)
}

type documentParser interface {
	Parse(path string) (document, error)
}

func collectTitles(root string, finder documentFinder, parser documentParser) ([]string, error) {
	paths, err := finder.Find(root)
	if err != nil {
		return nil, err
	}

	var titles []string
	var errs []error
	for _, path := range paths {
		doc, err := parser.Parse(path)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", path, err))
			continue
		}
		if doc.asyncAPI != "" && doc.title != "" {
			titles = append(titles, doc.title)
		}
	}
	return titles, errors.Join(errs...)
}
