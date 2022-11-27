TOOLS := github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1 \
		 google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0 \
		 google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1 \
		 github.com/bufbuild/buf/cmd/buf@v1.9.0 \
		 github.com/goreleaser/goreleaser@v1.13.0

PIP ?= pip

COVERAGE_FOLDER ?= coverage
CMD_FOLDER ?= ./cmd

COVERAGE_FILE := $(COVERAGE_FOLDER)/coverage.out
COVERAGE_FILE_HTML := $(COVERAGE_FOLDER)/coverage.html

TEST_COMMON_FLAGS ?= -v -failfast -timeout 1m

.PHONY: clean
clean:
	rm -rf $(COVERAGE_FOLDER)

.PHONY: install-tools
install-tools:
	for tool in $(TOOLS); do \
		go install $$tool; \
	done

	$(PIP) install pre-commit==2.20.0

.PHONY: setup
setup: install-tools
	pre-commit install

.PHONE: buf-build
buf-build:
	buf build

.PHONY: build
build: install-tools buf-build

.PHONY: buf-lint
buf-lint:
	buf lint

.PHONE: go-lint
go-lint:
	golangci-lint run

.PHONY: lint
lint: install-tools buf-lint go-lint

.PHONY: buf-generate
buf-generate:
	buf generate -o pkg/proto proto/internal

.PHONY: generate
generate: install-tools buf-generate

.PHONY: release-snapshot
release-snapshot: install-tools
	goreleaser release --snapshot --rm-dist

.PHONY:
test-unit:
	go test $(TEST_COMMON_FLAGS) -tags unit -race ./...

.PHONY: ensure-coverage
ensure-coverage:
	mkdir -p $(COVERAGE_FOLDER)

.PHONY: test-unit-coverage
test-unit-coverage: clean ensure-coverage
	go test $(TEST_COMMON_FLAGS) -tags unit -race -coverpkg ./... -covermode atomic -coverprofile $(COVERAGE_FILE) ./...

.PHONY: test-unit-coverage-html
test-unit-coverage-html: test-unit-coverage
	go tool cover -html $(COVERAGE_FILE) -o $(COVERAGE_FILE_HTML)
