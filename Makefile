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

## fmt: Run go fmt against code
.PHONY: fmt
fmt:
	go fmt ./...

## vet: Run go vet against code
.PHONY: vet
vet:
	go vet ./...

## test: Run go test
.PHONY: test
test:
	go test ./...

## release: Release a new version
.PHONY: release
release: test
	git tag -a "$(VERSION)" -m "$(VERSION)"
	git push origin "$(VERSION)"
	goreleaser release --rm-dist
