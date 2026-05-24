.PHONY: build test fmt run

build:
	go build -o bin/codeflow ./cmd/codeflow

fmt:
	gofmt -w $$(find . -name '*.go')

test:
	go test ./...

run:
	go run ./cmd/codeflow
