package main

import (
	"encoding/json"
	"os"
	"sort"
	"strconv"
	"strings"
)

type fileDocumentParser struct{}

func (fileDocumentParser) Parse(path string) (specification, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return specification{}, err
	}
	if strings.HasPrefix(strings.TrimSpace(string(data)), "{") {
		return parseJSONDocument(path, data)
	}
	return parseYAMLDocument(path, data), nil
}

type rawSpecification struct {
	asyncAPI    string
	openAPI     string
	name        string
	title       string
	description string
	version     string
	channels    []string
	operations  []string
	paths       []string
}

func parseJSONDocument(path string, data []byte) (specification, error) {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return specification{}, err
	}

	parsed := rawSpecification{
		asyncAPI:   stringValue(raw["asyncapi"]),
		openAPI:    stringValue(raw["openapi"]),
		name:       stringValue(raw["name"]),
		channels:   objectKeys(raw["channels"]),
		operations: objectKeys(raw["operations"]),
		paths:      objectKeys(raw["paths"]),
	}
	if info, ok := raw["info"].(map[string]any); ok {
		parsed.title = preferredTitle(stringValue(info["title"]), stringValue(info["name"]))
		parsed.description = stringValue(info["description"])
		parsed.version = stringValue(info["version"])
	}
	parsed.title = preferredTitle(parsed.title, parsed.name)
	return normalizeSpecification(path, parsed), nil
}

func parseYAMLDocument(path string, data []byte) specification {
	raw := rawSpecification{}
	inInfo := false
	section := ""
	structureIndent := 0

	for _, rawLine := range strings.Split(string(data), "\n") {
		line := strings.TrimSuffix(rawLine, "\r")
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		indent := len(line) - len(strings.TrimLeft(line, " "))
		if indent == 0 {
			inInfo = false
			section = ""
			structureIndent = 0
			key, value := splitField(trimmed)
			switch key {
			case "asyncapi":
				raw.asyncAPI = scalar(value)
			case "openapi":
				raw.openAPI = scalar(value)
			case "name":
				raw.name = scalar(value)
			case "info":
				info := flowInfo(value)
				raw.title = preferredTitle(info.title, info.name, raw.title)
				raw.description = preferredTitle(info.description, raw.description)
				raw.version = preferredTitle(info.version, raw.version)
				inInfo = value == ""
			case "channels", "operations", "paths":
				if value == "" {
					section = key
				}
			}
			continue
		}

		key, value := splitField(trimmed)
		if inInfo {
			switch key {
			case "title":
				if raw.title == "" {
					raw.title = scalar(value)
				}
			case "name":
				if raw.title == "" {
					raw.title = scalar(value)
				}
			case "description":
				if raw.description == "" {
					raw.description = scalar(value)
				}
			case "version":
				if raw.version == "" {
					raw.version = scalar(value)
				}
			}
		}
		if section == "" || indent == 0 {
			continue
		}
		if structureIndent == 0 {
			structureIndent = indent
		}
		if indent != structureIndent {
			continue
		}
		switch section {
		case "channels":
			raw.channels = appendUnique(raw.channels, key)
		case "operations":
			raw.operations = appendUnique(raw.operations, key)
		case "paths":
			raw.paths = appendUnique(raw.paths, key)
		}
	}

	raw.title = preferredTitle(raw.title, raw.name)
	return normalizeSpecification(path, raw)
}

func normalizeSpecification(path string, raw rawSpecification) specification {
	kind := kindUnknown
	specVersion := ""
	if raw.asyncAPI != "" {
		kind = kindAsyncAPI
		specVersion = raw.asyncAPI
	}
	if raw.openAPI != "" {
		kind = kindOpenAPI
		specVersion = raw.openAPI
	}

	return specification{
		source:      path,
		kind:        kind,
		specVersion: specVersion,
		version:     raw.version,
		title:       raw.title,
		name:        raw.name,
		description: raw.description,
		structure: specificationStructure{
			channels:   sortedCopy(raw.channels),
			operations: sortedCopy(raw.operations),
			paths:      sortedCopy(raw.paths),
		},
	}
}

func objectKeys(value any) []string {
	object, ok := value.(map[string]any)
	if !ok {
		return nil
	}
	keys := make([]string, 0, len(object))
	for key := range object {
		keys = append(keys, key)
	}
	return keys
}

func stringValue(value any) string {
	switch typed := value.(type) {
	case string:
		return typed
	case float64:
		return strconv.FormatFloat(typed, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(typed)
	default:
		return ""
	}
}

func appendUnique(values []string, value string) []string {
	if value == "" {
		return values
	}
	for _, existing := range values {
		if existing == value {
			return values
		}
	}
	return append(values, value)
}

func sortedCopy(values []string) []string {
	copied := append([]string(nil), values...)
	sort.Strings(copied)
	return copied
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

func flowInfo(value string) rawSpecification {
	var info rawSpecification
	value = strings.TrimSpace(value)
	if !strings.HasPrefix(value, "{") {
		return info
	}
	value = strings.TrimSuffix(strings.TrimPrefix(value, "{"), "}")
	for _, field := range strings.Split(value, ",") {
		key, val := splitField(strings.TrimSpace(field))
		switch key {
		case "title":
			info.title = scalar(val)
		case "name":
			info.name = scalar(val)
		case "description":
			info.description = scalar(val)
		case "version":
			info.version = scalar(val)
		}
	}
	return info
}

func preferredTitle(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
