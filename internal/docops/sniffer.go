package docops

import (
	"errors"
	"io"
	"os"
	"strings"
)

const sniffLimit = 64 * 1024

type FileSniffer struct{}

func (FileSniffer) Detect(path string) (Kind, error) {
	file, err := os.Open(path)
	if err != nil {
		return UnknownKind, err
	}
	defer file.Close()

	buffer := make([]byte, sniffLimit)
	n, err := io.ReadFull(file, buffer)
	if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
		return UnknownKind, err
	}
	return DetectSpecificationKind(string(buffer[:n])), nil
}

func DetectSpecificationKind(prefix string) Kind {
	for _, line := range strings.Split(prefix, "\n") {
		line = strings.TrimSuffix(line, "\r")
		if strings.HasPrefix(line, "asyncapi:") {
			return AsyncAPIKind
		}
		if strings.HasPrefix(line, "openapi:") {
			return OpenAPIKind
		}
		if strings.HasPrefix(line, "$schema:") && strings.Contains(line, "json-schema.org") {
			return JSONSchemaKind
		}
	}

	schema := strings.Index(prefix, `"$schema"`)
	if schema >= 0 && strings.Contains(prefix[schema:], "json-schema.org") {
		return JSONSchemaKind
	}
	if isAvroPrefix(prefix) {
		return AvroKind
	}

	async := strings.Index(prefix, `"asyncapi"`)
	open := strings.Index(prefix, `"openapi"`)
	switch {
	case async >= 0 && (open < 0 || async < open):
		return AsyncAPIKind
	case open >= 0:
		return OpenAPIKind
	default:
		return UnknownKind
	}
}

func isAvroPrefix(prefix string) bool {
	if !strings.Contains(prefix, `"type"`) || !strings.Contains(prefix, `"name"`) {
		return false
	}
	for _, avroType := range []string{`"record"`, `"error"`, `"enum"`, `"fixed"`} {
		if strings.Contains(prefix, avroType) {
			return true
		}
	}
	return false
}
