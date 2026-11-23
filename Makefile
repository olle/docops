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
.PHONY: build
build: fmt verify build-macos
	@echo "Build complete. Binaries in bin/"
	@ls -lh bin/*

build-macos:
	@mkdir -p bin
	@echo "Building docops for macOS..."
	GOOS=darwin GOARCH=arm64 $(GO) build -o bin/docops ./cmd/docops

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

