.DEFAULT_TARGET=help
VERSION:=$(shell cat VERSION)

## help: Display list of commands
.PHONY: help
help: Makefile
	@sed -n 's|^##||p' $< | column -t -s ':' | sed -e 's|^| |'

## build: Build tfautomv binary
.PHONY: build
build: fmt vet
	go build -o bin/tfautomv

## fmt: Format source code
.PHONY: fmt
fmt:
	go fmt ./...

## vet: Vet source code
.PHONY: vet
vet:
	go vet ./...

## test: Run unit tests
.PHONY: test
test:
	go test ./pkg/... -cover

## test-e2e: Run end-to-end tests
.PHONY: test-e2e
test-e2e:
	go test ./test/e2e/... -v

## release: Release a new version
.PHONY: release
release: test
	git tag -a "$(VERSION)" -m "$(VERSION)"
	git push origin "$(VERSION)"
	goreleaser release --clean --release-notes=docs/release-notes/$(VERSION).md

## release-dry-run: Test the release process without publishing
.PHONY: release-dry-run
release-dry-run:
	goreleaser release --snapshot --clean --skip=publish --release-notes=docs/release-notes/$(VERSION).md
