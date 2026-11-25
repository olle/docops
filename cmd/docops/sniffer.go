package main

import (
	"errors"
	"io"
	"os"
	"strings"
)

const sniffLimit = 64 * 1024

type documentSniffer interface {
	Detect(path string) (specificationKind, error)
}

type fileDocumentSniffer struct{}

func (fileDocumentSniffer) Detect(path string) (specificationKind, error) {
	file, err := os.Open(path)
	if err != nil {
		return kindUnknown, err
	}
	defer file.Close()

	buffer := make([]byte, sniffLimit)
	n, err := io.ReadFull(file, buffer)
	if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
		return kindUnknown, err
	}
	return detectSpecificationKind(string(buffer[:n])), nil
}

func detectSpecificationKind(prefix string) specificationKind {
	for _, line := range strings.Split(prefix, "\n") {
		line = strings.TrimSuffix(line, "\r")
		if strings.HasPrefix(line, "asyncapi:") {
			return kindAsyncAPI
		}
		if strings.HasPrefix(line, "openapi:") {
			return kindOpenAPI
		}
	}

	async := strings.Index(prefix, `"asyncapi"`)
	open := strings.Index(prefix, `"openapi"`)
	switch {
	case async >= 0 && (open < 0 || async < open):
		return kindAsyncAPI
	case open >= 0:
		return kindOpenAPI
	default:
		return kindUnknown
	}
}
