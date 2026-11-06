
.PHONY: all
all: build run

.PHONY: run
run: verify
	go run ./cmd/docops -h

##
## Cross-compile the npmprobe Go binary for multiple platforms.
## Outputs to bin/ directory with platform-specific names.
##
.PHONY: build
build: fmt verify build-macos
	@echo "Build complete. Binaries in bin/"
	@ls -lh bin/*

build-macos:
	@mkdir -p bin
	@echo "Building npmprobe for macOS..."
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

