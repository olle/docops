package main

import (
	"io/fs"
	"path/filepath"
	"strings"
)

var documentGlobs = []string{"*.yaml", "*.yml", "*.json"}

type documentGlobFinder struct{}

func (documentGlobFinder) Find(root string) ([]string, error) {
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
		for _, pattern := range documentGlobs {
			matched, err := filepath.Match(pattern, entry.Name())
			if err != nil {
				return err
			}
			if matched {
				paths = append(paths, path)
				return nil
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return paths, nil
}
