.PHONY: build clean test test-integration install dev-install deps fmt lint release-snapshot release man docs-gen check help

BINARY_NAME=exa
VERSION=$(shell git describe --tags --exact-match 2>/dev/null || git rev-parse --short HEAD 2>/dev/null || echo dev)
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION)"

all: build

build:
	@echo "Building $(BINARY_NAME)..."
	go build $(LDFLAGS) -o $(BINARY_NAME) .

clean:
	rm -f $(BINARY_NAME)
	rm -rf dist/
	go clean

test: build
	@echo "=== $(BINARY_NAME) smoke tests ==="
	@PASS=0; FAIL=0; \
	run_test() { \
		if eval ./$(BINARY_NAME) $$@ >/dev/null 2>&1; then \
			echo "  PASS: $$*"; PASS=$$((PASS+1)); \
		else \
			echo "  FAIL: $$*"; FAIL=$$((FAIL+1)); \
		fi; \
	}; \
	run_test --help; \
	run_test --version; \
	run_test docs; \
	run_test skill print; \
	run_test completion bash; \
	run_test completion zsh; \
	echo ""; \
	echo "Results: $$PASS passed, $$FAIL failed"; \
	[ "$$FAIL" -eq 0 ]

test-integration: build
	@echo "=== $(BINARY_NAME) integration tests ==="
	@if [ -z "$${EXA_API_KEY:-}" ]; then \
		echo "SKIP: EXA_API_KEY not set"; exit 0; \
	fi
	go test -v -run TestIntegration -count=1 -timeout 10m ./...

install: build
	sudo install -m 755 $(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)

dev-install: build
	sudo ln -sf $(PWD)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)

deps:
	go mod download
	go mod tidy

fmt:
	go fmt ./...

lint:
	golangci-lint run

release-snapshot:
	goreleaser release --snapshot --clean

release:
	goreleaser release --clean

check: fmt lint test

man:
	@echo "Generating man pages..."
	@mkdir -p docs/man
	go run ./cmd/gendocs man docs/man
	@echo "Man pages generated in docs/man/"

docs-gen:
	@echo "Generating markdown docs..."
	@mkdir -p docs/cli
	go run ./cmd/gendocs markdown docs/cli
	@echo "CLI docs generated in docs/cli/"

help:
	@echo "Available targets:"
	@echo "  build            - Build the binary"
	@echo "  clean            - Clean build artifacts"
	@echo "  test             - Run smoke tests (no API key needed)"
	@echo "  test-integration - Run integration tests (requires EXA_API_KEY)"
	@echo "  install          - Install to /usr/local/bin/"
	@echo "  dev-install      - Symlink for development"
	@echo "  deps             - Download and tidy dependencies"
	@echo "  fmt              - Format code"
	@echo "  lint             - Lint with golangci-lint"
	@echo "  release-snapshot - Local GoReleaser test build"
	@echo "  release          - Full GoReleaser release"
	@echo "  man              - Generate man pages via cobra/doc"
	@echo "  docs-gen         - Generate markdown CLI reference via cobra/doc"
	@echo "  check            - Run fmt, lint, test"
	@echo "  help             - Show this help"
