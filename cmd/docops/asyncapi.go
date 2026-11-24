package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var documentGlobs = []string{"*.yaml", "*.yml", "*.json"}

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

type globDocumentFinder struct{}

func (globDocumentFinder) Find(root string) ([]string, error) {
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

type fileDocumentParser struct{}

func (fileDocumentParser) Parse(path string) (document, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return document{}, err
	}
	if strings.HasPrefix(strings.TrimSpace(string(data)), "{") {
		return parseJSONDocument(data)
	}
	return parseYAMLDocument(data), nil
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

type jsonDocument struct {
	AsyncAPI string `json:"asyncapi"`
	Name     string `json:"name"`
	Info     struct {
		Title string `json:"title"`
		Name  string `json:"name"`
	} `json:"info"`
}

func parseJSONDocument(data []byte) (document, error) {
	var raw jsonDocument
	if err := json.Unmarshal(data, &raw); err != nil {
		return document{}, err
	}
	return document{
		asyncAPI: raw.AsyncAPI,
		title:    preferredTitle(raw.Info.Title, raw.Info.Name, raw.Name),
	}, nil
}

func parseYAMLDocument(data []byte) document {
	var doc document
	inInfo := false
	for _, rawLine := range strings.Split(string(data), "\n") {
		line := strings.TrimSuffix(rawLine, "\r")
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		indent := len(line) - len(strings.TrimLeft(line, " "))
		if indent == 0 {
			inInfo = false
			key, value := splitField(trimmed)
			switch key {
			case "asyncapi":
				doc.asyncAPI = scalar(value)
			case "info":
				title, name := flowTitle(value)
				doc.title = preferredTitle(title, name)
				inInfo = value == ""
			}
			continue
		}
		if inInfo {
			key, value := splitField(trimmed)
			switch key {
			case "title":
				if doc.title == "" {
					doc.title = scalar(value)
				}
			case "name":
				if doc.title == "" {
					doc.title = scalar(value)
				}
			}
		}
	}
	return doc
}

func splitField(s string) (string, string) {
	idx := strings.Index(s, ":")
	if idx < 0 {
		return s, ""
	}
	return strings.TrimSpace(s[:idx]), strings.TrimSpace(s[idx+1:])
}

func scalar(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if strings.HasPrefix(value, "\"") {
		if unquoted, err := strconv.Unquote(value); err == nil {
			return unquoted
		}
		if end := strings.LastIndex(value, "\""); end > 0 {
			return value[1:end]
		}
	}
	if strings.HasPrefix(value, "'") {
		if end := strings.Index(value[1:], "'"); end >= 0 {
			return strings.ReplaceAll(value[1:1+end], "''", "'")
		}
	}
	if idx := strings.Index(value, " #"); idx >= 0 {
		value = value[:idx]
	}
	return strings.TrimSpace(value)
}

func flowTitle(value string) (string, string) {
	value = strings.TrimSpace(value)
	if !strings.HasPrefix(value, "{") {
		return "", ""
	}
	value = strings.TrimSuffix(strings.TrimPrefix(value, "{"), "}")
	var title, name string
	for _, field := range strings.Split(value, ",") {
		key, val := splitField(strings.TrimSpace(field))
		switch key {
		case "title":
			title = scalar(val)
		case "name":
			name = scalar(val)
		}
	}
	return title, name
}

func preferredTitle(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
