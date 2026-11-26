package docops

import (
	"encoding/json"
	"os"
	"sort"
	"strconv"
	"strings"
)

type FileParser struct{}

func (FileParser) Parse(path string) (Specification, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Specification{}, err
	}
	if strings.HasPrefix(strings.TrimSpace(string(data)), "{") {
		return ParseJSONDocument(path, data)
	}
	return ParseYAMLDocument(path, data), nil
}

type rawSpecification struct {
	asyncAPI    string
	avroType    string
	openAPI     string
	jsonSchema  string
	id          string
	namespace   string
	name        string
	title       string
	description string
	version     string
	channels    []string
	operations  []string
	paths       []string
	properties  []string
	fields      []string
}

func ParseJSONDocument(path string, data []byte) (Specification, error) {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return Specification{}, err
	}

	parsed := rawSpecification{
		asyncAPI:    StringValue(raw["asyncapi"]),
		avroType:    StringValue(raw["type"]),
		openAPI:     StringValue(raw["openapi"]),
		jsonSchema:  StringValue(raw["$schema"]),
		id:          PreferredTitle(StringValue(raw["$id"]), StringValue(raw["id"])),
		namespace:   StringValue(raw["namespace"]),
		name:        StringValue(raw["name"]),
		title:       StringValue(raw["title"]),
		description: PreferredTitle(StringValue(raw["description"]), StringValue(raw["doc"])),
		version:     StringValue(raw["version"]),
		channels:    ObjectKeys(raw["channels"]),
		operations:  ObjectKeys(raw["operations"]),
		paths:       ObjectKeys(raw["paths"]),
		properties:  ObjectKeys(raw["properties"]),
		fields:      AvroFieldNames(raw["fields"]),
	}
	if info, ok := raw["info"].(map[string]any); ok {
		parsed.title = PreferredTitle(StringValue(info["title"]), StringValue(info["name"]))
		parsed.description = StringValue(info["description"])
		parsed.version = StringValue(info["version"])
	}
	parsed.title = PreferredTitle(parsed.title, parsed.name)
	return normalizeSpecification(path, parsed), nil
}

func ParseYAMLDocument(path string, data []byte) Specification {
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
			case "type":
				if raw.avroType == "" {
					raw.avroType = scalar(value)
				}
			case "namespace":
				raw.namespace = scalar(value)
			case "doc":
				if raw.description == "" {
					raw.description = scalar(value)
				}
			case "$schema":
				raw.jsonSchema = scalar(value)
			case "$id", "id":
				raw.id = scalar(value)
			case "name":
				raw.name = scalar(value)
			case "title":
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
			case "info":
				info := flowInfo(value)
				raw.title = PreferredTitle(info.title, info.name, raw.title)
				raw.description = PreferredTitle(info.description, raw.description)
				raw.version = PreferredTitle(info.version, raw.version)
				inInfo = value == ""
			case "channels", "operations", "paths", "properties", "fields":
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
		case "properties":
			raw.properties = appendUnique(raw.properties, key)
		case "fields":
			raw.fields = appendUnique(raw.fields, key)
		}
	}

	raw.title = PreferredTitle(raw.title, raw.name)
	return normalizeSpecification(path, raw)
}

func normalizeSpecification(path string, raw rawSpecification) Specification {
	kind := UnknownKind
	specVersion := ""
	if raw.asyncAPI != "" {
		kind = AsyncAPIKind
		specVersion = raw.asyncAPI
	}
	if raw.openAPI != "" {
		kind = OpenAPIKind
		specVersion = raw.openAPI
	}
	if raw.jsonSchema != "" {
		kind = JSONSchemaKind
		specVersion = schemaDialect(raw.jsonSchema)
	}
	if isAvroType(raw.avroType) {
		kind = AvroKind
	}

	id := raw.id
	if kind == AvroKind {
		id = PreferredTitle(id, qualifiedAvroName(raw.namespace, raw.name))
	}
	return Specification{
		Source:      path,
		Kind:        kind,
		SpecVersion: specVersion,
		ID:          id,
		Version:     raw.version,
		Title:       raw.title,
		Name:        raw.name,
		Description: raw.description,
		Structure: Structure{
			Channels:   SortedCopy(raw.channels),
			Operations: SortedCopy(raw.operations),
			Paths:      SortedCopy(raw.paths),
			Properties: SortedCopy(raw.properties),
			Fields:     SortedCopy(raw.fields),
		},
	}
}

func ObjectKeys(value any) []string {
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

func StringValue(value any) string {
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

func SortedCopy(values []string) []string {
	copied := append([]string(nil), values...)
	sort.Strings(copied)
	return copied
}

func splitField(s string) (string, string) {
	s = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(s), "- "))
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

func isAvroType(value string) bool {
	switch value {
	case "record", "error", "enum", "fixed":
		return true
	default:
		return false
	}
}

func AvroFieldNames(value any) []string {
	items, ok := value.([]any)
	if !ok {
		return nil
	}
	var names []string
	for _, item := range items {
		field, ok := item.(map[string]any)
		if !ok {
			continue
		}
		names = appendUnique(names, StringValue(field["name"]))
	}
	return names
}

func qualifiedAvroName(namespace, name string) string {
	if name == "" || namespace == "" || strings.HasPrefix(name, namespace+".") {
		return name
	}
	return namespace + "." + name
}

func schemaDialect(uri string) string {
	for _, dialect := range []string{"2020-12", "2019-09", "draft-07", "draft-06", "draft-04"} {
		if strings.Contains(uri, dialect) {
			return dialect
		}
	}
	return uri
}

func PreferredTitle(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
