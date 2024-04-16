# Makefile
.PHONY: all test build clean

all: test build

test:
	go test ./...

build:
	go build -o ActualBudgetNormalizer cmd/main.go

clean:
	rm -f ActualBudgetNormalizer
	go clean

lint:
	golangci-lint run

format:
	go fmt ./...

run: build
	./ActualBudgetNormalizer

.DEFAULT_GOAL := all