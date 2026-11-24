package main

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
)

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
	key, value, found := strings.Cut(s, ":")
	if !found {
		return s, ""
	}
	return strings.TrimSpace(key), strings.TrimSpace(value)
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
