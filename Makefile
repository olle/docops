# Run Go through mise so the version pinned in mise.toml is always used.
GO ?= mise exec -- go

.PHONY: all
all: build run

.PHONY: run
run: verify
	$(GO) run ./cmd/docops -h

##
## Cross-compile the docops Go binary for macOS.
## Outputs to bin/docops.
##
.PHONY: build build-macos
build: fmt verify build-macos
	@echo "Build complete. Binaries in bin/"
	@ls -lh bin/*

build-macos:
	@mkdir -p bin
	@echo "Building docops for macOS..."
	GOOS=darwin GOARCH=arm64 $(GO) build -o bin/docops ./cmd/docops

##
## Build and serve the example project from its own folder.
##
.PHONY: example
example: build-macos
	@echo "Serving examples/hello-world at http://localhost:8080"
	cd examples/hello-world && ../../bin/docops -server

##
## Regenerate embedded assets (e.g. minified CSS) from their sources.
##
.PHONY: generate
generate:
	$(GO) generate ./...

##
## Format all Go source files.
##
.PHONY: format fmt
format fmt:
	$(GO) fmt ./...

##
## Run all tests to verify the codebase.
##
.PHONY: verify v
verify v:
	$(GO) test ./...

