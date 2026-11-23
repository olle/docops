
.PHONY: all
all: build run

.PHONY: run
run: verify
	go run ./cmd/docops -h

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
	GOOS=darwin GOARCH=arm64 go build -o bin/docops ./cmd/docops

##
## Format all Go source files.
##
.PHONY: format fmt
format fmt:
	go fmt ./...

##
## Run all tests to verify the codebase.
##
.PHONY: verify v
verify v:
	go test ./...

