package main

import (
	"io/fs"
	"path/filepath"
	"strings"
)

var documentNames = map[string]bool{
	"asyncapi.yaml": true,
	"asyncapi.yml":  true,
	"asyncapi.json": true,
	"openapi.yaml":  true,
	"openapi.yml":   true,
	"openapi.json":  true,
}

type documentNameFinder struct{}

func (documentNameFinder) Find(root string) ([]string, error) {
	var paths []string
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if path != root && (strings.HasPrefix(entry.Name(), ".") || entry.Name() == "bin") {
				return filepath.SkipDir
			}
			return nil
		}
		if documentNames[entry.Name()] {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return paths, nil
}
